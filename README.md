# github-slack-emoji-reaction

## Slack token

Use a token with the following permissions: "channels:history", "reactions:read", "reactions:write".

1. "Create an App" at https://api.slack.com/apps.
2. "From an app manifest".
3. Choose the Slack workspace.
4. Copy paste the following manifest:
```json
{
    "display_information": {
        "name": "GSER"
    },
    "features": {
		"bot_user": {
			"display_name": "GSER"
		}
	},
    "settings": {
        "org_deploy_enabled": false,
        "socket_mode_enabled": false,
        "is_hosted": false,
        "token_rotation_enabled": false
    },
    "oauth_config": {
		"scopes": {
			"bot": [
				"channels:history",
				"reactions:read",
				"reactions:write"
			]
		}
	}
}
```
5. "Install to Workspace" under Settings / Basic Information.
6. "Allow".
7. Copy "Bot User OAuth Token" under Features / OAuth & Permissions.

## GitHub token

Use a token with the following permission: "Pull requests" (Read-only).

1. Settings.
2. Developer settings.
3. Personal access tokens.
4. Select "Pull requests" (Read-only).
