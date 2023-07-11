# github-slack-emoji-reaction


## Installation

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
