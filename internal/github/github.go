package github

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/go-github/v53/github"
)

// pull_request / pull_request_review event payload.
type eventPayload struct {
	PullRequest struct {
		HTMLURL string `json:"html_url"` // "https://github.com/slack-emoji-reaction/public-repo/pull/1"
		Number  int    `json:"number"`   // 1
	} `json:"pull_request"`
	Repository struct {
		Name  string `json:"name"` // "public-repo"
		Owner struct {
			Login string `json:"login"` // "slack-emoji-reaction"
		} `json:"owner"`
	} `json:"repository"`
}

func ParsePayload(payload []byte) (prUrl, owner, repo string, number int, _ error) {
	var event eventPayload
	if err := json.Unmarshal(payload, &event); err != nil {
		return "", "", "", 0, fmt.Errorf("could not unmarshal payload: %w", err)
	}
	return event.PullRequest.HTMLURL,
		event.Repository.Owner.Login,
		event.Repository.Name,
		event.PullRequest.Number, nil
}

type API struct {
	client *github.Client
}

func New(httpCLient *http.Client) API {
	return API{github.NewClient(httpCLient)}
}

type PullRequestStatus struct {
	Approved         bool
	ChangesRequested bool
	Commented        bool
	Closed           bool
	Merged           bool

	// True when a reviewer who has already submitted a review has been re-requested to review the PR.
	ReviewRequested bool
}

func (api API) PullRequestStatus(ctx context.Context, owner, repo string, number int) (PullRequestStatus, error) {
	pr, _, err := api.client.PullRequests.Get(ctx, owner, repo, number)
	if err != nil {
		return PullRequestStatus{},
			fmt.Errorf("PullRequests.Get(,%s, %s, %d): %w", owner, repo, number, err)
	}
	log.Printf("PR %s %s %d: %s\n", owner, repo, number, pr.String())
	status := PullRequestStatus{
		Merged: pr.GetMerged(),
	}
	if !status.Merged && pr.GetState() == "closed" {
		status.Closed = true
	}
	reviews, _, err := api.client.PullRequests.ListReviews(ctx, owner, repo, number, nil)
	if err != nil {
		return PullRequestStatus{},
			fmt.Errorf("PullRequests.ListReviews(,%s, %s, %d): %w", owner, repo, number, err)
	}
	latestByAuthor := map[int64]string{}
	for i, review := range reviews {
		log.Printf("review %d: %s\n", i, review.String())
		latestByAuthor[review.User.GetID()] = strings.ToLower(review.GetState())
	}
	// https://docs.github.com/en/rest/pulls/review-requests?apiVersion=2022-11-28#get-all-requested-reviewers-for-a-pull-request
	// Once a requested reviewer submits a review, they are no longer considered a requested reviewer.
	// Their review will instead be returned by the List reviews for a pull request operation.
	reviewers, _, err := api.client.PullRequests.ListReviewers(ctx, owner, repo, number, nil)
	if err != nil {
		return PullRequestStatus{},
			fmt.Errorf("PullRequests.ListReviewers(,%s, %s, %d): %w", owner, repo, number, err)
	}
	// If one reviewer has submitted a review AND is still in the list of requested reviewers, it means
	// that he has been requested to review the PR again by the PR author and his old review should be
	// considered dismissed.
	for _, reviewer := range reviewers.Users {
		closedOrMerged := status.Closed || status.Merged
		_, ok := latestByAuthor[reviewer.GetID()]
		if ok && !closedOrMerged {
			// Only mark the PR has re-requesting a review if it's not closed/merged.
			status.ReviewRequested = true
		}
		delete(latestByAuthor, reviewer.GetID())
	}
	for _, state := range latestByAuthor {
		switch state {
		case "approved":
			status.Approved = true
		case "changes_requested":
			status.ChangesRequested = true
		case "commented":
			status.Commented = true
		}
	}
	return status, nil
}
