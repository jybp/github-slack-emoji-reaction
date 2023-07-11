package slack

import (
	"context"
	"errors"
	"fmt"

	"github.com/slack-go/slack"
)

var (
	// Emojis available by defualt on Slack workspaces.
	EmojiApproved = "white_check_mark"
)

type API struct {
	client *slack.Client
}

func New(token string) API {
	return API{slack.New(token)}
}

func (api API) ReactChannel(ctx context.Context, channelID, url, emoji string, limit int) error {
	resp, err := api.client.GetConversationHistoryContext(ctx, &slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     limit,
	})
	if err != nil {
		return err
	}
	if len(resp.Error) > 0 {
		return fmt.Errorf("conversations.history: %s", resp.Error)
	}
	if !resp.Ok {
		return errors.New("conversations.history: ok is false")
	}
	for _, msg := range resp.Messages {
		if err := api.client.AddReactionContext(ctx, emoji,
			slack.NewRefToMessage(msg.Channel, msg.Timestamp)); err != nil {
			return err
		}
	}
	return nil
}

// React adds "emoji" reaction for the "count" most recent messages matching "url".
func (api API) ReactSearch(ctx context.Context, url, emoji string, count int) error {
	msgs, err := api.client.SearchMessagesContext(ctx, url, slack.SearchParameters{Count: count})
	if err != nil {
		return err
	}
	for _, msg := range msgs.Matches {
		if err := api.client.AddReactionContext(ctx, emoji,
			slack.NewRefToMessage(msg.Channel.ID, msg.Timestamp)); err != nil {
			return err
		}
	}
	return nil
}
