package main

import (
	"context"
	"fmt"
)

func checkGetBranchCommits(bruh GitServiceIFace) {

	fmt.Println("GetBranchCommits:")
	commits, _ := bruh.GetBranchCommits("jostanise", "rsa_encrypted_local_chat", "main")

	for _, commit := range commits {
		fmt.Printf("\tTitle:\t\t %v\n", commit.Title)
		fmt.Printf("\tHash:\t\t %v\n", commit.Hash)
		fmt.Printf("\tCreatedAt:\t %v\n", commit.CreatedAt)
		fmt.Println()

	}
	fmt.Println()
}

func checkGetUserInfo(bruh GitServiceIFace) {
	user, _ := bruh.GetUserInfo("jostanise")
	fmt.Println("GetUserInfo:")
	fmt.Println("\tUserName:\t", user.UserName)
	fmt.Println("\tFullName:\t", user.FullName)
	fmt.Println("\tFollowersCount:\t", user.FollowersCount)
	fmt.Println("\tFollowingCount:\t", user.FollowingCount)
	fmt.Println()
}

func checkGetUserRepositories(bruh GitServiceIFace) {
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
}

func checkGetRepositoryByName(bruh GitServiceIFace) {
	fmt.Println("GetRepositoryByName:")
	repo, _ := bruh.GetRepositoryByName("jostanise", "rsa_encrypted_local_chat")
	fmt.Println("\tName:\t\t\t", repo.Name)
	fmt.Println("\tLastUpdatedTime:\t", repo.LastUpdatedTime)
	fmt.Println("\tprogrammingLanguage:\t", repo.programmingLanguage)
	fmt.Println()
}

func checkGetRepositoryBranches(bruh GitServiceIFace) {
	fmt.Println("GetRepositoryBranches:")
	b, _ := bruh.GetRepositoryBranches("jostanise", "rsa_encrypted_local_chat")
	for i := 0; i < len(b); i++ {
		fmt.Println("\tBranch:", b[0].Name, "\tLast update:", b[0].UpdatedAt)
	}
	fmt.Println()
}

func checkGetRepositoryPullRequests(bruh GitServiceIFace) {
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
}

func checkGetIssues(bruh GitServiceIFace) {
	fmt.Println("GetIssues")
	issues, _ := bruh.GetIssues("google", "go-github")
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
}

func checkGetRepositoryContributors(bruh GitServiceIFace) {
	fmt.Println("GetRepositoryContributors:")
	contributors, _ := bruh.GetRepositoryContributors("google", "go-github")
	for i := 0; i < len(contributors); i++ {
		contributor := contributors[i]
		fmt.Println("\tUsername:\t", contributor.UserName)
		fmt.Println("\tFullname:\t", contributor.FullName)
		fmt.Println("\tFollowersCount:\t", contributor.FollowersCount)
		fmt.Println("\tFollowingCount:\t", contributor.FollowingCount)
		fmt.Println()
	}
	fmt.Println()
}

func checkGetRepositoryTags(bruh GitServiceIFace) {
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
}

func main() {

	// Authorizing a client
	bruh := NewGitHubService(context.TODO())

	// checkGetBranchCommits(bruh)
	// checkGetUserInfo(bruh)
	// checkGetUserRepositories(bruh)
	// checkGetRepositoryByName(bruh)
	// checkGetRepositoryBranches(bruh)
	// checkGetRepositoryPullRequests(bruh)
	// checkGetIssues(bruh)
	// checkGetRepositoryContributors(bruh)
	// checkGetRepositoryTags(bruh)

	// bruh.CreateBranch("jostanise", "rsa_encrypted_local_chat", "testbruh", "0480a292df58ba0bb4851bf828ed25efc56da813")
	// bruh.DeleteBranch("jostanise", "rsa_encrypted_local_chat", "testbruh")
	// bruh.CreateTag("jostanise", "rsa_encrypted_local_chat", "bruh", "0480a292df58ba0bb4851bf828ed25efc56da813")
	// bruh.DeleteTag("jostanise", "rsa_encrypted_local_chat", "bruh")
	// bruh.CreateRepository("bruuuh")
	// bruh.SetAccessToRepository("jostanise", "bruevich", "PeakIntegral")
	// bruh.DenyAccessToRepository("jostanise", "bruevich", "PeakIntegral")

	err := bruh.CreatePullRequest("jostanise", "rsa_encrypted_local_chat", "testbruh", "main", "testy besty")
	fmt.Printf("err: %v\n", err)

	// bruh.GetThreadsInfo("google", "go-github", 0)

}
