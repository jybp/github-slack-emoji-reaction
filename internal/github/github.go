package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v53/github"
)

type API struct {
	client *github.Client
}

func New(httpCLient *http.Client, token string) API {
	return API{github.NewClient(httpCLient)}
}

type PullRequestStatus struct {
	State string
}

func (api API) PullRequestStatus(ctx context.Context, url string) (PullRequestStatus, error) {
	var owner, repo string
	var number int
	n, err := fmt.Sscanf(url, "https://github.com/%s/%s/pull/%d", &owner, &repo, &number)
	if err != nil {
		return PullRequestStatus{}, fmt.Errorf("failed to fmt.Sscanf(%s): %w", url, err)
	}
	if n != 3 {
		return PullRequestStatus{}, fmt.Errorf("unexpected number of items parsed: %d", n)
	}
	pr, _, err := api.client.PullRequests.Get(ctx, owner, repo, number)
	if err != nil {
		return PullRequestStatus{},
			fmt.Errorf("PullRequests.Get(,%s, %s, %d): %w", owner, repo, number, err)
	}
	return PullRequestStatus{State: pr.GetState()}, nil
}
