package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/v39/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

func main() {
	fmt.Printf("Hello, world!\n")
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	ctx := context.Background()
	token := os.Getenv("RM_GITHUB_TOKEN")
	if token == "" {
		fmt.Println("Set the RM_GITHUB_TOKEN environment variable.")
		return
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	repoOwner := os.Getenv("RM_GITHUB_REPOSITORY_OWNER")
	repoName := os.Getenv("RM_GITHUB_REPOSITORY_NAME")
	if repoOwner == "" || repoName == "" {
		fmt.Println("Set the RM_GITHUB_REPOSITORY_OWNER and RM_GITHUB_REPOSITORY_NAME environment variables.")
		return
	}

	// Fetch all pull requests
	prs, _, err := client.PullRequests.List(ctx, repoOwner, repoName, nil)
	if err != nil {
		fmt.Printf("Error fetching pull requests: %v\n", err)
		return
	}

	for _, pr := range prs {
		// Fetch reviews for each PR
		reviews, _, err := client.PullRequests.ListReviews(ctx, repoOwner, repoName, *pr.Number, nil)
		if err != nil {
			fmt.Printf("Error fetching reviews for PR #%d: %v\n", *pr.Number, err)
			continue
		}

		if len(reviews) > 0 {
			// Assuming the first review in the list is the earliest
			firstReview := reviews[0]
			duration := firstReview.SubmittedAt.Sub(*pr.CreatedAt)
			fmt.Printf("PR #%d: Time to first review: %v\n", *pr.Number, duration)
		}
	}
}
