# Slack Discussion

Sends a notification to slack for a new discussion topic in a squad repo.

## Inputs

| Input           | Required  | Default | Description
| --------------- | --------- | ------- | -----------
| slack_token     | `true`    |         | A slack token that can post the message
| channel         | `true`    |         | The channel to send the message to

## Usage:

This build step:

```yaml
- name: Notify Slack
  uses: docker://ubio/slack-discussion:latest
  with:
    channel: tldr
    slack_token: ${{ secrets.SLACK_TOKEN }}
```

Will generate the [following message in Slack](https://app.slack.com/block-kit-builder/T02FBD280#%7B%22blocks%22:%5B%7B%22type%22:%22section%22,%22text%22:%7B%22type%22:%22mrkdwn%22,%22text%22:%22:package:%20*My%20Project*%20has%20been%20built%20and%20pushed%20to%20%60registry.hub.docker.com/automationcloud/my-project:27.8.1-rc1%60%22%7D%7D%5D%7D):

```markdown
:mega: *My Project* has been built and pushed to `registry.hub.docker.com/automationcloud/my-project:27.8.1-rc1`
```



https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#discussion_comment