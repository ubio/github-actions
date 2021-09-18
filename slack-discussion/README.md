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

Will generate the [following message in Slack](https://app.slack.com/block-kit-builder/T02FBD280#%7B%22blocks%22:%5B%7B%22type%22:%22section%22,%22text%22:%7B%22text%22:%22:speech_balloon:%20General%20-%20*Proxies%20Squad*%22,%22type%22:%22mrkdwn%22%7D%7D,%7B%22type%22:%22header%22,%22text%22:%7B%22type%22:%22plain_text%22,%22text%22:%22Ola%22%7D%7D,%7B%22type%22:%22section%22,%22text%22:%7B%22type%22:%22mrkdwn%22,%22text%22:%22hi%20friends%22%7D%7D,%7B%22type%22:%22divider%22%7D,%7B%22type%22:%22context%22,%22elements%22:%5B%7B%22type%22:%22image%22,%22image_url%22:%22https://avatars1.githubusercontent.com/u/21031067?v=4%22,%22alt_text%22:%22cute%20cat%22%7D,%7B%22type%22:%22mrkdwn%22,%22text%22:%22Author:%20%3Chttps://github.com/andrew-waters%7C@andrew-waters%3E%22%7D,%7B%22type%22:%22mrkdwn%22,%22text%22:%22November%2012,%202019%20@%2022:37%22%7D,%7B%22type%22:%22mrkdwn%22,%22text%22:%22%3Chttps://github.com/andrew-waters%7CView%20Discussion%3E%22%7D%5D%7D%5D%7D) depending on the discussion event passed by GitHub.

References:

 - [Discussion Payload](https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#discussion)
 