# Webhook

This action sends a message body to a URL using an (optional) secret for origin verification.

When the verification secret is supplied to the action, it will be set in a header called `X-ORIGIN-SECRET`.

Usage:

```yaml
- name: Notify Slack
  uses: ubio/github-actions/webhook@master
  with:
    channel: ${{ secrets.SLACK_CHANNEL }}
    message: "👋 from RSS(lack)"
    slack_token: ${{ secrets.SLACK_TOKEN }}
```
