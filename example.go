package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v44/github"
	"github.com/joho/godotenv"
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
	GetRepositoryBranches(owner string, repositoryName string) ([]*Branch, error)

	// CreateBranch создает новую ветку
	CreateBranch(repoName, branchName string) error

	// DeleteBranch удаляет указанную ветку
	DeleteBranch(repoName, branchName string) error

	// GetBranchCommits возвращает коммиты указанной ветки
	GetBranchCommits(userName, repositoryName, branchName string) ([]*Commit, error)

	// GetRepositoryPullRequests получает информацию о запросах на слияние
	GetRepositoryPullRequests(repositoryName string) ([]*PullRequest, error)

	// CreatePullRequest создает новый запрос на слияние
	CreatePullRequest(sourceBranch, destBranch, title string) error

	// GetThreadsInfo получает информацию об обсуждениях конкретного запроса на слияние
	GetThreadsInfo(repositoryName string, pullRequestID int) ([]*Thread, error)

	// GetIssues получает информацию об опубликованных проблемах репозитория
	GetIssues(repositoryName string) ([]*Issue, error)

	// GetRepositoryContributors получает список соавторов репозитория
	GetRepositoryContributors(repositoryName string) ([]*User, error)

	// GetRepositoryTags возвращает информацию о тегах репозитория
	GetRepositoryTags(userName, repositoryName string) ([]*Tag, error)

	// CreateTag создает новый тег
	CreateTag(title string) error

	// DeleteTag удаляет тег по имени
	DeleteTag(repositoryName, tagName string) error

	// SetAccessToRepository предоставляет доступ к репозиторию указанному пользователю
	SetAccessToRepository(oppoUserName, repositoryName string) error

	// DenyAccessToRepository закрывает доступ к репозиторию указанному пользователю
	DenyAccessToRepository(oppoUserName, repositoryName string) error
}

// Структура, реализующая интерфейс GitServiceIFace
type gitHubService struct {
	client *github.Client
}

// NewGitHubService - конструктор gitHubService
func NewGitHubService(ctx context.Context) GitServiceIFace {
	// Используем Oauth2.0 в качестве протокола аутентификации
	ts := oauth2.StaticTokenSource(
		// Передаем Oauth2.0-токен, который можно получить в настройках профиля GitHub
		&oauth2.Token{AccessToken: goDotEnvVariable("KEY")},
	)
	tc := oauth2.NewClient(ctx, ts)

	// Запросы к GitHub API будут отправлены от имени аутентифицированного пользователя
	client := github.NewClient(tc)

	return &gitHubService{
		client: client,
	}
}

func goDotEnvVariable(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println("Error loading .env file")
		os.Exit(1)
	}

	return os.Getenv(key)
}

func getLanguages(rpGH *github.Repository, ghs *gitHubService) []struct {
	Name           string
	PercentOfUsage float64
} {

	languages, _, _ := ghs.client.Repositories.ListLanguages(context.Background(), rpGH.GetOwner().GetLogin(), rpGH.GetName())

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

	return Languages
}

// Необходимо реализовать нижепредставленные методы в соответствии со структурой интерфейса
//                                   |
//                                   |
//                                   |
//                                   V

func (ghs *gitHubService) GetUserInfo(userName string) (*User, error) {

	user, _, err := ghs.client.Users.Get(context.Background(), userName)

	us := User{
		UserName:       *user.Login,
		FullName:       *user.Name,
		FollowersCount: *user.Followers,
		FollowingCount: *user.Following,
	}

	return &us, err
}

func (ghs *gitHubService) GetUserRepositories(userName string) ([]*Repository, error) {

	repos, _, err := ghs.client.Repositories.List(context.Background(), userName, nil)

	var Repos []*Repository

	for _, repo := range repos {
		r := Repository{
			Name:                repo.GetName(),
			Description:         repo.GetDescription(),
			Link:                repo.GetHTMLURL(),
			IsPrivate:           repo.GetPrivate(),
			StarsCount:          repo.GetStargazersCount(),
			ForksCount:          repo.GetForksCount(),
			LastUpdatedTime:     repo.GetUpdatedAt().Time,
			programmingLanguage: getLanguages(repo, ghs),
		}

		Repos = append(Repos, &r)
	}

	return Repos, err
}

func (ghs *gitHubService) GetRepositoryByName(userName, repositoryName string) (*Repository, error) {
	repo, _, err := ghs.client.Repositories.Get(context.Background(), userName, repositoryName)

	rp := Repository{
		Name:                repo.GetName(),
		Description:         repo.GetDescription(),
		Link:                repo.GetHTMLURL(),
		IsPrivate:           repo.GetPrivate(),
		StarsCount:          repo.GetStargazersCount(),
		ForksCount:          repo.GetForksCount(),
		LastUpdatedTime:     repo.GetUpdatedAt().Time,
		programmingLanguage: getLanguages(repo, ghs),
	}

	return &rp, err
}

func (ghs *gitHubService) CreateRepository(repositoryName string) error { // scared to check
	r := &github.Repository{Name: &repositoryName}
	_, _, err := ghs.client.Repositories.Create(context.Background(), "", r)

	return err
}

func (ghs *gitHubService) GetRepositoryBranches(owner string, repositoryName string) ([]*Branch, error) {
	branches, _, err := ghs.client.Repositories.ListBranches(context.Background(), owner, repositoryName, nil)

	var Branches []*Branch

	for _, branch := range branches {
		commit := branch.GetCommit() // DOESN'T RETURN A GOOD COMMIT

		// вырезаем SHA и по SHA ищем коммит
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/", owner, repositoryName)
		sha := strings.Replace(commit.GetURL(), url, "", 1)
		goodCommit, _, _ := ghs.client.Git.GetCommit(context.Background(), owner, repositoryName, sha)

		br := Branch{
			Name:      *branch.Name,
			UpdatedAt: goodCommit.GetAuthor().GetDate(),
		}
		Branches = append(Branches, &br)
	}

	return Branches, err
}

func (ghs *gitHubService) CreateBranch(repoName, branchName string) error { // <--- сложно

	// https://stackoverflow.com/questions/9506181/github-api-create-branch

	//Ref := "refs/heads/" + branchName

	ref := github.Reference{Ref: &branchName}
	ghs.client.Git.CreateRef(context.Background(), "", repoName, &ref)

	return fmt.Errorf("implement me")
}

func (ghs *gitHubService) DeleteBranch(repoName, branchName string) error { // <--- сложно

	//Ref := "refs/heads/" + branchName

	ghs.client.Git.DeleteRef(context.Background(), "", repoName, branchName)

	return fmt.Errorf("implement me")
}

func findParentsOfCommit(ghs *gitHubService, commit *github.Commit, userName string, repositoryName string) ([]*github.Commit, error) {

	var Parents []*github.Commit

	// вырезаем SHA и по SHA ищем коммит
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/commits/", userName, repositoryName)
	sha := strings.Replace(commit.GetURL(), url, "", 1)
	goodCommit, _, _ := ghs.client.Git.GetCommit(context.Background(), userName, repositoryName, sha)

	Parents = append(Parents, goodCommit)

	for _, commit := range goodCommit.Parents {
		grandParents, _ := findParentsOfCommit(ghs, commit, userName, repositoryName)
		Parents = append(Parents, grandParents...)
	}

	return Parents, fmt.Errorf("implement me")
}

func (ghs *gitHubService) GetBranchCommits(userName, repositoryName, branchName string) ([]*Commit, error) {
	br, _, _ := ghs.client.Repositories.GetBranch(context.Background(), userName, repositoryName, branchName, true)
	lastCommit := br.GetCommit().GetCommit()

	commits, _ := findParentsOfCommit(ghs, lastCommit, userName, repositoryName) // пока что возвращает пустоту

	var Commits []*Commit

	for _, commit := range commits {
		pr := Commit{
			Hash:      commit.GetSHA(),
			Title:     commit.GetMessage(),
			CreatedAt: commit.GetAuthor().GetDate(),
		}
		Commits = append(Commits, &pr)

	}

	return Commits, fmt.Errorf("implement me")

}

func (ghs *gitHubService) GetRepositoryPullRequests(repositoryName string) ([]*PullRequest, error) { // <--- no username?
	opts := github.PullRequestListOptions{
		State: "all",
	}
	pullRequests, _, _ := ghs.client.PullRequests.List(context.Background(), "google", repositoryName, &opts)

	var PullRequests []*PullRequest

	for _, request := range pullRequests {
		pr := PullRequest{
			ID:           int(*request.ID),
			Title:        *request.Title,
			SourceBranch: *request.Head.Ref,
			TargetBranch: *request.Base.Ref,
			IsClosed:     *request.Locked, // это вообще то? выдаёт только false
		}
		PullRequests = append(PullRequests, &pr)

	}

	return PullRequests, fmt.Errorf("implement me")
}

func (ghs *gitHubService) CreatePullRequest(sourceBranch, destBranch, title string) error {
	return fmt.Errorf("implement me")
}

func (ghs *gitHubService) GetThreadsInfo(repositoryName string, pullRequestID int) ([]*Thread, error) { // <--- no username? не понял как делать
	// t, _, _ := ghs.client.PullRequests.Get()

	return nil, fmt.Errorf("implement me")
}

func (ghs *gitHubService) GetIssues(repositoryName string) ([]*Issue, error) { // <--- no username?
	opts := github.IssueListByRepoOptions{
		State: "all",
	}

	issues, _, _ := ghs.client.Issues.ListByRepo(context.Background(), "google", repositoryName, &opts)

	var Issues []*Issue
	for _, issue := range issues {

		i := Issue{
			Title:                   *issue.Title,
			IsClosed:                issue.GetLocked(), // это вообще то? выдаёт только false
			ResolvedPullRequestLink: issue.GetPullRequestLinks().GetURL(),
			CreatedAt:               *issue.CreatedAt,
			UpdatedAt:               *issue.UpdatedAt,
		}
		Issues = append(Issues, &i)
	}

	return Issues, fmt.Errorf("implement me")

}

func (ghs *gitHubService) GetRepositoryContributors(repositoryName string) ([]*User, error) { // <--- no username?

	opts := github.ListContributorsOptions{}
	contributors, _, _ := ghs.client.Repositories.ListContributors(context.Background(), "google", repositoryName, &opts)

	var Users []*User

	for _, user := range contributors {
		id, _, _ := ghs.client.Users.GetByID(context.Background(), user.GetID())

		i := User{
			UserName:       user.GetLogin(),
			FullName:       user.GetName(), // works?
			FollowersCount: id.GetFollowers(),
			FollowingCount: id.GetFollowing(),
		}
		Users = append(Users, &i)
	}

	return Users, fmt.Errorf("implement me")
}

func (ghs *gitHubService) GetRepositoryTags(userName, repositoryName string) ([]*Tag, error) {

	opts := github.ListOptions{}
	tags, _, _ := ghs.client.Repositories.ListTags(context.Background(), userName, repositoryName, &opts)

	var Tags []*Tag

	for _, tag := range tags {
		// вырезаем SHA и по SHA ищем коммит

		// url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/", userName, repositoryName) // вот тут ошибка потому что нам нужен не коммит
		// sha := strings.Replace(*tag.Commit.URL, url, "", 1)
		sha := *(tag.Commit.SHA)

		goodTag, _, _ := ghs.client.Git.GetTag(context.Background(), userName, repositoryName, sha)

		t := Tag{
			Title:       tag.GetName(),
			Hash:        goodTag.GetSHA(),     // не работает
			Description: goodTag.GetMessage(), // не работает
			ZipLink:     tag.GetZipballURL(),
			CreatedAt:   goodTag.GetTagger().GetDate(), // не работает
		}
		Tags = append(Tags, &t)

	}

	return Tags, fmt.Errorf("implement me")
}

func (ghs *gitHubService) CreateTag(title string) error {
	return fmt.Errorf("implement me")
}

func (ghs *gitHubService) DeleteTag(repositoryName, tagName string) error { // <--- no username?
	return fmt.Errorf("implement me")
}

func (ghs *gitHubService) SetAccessToRepository(oppoUserName, repositoryName string) error {
	return fmt.Errorf("implement me")
}

func (ghs *gitHubService) DenyAccessToRepository(oppoUserName, repositoryName string) error {
	return fmt.Errorf("implement me")
}

func main() {

	// Authorizing a client
	bruh := NewGitHubService(context.TODO())

	// GetBranchCommits
	fmt.Println("GetBranchCommits:")
	commits, _ := bruh.GetBranchCommits("jostanise", "rsa_encrypted_local_chat", "main")

	for _, commit := range commits {
		fmt.Printf("\tTitle:\t\t %v\n", commit.Title)
		fmt.Printf("\tHash:\t\t %v\n", commit.Hash)
		fmt.Printf("\tCreatedAt:\t %v\n", commit.CreatedAt)
		fmt.Println()

	}
	fmt.Println()

	// GetUserInfo получает основную информацию о пользователе
	user, _ := bruh.GetUserInfo("jostanise")
	fmt.Println("GetUserInfo:")
	fmt.Println("\tUserName:\t", user.UserName)
	fmt.Println("\tFullName:\t", user.FullName)
	fmt.Println("\tFollowersCount:\t", user.FollowersCount)
	fmt.Println("\tFollowingCount:\t", user.FollowingCount)
	fmt.Println()

	// GetUserRepositories
	fmt.Println("GetUserRepositories:")
	repos, _ := bruh.GetUserRepositories("")
	for i := 0; i < len(repos); i++ {
		repo := repos[i]
		fmt.Println("\tName:\t\t\t", repo.Name)
		fmt.Println("\tDescription:\t\t", repo.Description)
		fmt.Println("\tIsPrivate:\t\t", repo.IsPrivate)
		fmt.Println("\tStarsCount:\t\t", repo.StarsCount)
		fmt.Println("\tForksCount:\t\t", repo.ForksCount)
		fmt.Println("\tLastUpdatedTime:\t", repo.LastUpdatedTime)
		fmt.Println("\tprogrammingLanguage:\t", repo.programmingLanguage)
		fmt.Println("\tLink:\t\t\t", repo.Link)
		fmt.Println()
	}
	fmt.Println()

	// GetRepositoryByName
	fmt.Println("GetRepositoryByName:")
	repo, _ := bruh.GetRepositoryByName("jostanise", "rsa_encrypted_local_chat")
	fmt.Println("\tName:\t\t\t", repo.Name)
	fmt.Println("\tLastUpdatedTime:\t", repo.LastUpdatedTime)
	fmt.Println("\tprogrammingLanguage:\t", repo.programmingLanguage)
	fmt.Println()

	// GetRepositoryBranches
	fmt.Println("GetRepositoryBranches:")
	b, _ := bruh.GetRepositoryBranches("jostanise", "rsa_encrypted_local_chat")
	for i := 0; i < len(b); i++ {
		fmt.Println("\tBranch:", b[0].Name, "\tLast update:", b[0].UpdatedAt)
	}
	fmt.Println()

	// GetRepositoryPullRequests
	fmt.Println("GetRepositoryPullRequests:")
	pullreqs, _ := bruh.GetRepositoryPullRequests("go-github")
	for i := 0; i < len(pullreqs); i++ {
		pr := pullreqs[i]
		fmt.Println("\tID:\t\t", pr.ID)
		fmt.Println("\tTitle:\t\t", pr.Title)
		fmt.Println("\tSourceBranch:\t", pr.SourceBranch)
		fmt.Println("\tTargetBranch:\t", pr.TargetBranch)
		fmt.Println("\tIsClosed:\t", pr.IsClosed)
		fmt.Println()
	}
	fmt.Println()

	// GetThreadsInfo получает информацию об обсуждениях конкретного запроса на слияние

	// GetIssues получает информацию об опубликованных проблемах репозитория
	fmt.Println("GetIssues")
	issues, _ := bruh.GetIssues("go-github")
	for i := 0; i < len(issues); i++ {
		issue := issues[i]
		fmt.Println("\tTitle:\t\t\t", issue.Title)
		fmt.Println("\tIsClosed:\t\t", issue.IsClosed)
		fmt.Println("\tCreatedAt:\t\t", issue.CreatedAt)
		fmt.Println("\tUpdatedAt:\t\t", issue.UpdatedAt)
		fmt.Println("\tResolvedPullRequestLink:", issue.ResolvedPullRequestLink)
		fmt.Println()
	}
	fmt.Println()

	// GetRepositoryContributors получает список соавторов репозитория
	fmt.Println("GetRepositoryContributors:")
	contributors, _ := bruh.GetRepositoryContributors("go-github")
	for i := 0; i < len(contributors); i++ {
		contributor := contributors[i]
		fmt.Println("\tUsername:\t", contributor.UserName)
		fmt.Println("\tFullname:\t", contributor.FullName)
		fmt.Println("\tFollowersCount:\t", contributor.FollowersCount)
		fmt.Println("\tFollowingCount:\t", contributor.FollowingCount)
		fmt.Println()
	}
	fmt.Println()

	// GetRepositoryTags возвращает информацию о тегах репозитория
	fmt.Println("GetRepositoryTags:")
	tags, _ := bruh.GetRepositoryTags("google", "go-github")
	for _, tag := range tags {
		fmt.Println("\tTitle:\t\t", tag.Title)
		fmt.Println("\tHash:\t\t", tag.Hash)
		fmt.Println("\tDescription:\t", tag.Description)
		fmt.Println("\tZipLink:\t", tag.ZipLink)
		fmt.Println("\tCreatedAt:\t", tag.CreatedAt)
		fmt.Println()
	}

	// CreateRepository создает репозиторий с указанным именем
	// CreateBranch создает новую ветку
	// DeleteBranch удаляет указанную ветку
	// CreatePullRequest создает новый запрос на слияние
	// CreateTag создает новый тег
	// DeleteTag удаляет тег по имени
	// SetAccessToRepository предоставляет доступ к репозиторию указанному пользователю
	// DenyAccessToRepository закрывает доступ к репозиторию указанному пользователю

}
