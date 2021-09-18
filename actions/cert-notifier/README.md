# Cert Notifier

This action takes the output of the cert checker and notifies slack of any upcoming renewals

## Inputs

| Input           | Required  | Default | Description
| --------------- | --------- | ------- | -----------
| warn_under_days | `true`    |         | The threshold for sending warnings to Slack
| channel         | `true`    |         | The channel to send warnings to
| certs           | `true`    |         | The certs to check and notifyy on - see the [cert-checker](../cert-checker)
| slack_token     | `true`    |         | A slack token that can post the message


## Example Usage

Usage (build action):

```yaml
- name: Check
  id: check
  uses: ubio/github-actions/cert-checker@master

- name: Notify
  uses: ubio/github-actions/cert-notifier@master
  with:
    warn_under_days: 30
    channel: "#general"
    certs: ${{ steps.check.outputs.result }}
    slack_token: ${{ secrets.SLACK_TOKEN }}
```

Usage (optimised):

```yaml
- name: Check
  id: check
  uses: docker://automationcloud/cert-checker:latest

- name: Notify
  uses: docker://automationcloud/cert-notifier:latest
  with:
    warn_under_days: 30
    channel: "#general"
    certs: ${{ steps.check.outputs.result }}
    slack_token: ${{ secrets.SLACK_TOKEN }}
```
