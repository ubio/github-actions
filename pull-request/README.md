# Pull Request

Creates a branch, commits the changes in `files` and then generates a pull request against a repository.

The PR author is set to the owner of the access token.

## Inputs

| Input                 | Required  | Default | Description
| --------------------- | --------- | ------- | -----------
| token                 | `true`    |         | A GitHub Personal Access Token which can access the target repo
| owner                 | `true`    |         | The owner of the repo to create the PR on (eg `ubio`)
| repository            | `true`    |         | The name of the repo to create the PR on (eg `my-repo`)
| message               | `true`    |         | The commit message
| files                 | `true`    |         | Comma-separated list of files to commit and their location. Example: `pull-request/README.md,pull-request/src/main.go`
| title                 | `true`    |         | The PR title
| body                  | `false`   |         | The body of the PR
| head                  | `true`    |         | Name of branch where changes are implemented
| base                  | `true`    |         | Name of branch where changes should be pulled into
| draft                 | `false`   | `false` | Whether the PR is in draft status
| maintainer_can_modify | `false`   | `true`  | Whether repo maintainers can modify the PR
| merge                 | `false`   | `false` | Whether to automatically merge the PR (if mergeable)

## Outputs

| Output  | Description
| ------- | ----------- 
| pr      | The URL to the created Pull Request
| merged  | true/false indicating whether the Pull Request was merged after creation

## Example Usage

Usage (build action):

```yaml
- name: Run
  uses: ubio/github-actions/pull-request@master
  with:
    token: ${{ secrets.ACCESS_TOKEN }}
    owner: "ubio"
    repository: "my-repo"
    message: "Update readme"
    files: "README.md"
    title: "My amazing PR"
    body: "These changes will change your life!"
    head: "my-feature"
    base: "master"
    draft: false
    maintainer_can_modify: true
```

Usage (optimised):

```yaml
- name: Run
  uses: docker://automationcloud/pull-request:latest
  with:
    token: ${{ secrets.ACCESS_TOKEN }}
    owner: "ubio"
    repository: "my-repo"
    message: "Update readme"
    files: "README.md"
    title: "My amazing PR"
    body: "These changes will change your life!"
    head: "my-feature"
    base: "master"
    draft: false
    maintainer_can_modify: true
```
