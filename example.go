package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/go-github/v44/github"
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
	GetRepositoryBranches(repositoryName string) ([]*Branch, error)

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
		// Токен необходимо передавать из переменных окружения!
		// Пример библиотеки: https://github.com/caarlos0/env
		&oauth2.Token{AccessToken: "... your access token ..."},
	)
	tc := oauth2.NewClient(ctx, ts)

	// Запросы к GitHub API будут отправлены от имени аутентифицированного пользователя
	client := github.NewClient(tc)

	return &gitHubService{
		client: client,
	}
}

// Необходимо реализовать нижепредставленные методы в соответствии со структурой интерфейса
//                                   |
//                                   |
//                                   |
//                                   V

func (ghs *gitHubService) GetUserInfo(userName string) (*User, error) {
	return nil, fmt.Errorf("implement me")
}

func (ghs *gitHubService) GetUserRepositories(userName string) ([]*Repository, error) {
	return nil, fmt.Errorf("implement me")
}

func (ghs *gitHubService) GetRepositoryByName(userName, repositoryName string) (*Repository, error) {
	return nil, fmt.Errorf("implement me")
}

func (ghs *gitHubService) CreateRepository(repositoryName string) error {
	return fmt.Errorf("implement me")
}

func (ghs *gitHubService) GetRepositoryBranches(repositoryName string) ([]*Branch, error) {
	return nil, fmt.Errorf("implement me")
}

func (ghs *gitHubService) CreateBranch(repoName, branchName string) error {
	return fmt.Errorf("implement me")
}

func (ghs *gitHubService) DeleteBranch(repoName, branchName string) error {
	return fmt.Errorf("implement me")
}

func (ghs *gitHubService) GetBranchCommits(userName, repositoryName, branchName string) ([]*Commit, error) {
	return nil, fmt.Errorf("implement me")
}

func (ghs *gitHubService) GetRepositoryPullRequests(repositoryName string) ([]*PullRequest, error) {
	return nil, fmt.Errorf("implement me")
}

func (ghs *gitHubService) CreatePullRequest(sourceBranch, destBranch, title string) error {
	return fmt.Errorf("implement me")
}

func (ghs *gitHubService) GetThreadsInfo(repositoryName string, pullRequestID int) ([]*Thread, error) {
	return nil, fmt.Errorf("implement me")
}

func (ghs *gitHubService) GetIssues(repositoryName string) ([]*Issue, error) {
	return nil, fmt.Errorf("implement me")
}

func (ghs *gitHubService) GetRepositoryContributors(repositoryName string) ([]*User, error) {
	return nil, fmt.Errorf("implement me")
}

func (ghs *gitHubService) GetRepositoryTags(userName, repositoryName string) ([]*Tag, error) {
	return nil, fmt.Errorf("implement me")
}

func (ghs *gitHubService) CreateTag(title string) error {
	return fmt.Errorf("implement me")
}

func (ghs *gitHubService) DeleteTag(repositoryName, tagName string) error {
	return fmt.Errorf("implement me")
}

func (ghs *gitHubService) SetAccessToRepository(oppoUserName, repositoryName string) error {
	return fmt.Errorf("implement me")
}

func (ghs *gitHubService) DenyAccessToRepository(oppoUserName, repositoryName string) error {
	return fmt.Errorf("implement me")
}
