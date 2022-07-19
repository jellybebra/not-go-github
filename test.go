package main

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
)

func checkGetBranchCommits(ghs GitServiceIFace) {
	fmt.Println("GetBranchCommits:")
	commits, _ := ghs.GetBranchCommits("jostanise", "rsa_encrypted_local_chat", "main")
	for _, commit := range commits {
		fmt.Printf("\tTitle:\t\t %v\n", commit.Title)
		fmt.Printf("\tHash:\t\t %v\n", commit.Hash)
		fmt.Printf("\tCreatedAt:\t %v\n", commit.CreatedAt)
		fmt.Println()
	}
	fmt.Println()
}

func checkGetUserInfo(ghs GitServiceIFace) {
	user, _ := ghs.GetUserInfo("jostanise")
	fmt.Println("GetUserInfo:")
	fmt.Println("\tUserName:\t", user.UserName)
	fmt.Println("\tFullName:\t", user.FullName)
	fmt.Println("\tFollowersCount:\t", user.FollowersCount)
	fmt.Println("\tFollowingCount:\t", user.FollowingCount)
	fmt.Println()
}

func checkGetUserRepositories(ghs GitServiceIFace) {
	fmt.Println("GetUserRepositories:")
	repos, _ := ghs.GetUserRepositories("")
	for _, repo := range repos {
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

func checkGetRepositoryByName(ghs GitServiceIFace) {
	fmt.Println("GetRepositoryByName:")
	repo, _ := ghs.GetRepositoryByName("jostanise", "rsa_encrypted_local_chat")
	fmt.Println("\tName:\t\t\t", repo.Name)
	fmt.Println("\tLastUpdatedTime:\t", repo.LastUpdatedTime)
	fmt.Println("\tprogrammingLanguage:\t", repo.programmingLanguage)
	fmt.Println()
}

func checkGetRepositoryBranches(ghs GitServiceIFace) {
	fmt.Println("GetRepositoryBranches:")
	b, _ := ghs.GetRepositoryBranches("jostanise", "rsa_encrypted_local_chat")
	for i := 0; i < len(b); i++ {
		fmt.Println("\tBranch:", b[0].Name, "\tLast update:", b[0].UpdatedAt)
	}
	fmt.Println()
}

func checkGetRepositoryPullRequests(ghs GitServiceIFace) {
	fmt.Println("GetRepositoryPullRequests:")
	prs, _ := ghs.GetRepositoryPullRequests("google", "go-github")
	for _, pr := range prs {
		fmt.Println("\tID:\t\t", pr.ID)
		fmt.Println("\tTitle:\t\t", pr.Title)
		fmt.Println("\tSourceBranch:\t", pr.SourceBranch)
		fmt.Println("\tTargetBranch:\t", pr.TargetBranch)
		fmt.Println("\tIsClosed:\t", pr.IsClosed)
		fmt.Println()
	}
	fmt.Println()
}

func checkGetIssues(ghs GitServiceIFace) {
	fmt.Println("GetIssues")
	issues, _ := ghs.GetIssues("google", "go-github")
	for _, issue := range issues {
		fmt.Println("\tTitle:\t\t\t", issue.Title)
		fmt.Println("\tIsClosed:\t\t", issue.IsClosed)
		fmt.Println("\tCreatedAt:\t\t", issue.CreatedAt)
		fmt.Println("\tUpdatedAt:\t\t", issue.UpdatedAt)
		fmt.Println("\tResolvedPullRequestLink:", issue.ResolvedPullRequestLink)
		fmt.Println()
	}
	fmt.Println()
}

func checkGetRepositoryContributors(ghs GitServiceIFace) {
	fmt.Println("GetRepositoryContributors:")
	contributors, _ := ghs.GetRepositoryContributors("google", "go-github")
	for _, contributor := range contributors {
		fmt.Println("\tUsername:\t", contributor.UserName)
		fmt.Println("\tFullname:\t", contributor.FullName)
		fmt.Println("\tFollowersCount:\t", contributor.FollowersCount)
		fmt.Println("\tFollowingCount:\t", contributor.FollowingCount)
		fmt.Println()
	}
	fmt.Println()
}

func checkGetRepositoryTags(ghs GitServiceIFace) {
	fmt.Println("GetRepositoryTags:")
	tags, _ := ghs.GetRepositoryTags("google", "go-github")
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
	// Load token
	godotenv.Load(".env")

	// Authorizing a client
	ghs, _ := NewGitHubService(context.Background())

	// Get info
	checkGetUserRepositories(ghs)
	checkGetBranchCommits(ghs)
	checkGetUserInfo(ghs)
	checkGetRepositoryByName(ghs)
	checkGetRepositoryBranches(ghs)
	checkGetRepositoryPullRequests(ghs)
	checkGetIssues(ghs)
	checkGetRepositoryContributors(ghs)
	checkGetRepositoryTags(ghs)

	// No output
	ghs.CreateBranch("jostanise", "rsa_encrypted_local_chat", "tessst", "0480a292df58ba0bb4851bf828ed25efc56da813")
	ghs.DeleteBranch("jostanise", "rsa_encrypted_local_chat", "tessst")
	ghs.CreateTag("jostanise", "rsa_encrypted_local_chat", "tessst", "0480a292df58ba0bb4851bf828ed25efc56da813")
	ghs.DeleteTag("jostanise", "rsa_encrypted_local_chat", "tessst")
	ghs.CreateRepository("tessst")
	ghs.SetAccessToRepository("jostanise", "bruevich", "PeakIntegral")
	ghs.DenyAccessToRepository("jostanise", "bruevich", "PeakIntegral")
	ghs.CreatePullRequest("jostanise", "rsa_encrypted_local_chat", "tessst", "main", "tesst_to_main")

	// Not implemented
	// ghs.GetThreadsInfo("google", "go-github", 0)
}
