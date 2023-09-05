
# Deployment guide

This guide will help you deploy governance for a GitHub project repository.

## Prerequisites



## Access token fine-grain permissions

### Repository permissions

| category | access |
| ----------- | ----------- |
| actions | read-write |
| admin | read-write |
| contents | read-write |
| environments | read-write |
| issues | read-only |
| meta | read-only |
| pull-requests | read-only |
| secrets | read-write |
| variables | read-write |
| workflows | read-write |

### Organization permissions

n/a



## Create a governance environment

Add a new environment to your GitHub project repository, named `governance`.

Using the GitHub UI, add an environment variable `GOV4GIT_RELEASE` pointing to the desired release of gov4git. For instance,

```GOV4GIT_RELEASE=v1.1.4```

