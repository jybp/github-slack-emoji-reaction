package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jybp/github-slack-emoji-reaction/internal/github"
	"github.com/jybp/github-slack-emoji-reaction/internal/slack"
	"golang.org/x/oauth2"
)

var (
	verbose bool
)

func init() {
	log.SetFlags(0)
	flag.BoolVar(&verbose, "v", false, "verbose mode")
	flag.Parse()
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("%+v\n", err)
		os.Exit(1)
	}
}

func run() error {
	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	if len(slackBotToken) == 0 {
		return errors.New("SLACK_BOT_TOKEN not set")
	}
	githubToken := os.Getenv("GITHUB_TOKEN")
	if len(githubToken) == 0 {
		return errors.New("GITHUB_TOKEN not set")
	}
	channelIDs := strings.Split(os.Getenv("SLACK_CHANNEL_IDS"), ",")
	if len(channelIDs) == 0 {
		return errors.New("SLACK_CHANNEL_IDS not set")
	}
	ghEventPath := os.Getenv("GITHUB_EVENT_PATH")
	if len(ghEventPath) == 0 {
		return errors.New("GITHUB_EVENT_PATH not set")
	}
	ghEvent, err := os.ReadFile(ghEventPath)
	if err != nil {
		return fmt.Errorf("could not read %s: %w", ghEventPath, err)
	}

	if !verbose {
		log.SetOutput(ioutil.Discard)
	}
	log.Printf("%s:\n%s\n", ghEventPath, string(ghEvent))

	url, owner, repo, number, err := github.ParsePayload(ghEvent)
	if err != nil {
		return fmt.Errorf("could not parse %s: %w", ghEventPath, err)
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	githubAPI := github.New(tc)
	status, err := githubAPI.PullRequestStatus(ctx, owner, repo, number)
	if err != nil {
		return err
	}

	log.Printf("status for %s: %+v\n", url, status)

	if len(os.Getenv("EMOJI_APPROVED")) > 0 {
		slack.EmojiApproved = os.Getenv("EMOJI_APPROVED")
	}
	if len(os.Getenv("EMOJI_CHANGES_REQUESTED")) > 0 {
		slack.EmojiChangesRequested = os.Getenv("EMOJI_CHANGES_REQUESTED")
	}
	if len(os.Getenv("EMOJI_COMMENTED")) > 0 {
		slack.EmojiCommented = os.Getenv("EMOJI_COMMENTED")
	}
	if len(os.Getenv("EMOJI_CLOSED")) > 0 {
		slack.EmojiClosed = os.Getenv("EMOJI_CLOSED")
	}
	if len(os.Getenv("EMOJI_MERGED")) > 0 {
		slack.EmojiMerged = os.Getenv("EMOJI_MERGED")
	}
	if len(os.Getenv("EMOJI_REVIEW_REQUESTED")) > 0 {
		slack.EmojiReviewRequested = os.Getenv("EMOJI_REVIEW_REQUESTED")
	}
	limit, repliesLimit := 40, 20
	if len(os.Getenv("SLACK_MESSAGES_LIMIT")) > 0 {
		limit, err = strconv.Atoi(os.Getenv("SLACK_MESSAGES_LIMIT"))
		if err != nil {
			return fmt.Errorf("could not parse SLACK_MESSAGES_LIMIT: %w", err)
		}
	}
	if len(os.Getenv("SLACK_REPLIES_LIMIT")) > 0 {
		repliesLimit, err = strconv.Atoi(os.Getenv("SLACK_REPLIES_LIMIT"))
		if err != nil {
			return fmt.Errorf("could not parse SLACK_REPLIES_LIMIT: %w", err)
		}
	}
	emojis := []slack.SetEmoji{
		{Name: slack.EmojiCommented, Set: status.Commented},
		{Name: slack.EmojiReviewRequested, Set: status.ReviewRequested},
		{Name: slack.EmojiChangesRequested, Set: status.ChangesRequested},
		{Name: slack.EmojiApproved, Set: status.Approved},
		{Name: slack.EmojiMerged, Set: status.Merged},
		{Name: slack.EmojiClosed, Set: status.Closed},
	}
	slackAPI := slack.New(http.DefaultClient, slackBotToken, channelIDs, limit, repliesLimit)
	if err := slackAPI.SetEmojis(ctx, url, emojis); err != nil {
		return fmt.Errorf("SetEmojis failed: %w", err)
	}
	return nil
}
