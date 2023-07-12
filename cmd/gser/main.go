package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"

	"github.com/jybp/github-slack-emoji-reaction/internal/slack"
)

var (
	url       string
	channelID string
	unreact   bool
)

func init() {
	log.SetFlags(0)
	flag.StringVar(&url, "url", "", `The GitHub Pull Request URL to add an amoji reaction to.`)
	flag.StringVar(&channelID, "channel", "", "The Slack channel ID to search 'url' into.")
	flag.BoolVar(&unreact, "unreact", false, "Unreacts instead of reacts.")
	flag.Parse()
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run() error {
	slackToken := os.Getenv("SLACK_TOKEN")
	if len(slackToken) == 0 {
		return errors.New("SLACK_TOKEN not set")
	}
	if len(url) == 0 || len(channelID) == 0 {
		flag.PrintDefaults()
		return nil
	}
	ctx := context.Background()
	api := slack.New(slackToken)
	if unreact {
		return api.UnreactChannel(ctx, channelID, url, slack.EmojiApproved, 100)
	}
	return api.ReactChannel(ctx, channelID, url, slack.EmojiApproved, 100)
}
