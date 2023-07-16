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
	flag.BoolVar(&verbose, "v", false, "Verbose outputs the event payload.")
	flag.Parse()
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("%+v", err)
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

	emojis := map[string]bool{
		slack.EmojiApproved:         status.Approved,
		slack.EmojiChangesRequested: status.ChangesRequested,
		slack.EmojiCommented:        status.Commented,
		slack.EmojiClosed:           status.Closed,
		slack.EmojiMerged:           status.Merged,
	}
	slackAPI := slack.New(http.DefaultClient, slackBotToken)
	if err := slackAPI.SetEmojis(ctx, url, channelIDs, emojis); err != nil {
		return fmt.Errorf("unreact failed: %w", err)
	}
	return nil
}
