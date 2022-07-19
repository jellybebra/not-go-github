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
						testCase.owner, testCase.repo, exp, res)
					break
				}
			}
		} else {
			t.Errorf("Incorrect amount of contributors for %s/%s: expected %v, got %v",
				testCase.owner, testCase.repo, len(testCase.expected), len(contributors))
		}
	}
}

func TestGetRepositoryTags(t *testing.T) {
	const layout = "2006-01-02 15:04:05 -0700 MST"
	timeParsed, _ := time.Parse(layout, "2021-10-12 15:20:05 +0000 UTC")

	// Arrange
	testTable := []struct {
		owner    string
		repo     string
		expected []Tag
	}{
		{
			owner: "jostanise",
			repo:  "rsa_encrypted_local_chat",
			expected: []Tag{
				{
					Title: "v1.0",
					Hash:  "0480a292df58ba0bb4851bf828ed25efc56da813",
					Description: `# Никто больше не узнает, о чём ты разговариваешь 🗣️
					При запуске клиента автоматически генерируется пара ключей, которая используется на протяжении кратковременного (рекомендованно) общения. Сервер видит лишь открытые ключи шифрования и зашифрованные сообщения.
					
					# Просто использовать 😆
					
					- Если хочешь поговорить с кем-то по локальной сети, запусти **_local.exe_**
					- Если хочешь поговорить с кем-то удаленно, запусти _**server.exe**_ где угодно, подключись к нему и общайся!`,
					ZipLink:   "https://api.github.com/repos/jostanise/rsa_encrypted_local_chat/zipball/refs/tags/v1.0",
					CreatedAt: timeParsed,
				},
			},
		},
	}

	// Act
	ghs := getGHS()

	for _, testCase := range testTable {
		tags, _ := ghs.GetRepositoryTags(testCase.owner, testCase.repo)

		// Assert
		if len(tags) == len(testCase.expected) {
			for i := 0; i < len(tags); i++ {
				res := *tags[i]
				exp := testCase.expected[i]

				// exp.Description нужно нормально скопировать

				// if res.Description != exp.Description {
				// 	t.Errorf("Incorrect Description:\n\nexpected:\n\n %v\n\ngot\n\n%v", res.Description, exp.Description)
				// }

				if res.Title != exp.Title {
					t.Errorf("Incorrect Title: expected %v, got %v", res.Title, exp.Title)
				}

				if res.Hash != exp.Hash {
					t.Errorf("Incorrect Hash: expected %v, got %v", res.Hash, exp.Hash)
				}

				if res.ZipLink != exp.ZipLink {
					t.Errorf("Incorrect ZipLink: expected %v, got %v", res.ZipLink, exp.ZipLink)
				}

				if res.CreatedAt != exp.CreatedAt {
					t.Errorf("Incorrect CreatedAt: expected %v, got %v", res.CreatedAt, exp.CreatedAt)
				}
			}
		} else {
			t.Errorf("Incorrect amount of tags for %s/%s: expected %v, got %v",
				testCase.owner, testCase.repo, len(testCase.expected), len(tags))
		}
	}
}

func stringToTime(t string) time.Time {
	const layout = "2006-01-02 15:04:05 -0700 MST"
	parsedTime, _ := time.Parse(layout, t)
	return parsedTime
}

func TestGetRepositoryBranches(t *testing.T) {
	// Arrange
	testTable := []struct {
		owner    string
		repo     string
		expected []Branch
	}{
		{
			owner: "jostanise",
			repo:  "rsa_encrypted_local_chat",
			expected: []Branch{
				{
					Name:      "main",
					UpdatedAt: stringToTime("2021-10-12 15:20:05 +0000 UTC"),
				},
			},
		},
		{
			owner: "PeakIntegral",
			repo:  "cppLessons",
			expected: []Branch{
				{
					Name:      "create-sec-hero",
					UpdatedAt: stringToTime("2022-03-02 20:11:13 +0000 UTC"),
				},
				{
					Name:      "main",
					UpdatedAt: stringToTime("2022-03-05 18:17:54 +0000 UTC"),
				},
				{
					Name:      "patch-1",
					UpdatedAt: stringToTime("2022-03-05 22:52:20 +0000 UTC"),
				},
				{
					Name:      "yura",
					UpdatedAt: stringToTime("2022-03-02 14:59:38 +0000 UTC"),
				},
			},
		},
	}

	// Act
	ghs := getGHS()

	for _, testCase := range testTable {
		branches, _ := ghs.GetRepositoryBranches(testCase.owner, testCase.repo)

		// Assert
		if len(branches) == len(testCase.expected) {
			for i := 0; i < len(branches); i++ {
				res := *branches[i]
				exp := testCase.expected[i]

				if res != exp {
					t.Errorf("Incorrect branch data for %s/%s: expected %v, got %v",
						testCase.owner, testCase.repo, exp, res)
					break
				}
			}
		} else {
			t.Errorf("Incorrect amount of branches for %s/%s: expected %v, got %v",
				testCase.owner, testCase.repo, len(testCase.expected), len(branches))
		}
	}
}

/*
Get:

- GetBranchCommits
- GetUserRepositories
- GetIssues
- GetRepositoryPullRequests
- GetThreadsInfo (not implemented)

Ничего не возвращают:

- CreateRepository
- CreateBranch
- CreatePullRequest
- CreateTag
- DeleteTag
- SetAccessToRepository
- DenyAccessToRepository
- DeleteBranch

*/
