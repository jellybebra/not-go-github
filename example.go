package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
)

// User хранит краткую информацию о пользователе
type User struct {
	UserName       string // GitHub username пользователя
	FullName       string // Полное имя пользователя
	FollowersCount int    // Количество подписчиков
	FollowingCount int    // Количество подписок
}

// Repository хранит информацию о репозиториях пользователя
type Repository struct {
	Name            string    // Название репозитория
	Description     string    // Краткое описание репозитория
	Link            string    // Ссылка на репозиторий
	IsPrivate       bool      // Приватный репозиторий или открытый
	StarsCount      int       // Количество звезд
	ForksCount      int       // Количество форков (ответвлений, сделанных другими пользователями)
	LastUpdatedTime time.Time // Время последнего изменения

	// programmingLanguage хранит информацию об используемых языках программирования
	programmingLanguage []struct {
		Name           string  // Название языка программирования
		PercentOfUsage float64 // Процент использования в репозитории
	}
}

type Branch struct {
	Name      string    // Название ветки
	UpdatedAt time.Time // Дата последнего обновления
}

type Commit struct {
	Hash      string    // SHA коммита
	Title     string    // Сообщение коммита
	CreatedAt time.Time // Дата создания коммита
}

type Issue struct {
	Title                   string    // Тема issue
	IsClosed                bool      // Актуальная или разрешенная проблема
	ResolvedPullRequestLink string    // Ссылка на PR, в котором разрешена проблема
	CreatedAt               time.Time // Дата создания
	UpdatedAt               time.Time // Дата обновления
}

type PullRequest struct {
	ID           int    // Номер запроса на слияние (отображен в url как /pulls/{id})
	Title        string // Название запроса на слияние
	SourceBranch string // Название ветки-источника
	TargetBranch string // Название ветки-назначения
	IsClosed     bool   // Закрыт или открыт
}

type Thread struct {
	IsResolved bool // Закрытое или открытое обсуждение
}

type Tag struct {
	Title       string    // Название тега
	Hash        string    // SHA, хэш
	Description string    // Описание тега
	ZipLink     string    // Ссылка на скачивание архива
	CreatedAt   time.Time // Дата создания
}

type GitServiceIFace interface {
	// GetUserInfo получает основную информацию о пользователе
	GetUserInfo(userName string) (*User, error)

	// GetUserRepositories получает список всех репозиториев пользователя
	GetUserRepositories(userName string) ([]*Repository, error)

	// GetRepositoryByName получает информацию об указанном репозитории
	GetRepositoryByName(userName, repositoryName string) (*Repository, error)

	// CreateRepository создает репозиторий с указанным именем
	CreateRepository(repositoryName string) error

	// GetRepositoryBranches получает список всех веток репозитория
	GetRepositoryBranches(owner, repositoryName string) ([]*Branch, error)

	// CreateBranch создает новую ветку
	CreateBranch(userName, repoName, branchName, sha string) error

	// DeleteBranch удаляет указанную ветку
	DeleteBranch(userName, repoName, branchName string) error

	// GetBranchCommits возвращает коммиты указанной ветки
	GetBranchCommits(userName, repositoryName, branchName string) ([]*Commit, error)

	// GetRepositoryPullRequests получает информацию о запросах на слияние
	GetRepositoryPullRequests(userName, repositoryName string) ([]*PullRequest, error)

	// CreatePullRequest создает новый запрос на слияние
	CreatePullRequest(userName, repoName, sourceBranch, destBranch, title string) error

	// GetThreadsInfo получает информацию об обсуждениях конкретного запроса на слияние
	GetThreadsInfo(userName, repositoryName string, pullRequestID int) ([]*Thread, error)

	// GetIssues получает информацию об опубликованных проблемах репозитория
	GetIssues(userName, repositoryName string) ([]*Issue, error)

	// GetRepositoryContributors получает список соавторов репозитория
	GetRepositoryContributors(userName, repositoryName string) ([]*User, error)

	// GetRepositoryTags возвращает информацию о тегах репозитория
	GetRepositoryTags(userName, repositoryName string) ([]*Tag, error)

	// CreateTag создает новый тег
	CreateTag(userName, repositoryName, title, sha string) error

	// DeleteTag удаляет тег по имени
	DeleteTag(userName, repositoryName, tagName string) error

	// SetAccessToRepository предоставляет доступ к репозиторию указанному пользователю
	SetAccessToRepository(owner, repositoryName, oppoUserName string) error

	// DenyAccessToRepository закрывает доступ к репозиторию указанному пользователю
	DenyAccessToRepository(owner, repositoryName, oppoUserName string) error
}

// Структура, реализующая интерфейс GitServiceIFace
type gitHubService struct {
	client *github.Client
}

// NewGitHubService - конструктор gitHubService
func NewGitHubService(ctx context.Context) (GitServiceIFace, error) {
	// Используем Oauth2.0 в качестве протокола аутентификации
	ts := oauth2.StaticTokenSource(
		// Передаем Oauth2.0-токен, который можно получить в настройках профиля GitHub
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	// Запросы к GitHub API будут отправлены от имени аутентифицированного пользователя
	client := github.NewClient(tc)

	return &gitHubService{client: client}, nil
}

func getLanguages(rpGH *github.Repository, ghs *gitHubService) ([]struct {
	Name           string
	PercentOfUsage float64
}, error) {
	languages, _, err := ghs.client.Repositories.ListLanguages(context.Background(), rpGH.GetOwner().GetLogin(), rpGH.GetName())
	if err != nil {
		return nil, err
	}

	sum := 0
	for _, bytes := range languages {
		sum += bytes
	}

	var Languages []struct {
		Name           string
		PercentOfUsage float64
	}

	for lang, bytes := range languages {
		l := struct {
			Name           string
			PercentOfUsage float64
		}{
			lang,
			float64(bytes) / float64(sum),
		}
		Languages = append(Languages, l)
	}

	return Languages, nil
}

func findParentsOfCommit(ghs *gitHubService, commit *github.Commit, userName string, repositoryName string) ([]*github.Commit, error) {
	// Вырезаем SHA и по SHA ищем коммит (попробовать переделать)
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/commits/", userName, repositoryName)
	sha := strings.Replace(commit.GetURL(), url, "", 1)
	goodCommit, _, err := ghs.client.Git.GetCommit(context.Background(), userName, repositoryName, sha)
	if err != nil {
		return nil, fmt.Errorf("get commit: %w", err)
	}

	var Parents []*github.Commit
	Parents = append(Parents, goodCommit)

	for _, commit := range goodCommit.Parents {
		grandParents, err := findParentsOfCommit(ghs, commit, userName, repositoryName)
		if err != nil {
			return nil, fmt.Errorf("find parents of commit: %w", err)
		}

		Parents = append(Parents, grandParents...)
	}

	return Parents, nil
}

// Необходимо реализовать нижепредставленные методы в соответствии со структурой интерфейса
//                                   |
//                                   |
//                                   |
//                                   V

func (ghs *gitHubService) GetUserInfo(userName string) (*User, error) {
	ghUser, _, err := ghs.client.Users.Get(context.Background(), userName)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	user := User{
		UserName:       ghUser.GetLogin(),
		FullName:       ghUser.GetName(),
		FollowersCount: ghUser.GetFollowers(),
		FollowingCount: ghUser.GetFollowing(),
	}

	return &user, nil
}

func (ghs *gitHubService) GetUserRepositories(userName string) ([]*Repository, error) {
	opts := github.RepositoryListOptions{}

	// УБРАТЬ ПОТОМ
	repos, _, err := ghs.client.Repositories.List(context.Background(), userName, &opts)
	if err != nil {
		return nil, fmt.Errorf("list user repos: %w", err)
	}

	var Repos []*Repository
	for _, r := range repos {
		langs, err := getLanguages(r, ghs)
		if err != nil {
			return nil, fmt.Errorf("GetUserRepositories: %w", err)
		}

		repo := Repository{
			Name:                r.GetName(),
			Description:         r.GetDescription(),
			Link:                r.GetHTMLURL(),
			IsPrivate:           r.GetPrivate(),
			StarsCount:          r.GetStargazersCount(),
			ForksCount:          r.GetForksCount(),
			LastUpdatedTime:     r.GetUpdatedAt().Time,
			programmingLanguage: langs,
		}

		Repos = append(Repos, &repo)
	}

	return Repos, nil
}

func (ghs *gitHubService) GetRepositoryByName(userName, repositoryName string) (*Repository, error) {
	repo, _, err := ghs.client.Repositories.Get(context.Background(), userName, repositoryName)
	if err != nil {
		return nil, fmt.Errorf("get repo: %w", err)
	}

	langs, err := getLanguages(repo, ghs)
	if err != nil {
		return nil, fmt.Errorf("get langs for repo: %w", err)
	}

	rp := Repository{
		Name:                repo.GetName(),
		Description:         repo.GetDescription(),
		Link:                repo.GetHTMLURL(),
		IsPrivate:           repo.GetPrivate(),
		StarsCount:          repo.GetStargazersCount(),
		ForksCount:          repo.GetForksCount(),
		LastUpdatedTime:     repo.GetUpdatedAt().Time,
		programmingLanguage: langs,
	}

	return &rp, nil
}

func (ghs *gitHubService) CreateRepository(repositoryName string) error {
	repo := &github.Repository{Name: &repositoryName}
	_, _, err := ghs.client.Repositories.Create(context.Background(), "", repo)
	return err
}

func (ghs *gitHubService) GetRepositoryBranches(owner, repositoryName string) ([]*Branch, error) {
	opts := github.BranchListOptions{}
	branches, _, err := ghs.client.Repositories.ListBranches(context.Background(), owner, repositoryName, &opts)
	if err != nil {
		return nil, fmt.Errorf("list branches: %w", err)
	}

	var Branches []*Branch
	for _, branch := range branches {
		badCommit := branch.GetCommit() // returns "broken" UpdatedAt field
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/", owner, repositoryName)
		sha := strings.Replace(badCommit.GetURL(), url, "", 1)

		goodCommit, _, err := ghs.client.Git.GetCommit(context.Background(), owner, repositoryName, sha)
		if err != nil {
			return nil, fmt.Errorf("get commit: %w", err)
		}

		br := Branch{
			Name:      *branch.Name,
			UpdatedAt: goodCommit.GetAuthor().GetDate(),
		}
		Branches = append(Branches, &br)
	}

	return Branches, nil
}

func (ghs *gitHubService) CreateBranch(userName, repoName, branchName, sha string) error {
	// https://stackoverflow.com/questions/9506181/github-api-create-branch
	// https://api.github.com/repos/jostanise/rsa_encrypted_local_chat/branches
	// https://api.github.com/repos/jostanise/rsa_encrypted_local_chat/git/refs/heads
	// https://docs.github.com/en/rest/git/refs#create-a-reference

	ref := "refs/heads/" + branchName
	obj := github.GitObject{SHA: &sha}
	ghref := github.Reference{Ref: &ref, Object: &obj}
	_, _, err := ghs.client.Git.CreateRef(context.Background(), userName, repoName, &ghref)
	return err
}

func (ghs *gitHubService) DeleteBranch(userName, repoName, branchName string) error {
	ref := "refs/heads/" + branchName
	_, err := ghs.client.Git.DeleteRef(context.Background(), userName, repoName, ref)
	return err
}

func (ghs *gitHubService) GetBranchCommits(userName, repositoryName, branchName string) ([]*Commit, error) {
	br, _, err := ghs.client.Repositories.GetBranch(context.Background(), userName, repositoryName, branchName, true)
	if err != nil {
		return nil, fmt.Errorf("GetBranchCommits: %w", err)
	}

	lastCommit := br.GetCommit().GetCommit()
	commits, err := findParentsOfCommit(ghs, lastCommit, userName, repositoryName)
	if err != nil {
		return nil, fmt.Errorf("GetBranchCommits: %w", err)
	}

	var Commits []*Commit
	for _, c := range commits {
		commit := Commit{
			Hash:      c.GetSHA(),
			Title:     c.GetMessage(),
			CreatedAt: c.GetAuthor().GetDate(),
		}
		Commits = append(Commits, &commit)
	}

	return Commits, nil
}

func (ghs *gitHubService) GetRepositoryPullRequests(userName, repositoryName string) ([]*PullRequest, error) { // <--- no username?
	opts := github.PullRequestListOptions{State: "all"}
	pullRequests, _, err := ghs.client.PullRequests.List(context.Background(), userName, repositoryName, &opts)
	if err != nil {
		return nil, fmt.Errorf("list pull requests: %w", err)
	}

	var PullRequests []*PullRequest
	for _, r := range pullRequests {
		request := PullRequest{
			ID:           int(*r.ID),
			Title:        *r.Title,
			SourceBranch: *r.Head.Ref,
			TargetBranch: *r.Base.Ref,
			IsClosed:     *r.Locked, // не то? (выдаёт только false)
		}
		PullRequests = append(PullRequests, &request)
	}

	return PullRequests, nil
}

func (ghs *gitHubService) CreatePullRequest(userName, repoName, sourceBranch, destBranch, title string) error {
	pull := github.NewPullRequest{Title: &title, Head: &sourceBranch, Base: &destBranch}
	_, _, err := ghs.client.PullRequests.Create(context.Background(), userName, repoName, &pull)
	return err
}

func (ghs *gitHubService) GetThreadsInfo(userName, repositoryName string, pullRequestID int) ([]*Thread, error) { // must use GraphQL
	// REST API doesn't give "IsResolved" info
	// https://github.community/t/resolve-a-pr-review-comment-through-api/254182
	// https://github.com/shurcooL/githubv4
	// https://docs.github.com/en/graphql/overview/about-the-graphql-api
	// https://docs.github.com/en/graphql/reference/objects#repository

	return nil, fmt.Errorf("function is not implemented")
}

func (ghs *gitHubService) GetIssues(userName, repositoryName string) ([]*Issue, error) {
	opts := github.IssueListByRepoOptions{State: "all"}
	issues, _, err := ghs.client.Issues.ListByRepo(context.Background(), userName, repositoryName, &opts)
	if err != nil {
		return nil, fmt.Errorf("list issues by repo: %w", err)
	}

	var Issues []*Issue
	for _, issue := range issues {
		i := Issue{
			Title:                   *issue.Title,
			IsClosed:                issue.GetLocked(),                    // это вообще то? выдаёт только false
			ResolvedPullRequestLink: issue.GetPullRequestLinks().GetURL(), // Если не "", то IsClosed = true?
			CreatedAt:               *issue.CreatedAt,
			UpdatedAt:               *issue.UpdatedAt,
		}
		Issues = append(Issues, &i)
	}

	return Issues, nil
}

func (ghs *gitHubService) GetRepositoryContributors(userName, repositoryName string) ([]*User, error) {
	opts := github.ListContributorsOptions{}
	contributors, _, err := ghs.client.Repositories.ListContributors(context.Background(), userName, repositoryName, &opts)
	if err != nil {
		return nil, fmt.Errorf("list contributors: %w", err)
	}

	var Users []*User
	for _, user := range contributors {
		id, _, err := ghs.client.Users.GetByID(context.Background(), user.GetID())
		if err != nil {
			return nil, fmt.Errorf("get user by ID: %w", err)
		}

		user := User{
			UserName:       id.GetLogin(),
			FullName:       id.GetName(),
			FollowersCount: id.GetFollowers(),
			FollowingCount: id.GetFollowing(),
		}
		Users = append(Users, &user)
	}

	return Users, nil
}

func (ghs *gitHubService) GetRepositoryTags(userName, repositoryName string) ([]*Tag, error) {
	opts := github.ListOptions{}
	tags, _, err := ghs.client.Repositories.ListTags(context.Background(), userName, repositoryName, &opts)
	if err != nil {
		return nil, fmt.Errorf("list repo tags: %w", err)
	}

	var Tags []*Tag
	for _, tag := range tags {
		release, _, err := ghs.client.Repositories.GetReleaseByTag(context.Background(), userName, repositoryName, tag.GetName())
		if err != nil {
			return nil, fmt.Errorf("get release by tag: %w", err)
		}

		t := Tag{
			Title:       tag.GetName(),
			Hash:        tag.GetCommit().GetSHA(), // оно?
			Description: *release.Body,
			ZipLink:     tag.GetZipballURL(),
			CreatedAt:   release.GetCreatedAt().Time,
		}
		Tags = append(Tags, &t)
	}

	return Tags, nil
}

func (ghs *gitHubService) CreateTag(owner, repo, title, sha string) error {
	ref := "refs/tags/" + title
	obj := github.GitObject{SHA: &sha}
	ghref := github.Reference{Ref: &ref, Object: &obj}
	_, _, err := ghs.client.Git.CreateRef(context.Background(), owner, repo, &ghref)
	return err
}

func (ghs *gitHubService) DeleteTag(owner, repositoryName, tagName string) error {
	ref := "refs/tags/" + tagName
	_, err := ghs.client.Git.DeleteRef(context.Background(), owner, repositoryName, ref)
	return err
}

func (ghs *gitHubService) SetAccessToRepository(owner, repositoryName, oppoUserName string) error {
	opts := github.RepositoryAddCollaboratorOptions{Permission: "pull"}
	_, _, err := ghs.client.Repositories.AddCollaborator(context.Background(), owner, repositoryName, oppoUserName, &opts)
	return err
}

func (ghs *gitHubService) DenyAccessToRepository(owner, repositoryName, oppoUserName string) error {
	_, err := ghs.client.Repositories.RemoveCollaborator(context.Background(), owner, repositoryName, oppoUserName)
	if err != nil {
		return fmt.Errorf("removing collaborator: %w", err)
	}

	// Для случая, когда пользователь не принял приглашение
	opts := github.ListOptions{}
	invites, _, err := ghs.client.Repositories.ListInvitations(context.Background(), owner, repositoryName, &opts)
	if err != nil {
		return fmt.Errorf("list invitations: %w", err)
	}

	for _, invite := range invites {
		login := invite.Invitee.GetLogin()
		if login == oppoUserName {
			_, err := ghs.client.Repositories.DeleteInvitation(context.Background(), owner, repositoryName, *invite.ID)
			if err != nil {
				return fmt.Errorf("delete invitation: %w", err)
			}
			break
		}
	}

	return nil
}
