
# Deployment guide

This guide will help you deploy governance for a GitHub project repository.

## Install gov4git on your local machine

Make sure you have the [Go language installed](https://golang.org/doc/install) on your local machine.

Install `gov4git` on your local machine:

```bash
go get github.com/gov4git/gov4git/gov4git@latest
```

Verify `gov4git` is installed:

```bash
gov4git version
```

## Prepare a GitHub access token for automation

The governance system is deployed in the form of GitHub actions, installed in a newly created governance repo in the same GitHub organization as your project repo.

First, you must create a [GitHub access token](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token) which will be used by the automation logic.

Your token must have permissions to:
- create new repositories in your GitHub organization
- manage GitHub Actions within these repositories (i.e. environments, secrets, variables, workflows)
- read the issues and pull requests of your project repository

XXX

#### Repository permissions

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

#### Organization permissions

None.



## Deploy governance for your project repository

XXX

## How does governance integration with GitHub work?

XXX
