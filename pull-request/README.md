# Pull Request

Creates a pull request against the current repository.

## Inputs

| Input                 | Required  | Default | Description
| --------------------- | --------- | ------- | -----------
| token                 | `true`    |         | A GitHub Personal Access Token which can access the target repo
| owner                 | `true`    |         | The owner of the repo to create the PR on (eg `universalbasket`)
| repository            | `true`    |         | The name of the repo to create the PR on (eg `my-repo`)
| title                 | `true`    |         | The PR title
| body                  | `false`   |         | The body of the PR
| head                  | `true`    |         | Name of branch where changes are implemented
| base                  | `true`    |         | Name of branch where changes should be pulled into
| draft                 | `false`   | `false` | Whether the PR is in draft status
| maintainer_can_modify | `false`   | `true`  | Whether repo maintainers can modify the PR

## Example Usage

```yaml
- name: Run
  uses: universalbasket/github-actions/pull-request@master
  with:
    token: ${{ secrets.ACCESS_TOKEN }}
    owner: "universalbasket"
    repository: "my-repo"
    title: "My amazing PR"
    body: "These changes will change your life!"
    head: "my-feature"
    base: "master"
    draft: false
    maintainer_can_modify: true
```
