package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-github/v53/github"
)

// pull_request / pull_request_review event payload.
type actionPayload struct {
	Action      string `json:"action"` // "closed", "submitted"
	PullRequest struct {
		Links struct {
			HTML struct {
				Href string `json:"href"` // "https://github.com/slack-emoji-reaction/public-repo/pull/1"
			} `json:"html"`
		} `json:"_links"`
		Base struct {
			Repo struct {
				Name  string `json:"name"` // "public-repo"
				Owner struct {
					HTMLURL string `json:"html_url"` // "https://github.com/slack-emoji-reaction"
					Login   string `json:"login"`
				} `json:"owner"`
			} `json:"repo"`
		} `json:"base"`
		// Draft   bool   `json:"draft"`
		HTMLURL string `json:"html_url"` // "https://github.com/slack-emoji-reaction/public-repo/pull/1"
		Number  int    `json:"number"`   // 1
		// RequestedReviewers []any  `json:"requested_reviewers"`
		// RequestedTeams     []any  `json:"requested_teams"`
		State string `json:"state"` // "open"
	} `json:"pull_request"`
	Repository struct {
		FullName string `json:"full_name"` // "slack-emoji-reaction/public-repo"
		Name     string `json:"name"`      // "public-repo"
		Owner    struct {
			Login string `json:"login"` // "slack-emoji-reaction"
		} `json:"owner"`
	} `json:"repository"`

	// Only present for "pull_request_review". "action" is "submitted".
	// Review struct {
	// 	State string `json:"state"` // "commented"
	// } `json:"review"`
}

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
