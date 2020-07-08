# Cert Notifier

This action takes the output of the cert checker and notifies slack of any upcoming renewals

Cert Checker usage:

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
