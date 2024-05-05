package main

import (
	"context"
	"fmt"
	"os"

	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
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

	// Store the time to first review for each PR
	times := make([]float64, 0)

	for _, pr := range prs {
		// Fetch all reviews for the PR
		reviews, _, err := client.PullRequests.ListReviews(ctx, repoOwner, repoName, *pr.Number, nil)
		if err != nil {
			fmt.Printf("Error fetching reviews for PR %d: %v\n", *pr.Number, err)
			return
		}
		fmt.Printf("PR %d has %d reviews\n", *pr.Number, len(reviews))

		// Calculate the time to first review
		if len(reviews) > 0 {
			timeToFirstReview := reviews[0].SubmittedAt.Sub(*pr.CreatedAt).Hours()
			times = append(times, timeToFirstReview)
		}
	}

	// Draw a graph of the average time to first review for each PR
	drawAverageFirstReviewTimeGraph(&times)
}

func drawAverageFirstReviewTimeGraph(times *[]float64) {
	p := plot.New()
	p.Title.Text = "Average Time to First Review"
	p.X.Label.Text = "Lines of Code Changed"
	p.Y.Label.Text = "Time to First Review (hours)"

	// Create a plotter.Values value and fill it with the times
	v := make(plotter.Values, len(*times))
	copy(v, *times)

	// Create a histogram of our values drawn
	h, err := plotter.NewHist(v, 16)
	if err != nil {
		panic(err)
	}
	h.Normalize(1)

	// Add the histogram to the plot
	p.Add(h)

	// Save the plot to a PNG file
	if err := p.Save(400, 400, "hist.png"); err != nil {
		panic(err)
	}
}
