# github-slack-emoji-reaction

Automatically add Slack emoji reactions to messages with Pull Request links.

![preview](docs/preview.gif?raw=true)

## Installation

1. Add the GitHub workflow below in your repository under `.github/workflows/github-slack-emoji-reaction.yml`.
  ```yaml
  name: GitHub Slack Emoji Reaction
  on:
    pull_request_review:
      types: [submitted, edited, dismissed]
    pull_request:
      types: [closed, reopened, review_requested]
  check_suite:
    types: [completed]
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
              ref: 'v1.3.0'
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
              SLACK_MESSAGES_LIMIT: 100
              SLACK_REPLIES_LIMIT: 100
  ```

2. Set a Slack bot token in the `SLACK_BOT_TOKEN` GitHub secret with the following permissions: `channels:history`, `reactions:read`, `reactions:write`.
<img src="docs/permissions.png?raw=true" alt="permissions" width="500">
<img src="docs/copy-token.png?raw=true" alt="copy-token" width="500">
<img src="docs/github1.png?raw=true" alt="github1" width="300">
<img src="docs/github2.png?raw=true" alt="github2" width="500">

3. Update `SLACK_CHANNEL_IDS` with the channel IDs to cover.
<img src="docs/channel-id.png?raw=true" alt="channel-id" width="500">

4. Optionally customize the Slack emojis with the `EMOJI_*` environment variables.
5. Optionally customize the number of messages looked up on each channel with the `SLACK_MESSAGES_LIMIT` and `SLACK_REPLIES_LIMIT` environment variables.
6. Invite the bot to the channels to cover in Slack.
<img src="docs/invite-app-bot.png?raw=true" alt="invite-app-bot" width="500">
