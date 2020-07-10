# Cert Checker

This action wraps [genkiroid/cert](https://github.com/genkiroid/cert) in a Dockerfile and outputs the result.


## Inputs

| Input | Required  | Default | Description
| ----- | --------- | ------- | -----------
| cmd   | `true`    |         | The command to run. See [genkiroid/cert](https://github.com/genkiroid/cert) for options


## Example Usage

Usage (build action):

```yaml
- name: Run
  uses: ubio/github-actions/cert-checker@master
  with:
    cmd: "-f json api.automationcloud.net"
```

Usage (optimised):

```yaml
- name: Run
  uses: docker://automationcloud/cert-checker:latest
  with:
    cmd: "-f json api.automationcloud.net"
```
