package slack

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode"

	"github.com/slack-go/slack"
)

var (
	// Emojis available by default on Slack workspaces.

	// Review state:
	// https://docs.github.com/en/rest/pulls/reviews?apiVersion=2022-11-28#list-reviews-for-a-pull-request
	EmojiApproved         = "white_check_mark"        // ✅
	EmojiChangesRequested = "x"                       // ❌
	EmojiCommented        = "speech_balloon"          // 💬
	EmojiReviewRequested  = "arrows_counterclockwise" // 🔄

	// PR state:
	// https://docs.github.com/en/rest/pulls/pulls?apiVersion=2022-11-28#get-a-pull-request
	EmojiClosed = "no_entry"            // ⛔
	EmojiMerged = "large_purple_square" // 🟪
)

type API struct {
	client              *slack.Client
	channelIDs          []string
	limit, repliesLimit int
}

func New(httpCLient *http.Client, token string, channelIDs []string, limit, repliesLimit int) API {
	return API{slack.New(token, slack.OptionHTTPClient(httpCLient)), channelIDs, limit, repliesLimit}
}

type SetEmoji struct {
	Name string
	Set  bool
}

func (api API) SetEmojis(ctx context.Context, match string, emojis []SetEmoji) error {
	for _, channelID := range api.channelIDs {
		resp, err := api.client.GetConversationHistoryContext(ctx, &slack.GetConversationHistoryParameters{
			ChannelID: channelID,
			Limit:     api.limit,
		})
		if err != nil {
			return fmt.Errorf("GetConversationHistoryContext(,%s): %w", channelID, err)
		}
		log.Printf("%d messages found in channel %s\n", len(resp.Messages), channelID)
		for _, msg := range resp.Messages {
			msgWithReplies := []slack.Message{msg}
			if msg.ReplyCount > 0 {
				replies, _, _, _ := api.client.GetConversationRepliesContext(ctx, &slack.GetConversationRepliesParameters{
					ChannelID: channelID,
					Timestamp: msg.Timestamp,
					Limit:     api.repliesLimit,
				})
				msgWithReplies = append(msgWithReplies, replies...)
				log.Printf("%d replies found in channel %s\n", len(replies), channelID)
			}
			for _, msg := range msgWithReplies {
				idx := strings.LastIndex(msg.Text, match)
				if idx == -1 {
					continue
				}
				rest := msg.Text[idx+len(match):]
				if len(rest) > 0 && unicode.IsDigit(rune(rest[0])) {
					continue
				}
				ref := slack.NewRefToMessage(channelID, msg.Timestamp)

				log.Printf("message %+v setting emojis %+v\n", ref, emojis)
				// Forcing consistent emoji order is not possible.
				for _, emoji := range emojis {
					if emoji.Set {
						if err := api.client.AddReactionContext(ctx, emoji.Name, ref); err != nil &&
							err.Error() != "already_reacted" {
							return fmt.Errorf("AddReactionContext(ctx,%s,%+v): %w", emoji.Name, ref, err)
						}
						continue
					}
					if err := api.client.RemoveReactionContext(ctx, emoji.Name, ref); err != nil &&
						err.Error() != "no_reaction" {
						return fmt.Errorf("RemoveReactionContext(ctx,%s,%+v): %w", emoji.Name, ref, err)
					}
				}
			}
		}
	}
	return nil
}
