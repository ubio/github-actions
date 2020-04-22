# Cert Checker

This action wraps [genkiroid/cert](https://github.com/genkiroid/cert) in a Dockerfile and outputs the result.

Usage:

```yaml
- name: Run
  uses: universalbasket/github-actions/cert-checker@master
  with:
    cmd: "-f json api.automationcloud.net"
```
