
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

## Prepare a GitHub user and access token for automation

The governance system is deployed in the form of a GitHub action, installed in a newly created governance repo in the same GitHub organization as your project repo.

The governance automation — invoked by the GitHub action — executes on behalf of a dedicated GitHub user, which represents the governance system itself. You must **create a new GitHub user, designated as the governance automation user** — and name it appropriately, as it will speak to the community users via GitHub comments.

For instance, the _Plurality Book Project_ — *@pluralitybook* on GitHub — uses a dedicated user called _Plurality Book DAO_ — *@pluralitybook-dao* on GitHub — to operate the governance system.

**Invite the automation user to your organization with Owner privileges.**

**Create a [GitHub access token](https://docs.github.com/en/github/authenticating-to-github/creating-a-personal-access-token) for the automation user.**

Your token must have permissions to:
- create new repositories in your GitHub organization
- manage GitHub Actions within these repositories (i.e. environments, secrets, variables, workflows)
- read the issues and pull requests of your project repository

To create the token:
- Go to your GitHub profile settings
- Click on "Developer settings"
- Click on "Personal access tokens"
- Click on "Fine-grained tokens"
- Click on "Generate new token"
- Choose a name, description and expiration for your token
- Under "Resource owner" pick the user or organization that owns the project repository
- Under "Repository access" select "All repositories"
- Under "Repository permissions" make the following choices:

     | category | access |
     | ----------- | ----------- |
     | actions | read-write |
     | admin | read-write |
     | contents | read-write |
     | environments | read-write |
     | issues | read-write |
     | pull-requests | read-write |
     | secrets | read-write |
     | variables | read-write |
     | workflows | read-write |

- Under "Organization permissions" make the following choices:

     | category | access |
     | ----------- | ----------- |
     | members | read-only |

- Click on "Generate token" and write down the generated token


## Deploy governance for your project repository

You are now ready to deploy governance on your project repository. This can be accomplished with a single command:

```bash
gov4git -v github deploy \
     --token=$YOUR_ACCESS_TOKEN \
     --project=$PROJECT_OWNER/$PROJECT_REPO \
     --release=$GOV4GIT_RELEASE
```

Here `$GOV4GIT_RELEASE` specifies the gov4git release on GitHub that you want to use for the deployment.

## What does the deployment command do?

During the deployment, the following steps are performed:

- Two new repositories — one public, one private — are created within the GitHub organization of your project repository. The public repository is named `$PROJECT_REPO-gov.public` and the private repository is named `$PROJECT_REPO-gov.private`.

- Both repositories are initialized with a newly-generated identity for your governance system. This step corresponds to the `gov4git init-gov` command.

- One GitHub actions are created in the public governance repository, named `.github/workflows/gov4git_cron.yml`. This action is accompanied by a helper script `.github/scripts/gov4git_cron.sh`. The action runs every two minutes. It is responsible for:
     - Reading all issues and pull requests from your project repositories and updating the governance system accordingly, and
     - Fetching votes and other service requests by your community members and incorporating them into the governance system.

- A new GitHub environment called `gov4git:governance` is created, where the GitHub action `gov4git_cron.yml` runs. This environment contains a set of variables:
  - `GOV4GIT_RELEASE` is the gov4git release to use for the automation
  - `GOV_PUBLIC_REPO_URL` is the HTTPS URL of the public governance repository
  - `GOV_PRIVATE_REPO_URL` is the HTTPS URL of the private governance repository
  - `PROJECT_OWNER` is the GitHub user or organization owning your project repository
  - `PROJECT_REPO` is the name of your project repository
  - `SYNC_GITHUB_FREQ` is the number of seconds between updates from GitHub.
  - `SYNC_COMMUNITY_FREQ`is the number of seconds between updates from community members.
  - `SYNC_FETCH_PAR` is the number of parallel repository fetches performed during updates from community members.

     Additionally, a GitHub secret called `ORGANIZER_GITHUB_TOKEN` is created in the public governance repository. This secret contains the GitHub access token you provided to the deployment command. It is used by the GitHub actions to access your project repository, as well as the governance repositories.

## Managing your deployment

All governance write operations are strictly contained within the public and private governance repositories. In particular, governance does not write to your project repository.

The entire state of governance is captured in the most recent commit of the public and private governance repositories.

This allows you to perform a few basic administrative tasks, using standard git and GitHub operations:

- _If you want to stop and erase a deployment_, delete the public and private governance repositories.

- _If you want to stop a deployment but keep the governance state_, edit the GitHub actions and comment out the cron triggers.

- _If you want to reduce the size of your governance repositories_, archive their `main`-branch history in another branch or repository, and reset the `main` branch to contain only the most recent commit.

- _If you want to upgrade your deployment to a newer version of gov4git_, edit the GitHub environment and change the `GOV4GIT_RELEASE` variable to the new release.
