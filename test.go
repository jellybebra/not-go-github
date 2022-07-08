package main

import (
	"context"
	"fmt"
)

func main() {

	// Authorizing a client
	bruh := NewGitHubService(context.TODO())

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
		fmt.Println("\tLastUpdatedTime:\t", repo.LastUpdatedTime)
		fmt.Println("\tprogrammingLanguage:\t", repo.programmingLanguage)
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

	// GetBranchCommits

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
	fmt.Println("GetRepositoryContributors")
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

	// CreateRepository создает репозиторий с указанным именем
	// CreateBranch создает новую ветку
	// DeleteBranch удаляет указанную ветку
	// CreatePullRequest создает новый запрос на слияние
	// CreateTag создает новый тег
	// DeleteTag удаляет тег по имени
	// SetAccessToRepository предоставляет доступ к репозиторию указанному пользователю
	// DenyAccessToRepository закрывает доступ к репозиторию указанному пользователю

}
