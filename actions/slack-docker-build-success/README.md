# Slack Docker Build Success

Sends a notification to slack for a successful Docker image build

## Inputs

| Input           | Required  | Default | Description
| --------------- | --------- | ------- | -----------
| slack_token     | `true`    |         | A slack token that can post the message
| channel         | `true`    |         | The channel to send the message to
| name            | `true`    |         | The name of the project to display in the message
| registry        | `true`    |         | The registry name - eg `registry.hub.docker.com`
| namespace       | `true`    |         | The registry namespace - eg `automationcloud`
| image           | `true`    |         | The image name - eg `my-project`
| tag             | `true`    |         | The image tag - eg `27.8.1-rc1`

## Usage:

This build step:

```yaml
- name: Notify Slack
  uses: docker://automationcloud/slack-docker-build-success:latest
  with:
    channel: ${{ secrets.SLACK_CHANNEL }}
    slack_token: ${{ secrets.SLACK_TOKEN }}
    name: My Project
    registry: registry.hub.docker.com
    namespace: automationcloud
    image: my-project
    tag: 27.8.1-rc1
```

Will generate the [following message in Slack](https://app.slack.com/block-kit-builder/T02FBD280#%7B%22blocks%22:%5B%7B%22type%22:%22section%22,%22text%22:%7B%22type%22:%22mrkdwn%22,%22text%22:%22:package:%20*My%20Project*%20has%20been%20built%20and%20pushed%20to%20%60registry.hub.docker.com/automationcloud/my-project:27.8.1-rc1%60%22%7D%7D%5D%7D):

```markdown
:package: *My Project* has been built and pushed to `registry.hub.docker.com/automationcloud/my-project:27.8.1-rc1`
```
