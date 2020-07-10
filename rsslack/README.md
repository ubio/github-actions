# RSS(lack)

This action sends a really simple slack notification to a given channel. Nothing fancy, just a simple message to a channel.

Usage:

```yaml
- name: Notify Slack
  uses: ubio/github-actions/slack-notifier@master
  with:
    channel: ${{ secrets.SLACK_CHANNEL }}
    message: "👋 from RSS(lack)"
    slack_token: ${{ secrets.SLACK_TOKEN }}
```
