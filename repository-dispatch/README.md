# Repository Dispatch

Creates a [repository dispatch](https://help.github.com/en/actions/reference/events-that-trigger-workflows#external-events-repository_dispatch) event.

This is a custom event you can describe which can in turn be used to trigger a workflow on the target repo.

## Inputs

| Input       | Required  | Default | Description
| ----------- | --------- | ------- | -----------
| token       | `true`    |         | A GitHub Personal Access Token which can access the target repo
| owner       | `true`    |         | The owner of the repo to send the dispatch event to (eg `universalbasket`)
| repository  | `true`    |         | The name of the repo to send the dispatch event to (eg `my-repo`)
| event       | `true`    |         | The event type
| payload     | `false`   | `{}`    | JSON payload with data that your target action or worklow may use

## Example Usage

```yaml
- name: Run
  uses: universalbasket/github-actions/repository-dispatch@master
  with:
    token: ${{ secrets.ACCESS_TOKEN }}
    owner: "github-owner"
    repository: "repo-name"
    event: "your-event"
    payload: '{"extra":"info"}'
```
