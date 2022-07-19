package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func getGHS() GitServiceIFace {
	godotenv.Load(".env")

	ghs, err := NewGitHubService(context.TODO())
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	return ghs
}

func TestGetUserInfo(t *testing.T) {
	// Arrange
	testTable := []struct {
		username string
		expected User
	}{
		{
			username: "jostanise",
			expected: User{
				UserName:       "jostanise",
				FullName:       "Mikhail Chestneyshy",
				FollowersCount: 0,
				FollowingCount: 0,
			},
		},
		{
			username: "PeakIntegral",
			expected: User{
				UserName:       "PeakIntegral",
				FullName:       "",
				FollowersCount: 0,
				FollowingCount: 0,
			},
		},
	}

	// Act
	ghs := getGHS()

	for _, testCase := range testTable {
		presult, _ := ghs.GetUserInfo(testCase.username)
		result := *presult

		// sresult := result.UserName + " " + result.FullName + " " + fmt.Sprint(result.FollowersCount) + " " + fmt.Sprint(result.FollowingCount)
		// t.Logf("Calling GetUserInfo(%v), result: %s",
		// 	testCase.username, sresult)

		// Assert
		if result != testCase.expected {
			t.Errorf("Incorrect result for %s", testCase.expected.UserName)
		}
	}
}

func TestGetRepositoryByName(t *testing.T) {
	// Arrange
	testTable := []struct {
		owner    string
		repo     string
		expected Repository
	}{
		{
			owner: "jostanise",
			repo:  "not-go-github",
			expected: Repository{
				Name:            "not-go-github",
				Description:     "Библиотека для управления GitHub на Go",
				Link:            "https://github.com/jostanise/not-go-github",
				IsPrivate:       true,
				StarsCount:      0,
				ForksCount:      0,
				LastUpdatedTime: time.Time{}, // как "2022-07-08 09:33:31 +0000 UTC" преобразовать для теста
				programmingLanguage: []struct {
					Name           string
					PercentOfUsage float64
				}{
					{
						Name:           "Go",
						PercentOfUsage: 1,
					},
				},
			},
		},

		{
			owner: "jostanise",
			repo:  "rsa_encrypted_local_chat",
			expected: Repository{
				Name:            "rsa_encrypted_local_chat",
				Description:     "Secure chatting with a friend over local network.",
				Link:            "https://github.com/jostanise/rsa_encrypted_local_chat",
				IsPrivate:       false,
				StarsCount:      0,
				ForksCount:      0,
				LastUpdatedTime: time.Time{}, // как "2021-10-12 15:20:12 +0000 UTC" преобразовать для теста
				programmingLanguage: []struct {
					Name           string
					PercentOfUsage float64
				}{
					{
						Name:           "Python",
						PercentOfUsage: 0.9896706768744683,
					},
					{
						Name:           "Batchfile",
						PercentOfUsage: 0.010329323125531656,
					},
				},
			},
		},
	}

	// Act
	ghs := getGHS()

	for _, testCase := range testTable {
		repo, _ := ghs.GetRepositoryByName(testCase.owner, testCase.repo)
		result := *repo

		// Assert
		if result.Name != testCase.expected.Name {
			t.Errorf("Incorrect NAME result for %s/%s", testCase.owner, testCase.repo)
		}

		if result.Description != testCase.expected.Description {
			t.Errorf("Incorrect DESCRIPTION result for %s/%s", testCase.owner, testCase.repo)
		}

		// и т.д. доделать

		// slice can only be compared to nil, придётся сравнивать элементы обоих слайсов
		if result.programmingLanguage[0] != testCase.expected.programmingLanguage[0] {
			t.Errorf("Incorrect result for %s/%s", testCase.owner, testCase.repo)
		}
	}
}

func TestGetRepositoryContributors(t *testing.T) {
	// Arrange
	testTable := []struct {
		owner    string
		repo     string
		expected []User
	}{
		{
			owner: "jostanise",
			repo:  "rsa_encrypted_local_chat",
			expected: []User{
				{
					UserName:       "jostanise",
					FullName:       "Mikhail Chestneyshy",
					FollowersCount: 0,
					FollowingCount: 0,
				},
				{
					UserName:       "PeakIntegral",
					FullName:       "",
					FollowersCount: 0,
					FollowingCount: 0,
				},
			},
		},
		{
			owner: "jostanise",
			repo:  "jostanise",
			expected: []User{
				{
					UserName:       "jostanise",
					FullName:       "Mikhail Chestneyshy",
					FollowersCount: 0,
					FollowingCount: 0,
				},
			},
		},
	}

	// Act
	ghs := getGHS()

	for _, testCase := range testTable {
		contributors, _ := ghs.GetRepositoryContributors(testCase.owner, testCase.repo)

		// Assert
		if len(contributors) == len(testCase.expected) {
			for i := 0; i < len(contributors); i++ {
				res := *contributors[i]
				exp := testCase.expected[i]

				if res != exp {
					t.Errorf("Incorrect user data for %s/%s: expected %v, got %v",
						testCase.owner, testCase.repo, res, exp)
					break
				}
			}
		} else {
			t.Errorf("Incorrect amount of contributors for %s/%s: expected %v, got %v",
				testCase.owner, testCase.repo, len(testCase.expected), len(contributors))
		}
	}
}

/*
Get:

- GetBranchCommits
- GetUserRepositories
- GetRepositoryTags
- GetIssues
- GetRepositoryPullRequests
- GetRepositoryBranches
- GetThreadsInfo (not implemented)

Ничего не возвращает:

- CreateRepository
- CreateBranch
- CreatePullRequest
- CreateTag
- DeleteTag
- SetAccessToRepository
- DenyAccessToRepository
- DeleteBranch

*/
