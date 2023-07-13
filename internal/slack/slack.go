package slack

import (
	"context"
	"fmt"
	"strings"

	"github.com/slack-go/slack"
)

var (
	// Emojis available by default on Slack workspaces.

	// Review state:
	// https://docs.github.com/en/rest/pulls/reviews?apiVersion=2022-11-28#list-reviews-for-a-pull-request
	EmojiApprove        = "white_check_mark" // ✅ APPROVE
	EmojiComment        = "speech_balloon"   // 💬 COMMENT
	EmojiRequestChanges = "x"                // ❌ REQUEST_CHANGES

	// PR state:
	// https://docs.github.com/en/rest/pulls/pulls?apiVersion=2022-11-28#get-a-pull-request
	EmojiDraft  = ""
	EmojiOpen   = ""
	EmojiMerged = "larged_purple_square" // 🟪
	EmojiClosed = "no_entry"             // ⛔

	// Error strings returned by the slack API,
	errAlreadyReacted = "already_reacted"
	errNoReaction     = "no_reaction"
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
		if !strings.Contains(msg.Text, url) {
			continue
		}
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
		if !strings.Contains(msg.Text, url) {
			continue
		}
		ref := slack.NewRefToMessage(channelID, msg.Timestamp)
		err := api.client.RemoveReactionContext(ctx, emoji, ref)
		if err != nil && err.Error() != errNoReaction {
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
