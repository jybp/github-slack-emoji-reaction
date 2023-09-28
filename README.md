# github-slack-emoji-reaction

1. Add the GitHub workflow below in your repository.
2. Set a Slack bot token in the `SLACK_BOT_TOKEN` secret with the following permissions: `channels:history`, `reactions:read`, `reactions:write`.
3. Update `SLACK_CHANNEL_IDS` with the channel IDs to cover.
4. Optionally customize the Slack emojis.
5. Invite the bot to the channels in Slack.

```
name: GitHub Slack Emoji Reaction
on:
  pull_request_review:
    types: [submitted]
  pull_request:
    types: [closed, reopened, review_requested]
permissions:
  pull-requests: read
jobs:
  gser:
    runs-on: ubuntu-latest
    continue-on-error: true
    steps:
        - uses: actions/checkout@v4
          with:
            repository: jybp/github-slack-emoji-reaction
            ref: 'v1.1.1'
        - uses: actions/setup-go@v4
          with:
            go-version: '1.20'
        - run: go run cmd/gser/main.go -v
          env:
            GITHUB_TOKEN: "${{ secrets.GITHUB_TOKEN }}"
            SLACK_BOT_TOKEN: "${{ secrets.SLACK_BOT_TOKEN }}"
            SLACK_CHANNEL_IDS: C05GGMQ2R61,C05GN394WCC
            EMOJI_APPROVED: white_check_mark
            EMOJI_CHANGES_REQUESTED: x
            EMOJI_COMMENTED: speech_balloon
            EMOJI_CLOSED: no_entry
            EMOJI_MERGED: large_purple_square
            EMOJI_REVIEW_REQUESTED: arrows_counterclockwise
```