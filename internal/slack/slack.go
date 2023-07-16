package slack

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/slack-go/slack"
)

var (
	// Emojis available by default on Slack workspaces.

	// Review state:
	// https://docs.github.com/en/rest/pulls/reviews?apiVersion=2022-11-28#list-reviews-for-a-pull-request
	EmojiApproved         = "white_check_mark" // ✅
	EmojiCommented        = "speech_balloon"   // 💬
	EmojiChangesRequested = "x"                // ❌

	// PR state:
	// https://docs.github.com/en/rest/pulls/pulls?apiVersion=2022-11-28#get-a-pull-request
	EmojiMerged = "larged_purple_square" // 🟪
	EmojiClosed = "no_entry"             // ⛔
)

type API struct {
	client *slack.Client
}

func New(httpCLient *http.Client, token string) API {
	return API{slack.New(token, slack.OptionHTTPClient(httpCLient))}
}

func (api API) SetEmojis(ctx context.Context, match string, channelIDs []string, emojis map[string]bool) error {
	for _, channelID := range channelIDs {
		resp, err := api.client.GetConversationHistoryContext(ctx, &slack.GetConversationHistoryParameters{
			ChannelID: channelID,
			Limit:     100,
		})
		if err != nil {
			return fmt.Errorf("GetConversationHistoryContext(,%s): %w", channelID, err)
		}
		for _, msg := range resp.Messages {
			if !strings.Contains(msg.Text, match) {
				continue
			}
			ref := slack.NewRefToMessage(channelID, msg.Timestamp)
			for emoji, set := range emojis {
				if set {
					if err := api.client.AddReactionContext(ctx, emoji, ref); err != nil &&
						err.Error() != "already_reacted" {
						return fmt.Errorf("AddReactionContext(ctx,%s,%+v): %w", emoji, ref, err)
					}
					continue
				}
				if err := api.client.RemoveReactionContext(ctx, emoji, ref); err != nil &&
					err.Error() != "no_reaction" {
					return fmt.Errorf("RemoveReactionContext(ctx,%s,%+v): %w", emoji, ref, err)
				}
			}
		}
	}
	return nil
}
