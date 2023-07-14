package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jybp/github-slack-emoji-reaction/internal/slack"
)

var (
	unreact bool
)

func init() {
	log.SetFlags(0)
	flag.BoolVar(&unreact, "unreact", false, "Unreacts instead of reacts.")
	flag.Parse()
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
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

	fmt.Printf("%s:\n%s\n", ghEventPath, string(ghEvent))
	ctx := context.Background()

	_ = githubToken
	// githubAPI := github.New(http.DefaultClient, githubToken)
	// status, err := githubAPI.PullRequestStatus(ctx, url)
	// if err != nil {
	// 	return err
	// }
	// _ = status

	slackAPI := slack.New(http.DefaultClient, slackBotToken)
	for _, channelID := range channelIDs {
		if unreact {
			if err := slackAPI.UnreactChannel(ctx, channelID, "a", slack.EmojiApprove, 100); err != nil {
				return fmt.Errorf("unreact failed: %w", err)
			}
			continue
		}
		if err := slackAPI.ReactChannel(ctx, channelID, "a", slack.EmojiApprove, 100); err != nil {
			return fmt.Errorf("react failed: %w", err)
		}
	}
	return nil
}
