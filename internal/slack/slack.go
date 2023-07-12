package slack

import (
	"context"
	"fmt"

	"github.com/slack-go/slack"
)

var (
	// Emojis available by defualt on Slack workspaces.
	EmojiApproved = "white_check_mark"

	// Error strings returned by the slack API,
	errAlreadyReacted = "already_reacted"
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
		return fmt.Errorf("GetConversationHistoryContext(): %w", err)
	}
	for _, msg := range resp.Messages {
		ref := slack.NewRefToMessage(channelID, msg.Timestamp)
		err := api.client.AddReactionContext(ctx, emoji, ref)
		if err != nil && err.Error() != errAlreadyReacted {
			return fmt.Errorf("AddReactionContext(ctx,%s,%+v): %w", emoji, ref, err)
		}
	}
	return nil
}

func (api API) UnreactChannel(ctx context.Context, channelID, url, emoji string, limit int) error {
	resp, err := api.client.GetConversationHistoryContext(ctx, &slack.GetConversationHistoryParameters{
		ChannelID: channelID,
		Limit:     limit,
	})
	if err != nil {
		return fmt.Errorf("GetConversationHistoryContext(): %w", err)
	}
	for _, msg := range resp.Messages {
		ref := slack.NewRefToMessage(channelID, msg.Timestamp)
		err := api.client.RemoveReactionContext(ctx, emoji, ref)
		if err != nil && err.Error() != errAlreadyReacted {
			return fmt.Errorf("AddReactionContext(ctx,%s,%+v): %w", emoji, ref, err)
		}
	}
	return nil
}

// React adds "emoji" reaction for the "count" most recent messages matching "url".
func (api API) ReactSearch(ctx context.Context, url, emoji string, count int) error {
	msgs, err := api.client.SearchMessagesContext(ctx, url, slack.SearchParameters{Count: count})
	if err != nil {
		return fmt.Errorf("SearchMessagesContext(): %w", err)
	}
	for _, msg := range msgs.Matches {
		ref := slack.NewRefToMessage(msg.Channel.ID, msg.Timestamp)
		err := api.client.AddReactionContext(ctx, emoji, ref)
		if err != nil && err.Error() != errAlreadyReacted {
			return fmt.Errorf("AddReactionContext(ctx,%s,%+v): %w", emoji, ref, err)
		}
	}
	return nil
}
