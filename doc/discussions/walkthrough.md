## Initial setup
### Prepare your workstation

Install gov4git.

Create a config file `~/.gov4git/config.json` on your workstation with contents

```
{
     "public_url": PUBLIC_REPO,
     "private_url": PRIVATE_REPO,
     "community_url": COMMUNITY_REPO,
     "community_branch": COMMUNITY_BRANCH
}
```

Create empty repos `PUBLIC_REPO` and `PRIVATE_REPO` on your git provider (e.g. GitHub).

`COMMUNITY_REPO` and `COMMUNITY_BRANCH` are the community repo and the name of the main branch (usually `master` or `main`), respectively.

### Generate your identity

Initialize your (personal) gov4git identity:
```sh
gov4git init
```
This creates new public and private keys and stored them in `PUBLIC_REPO` and `PRIVATE_REPO`, respectively.

## Manage community users

### Add a user to the community

```sh
gov4git gov user add --name=petar --url=git@github.com:petar/gov4git.public.git
```

- We could support import from CSV to facilitate Excel users.

### List community users

```sh
gov4git gov user list
```

- Could support a CSV printout to facilitate dumping into Excel?

### Set/get user properties

Users can be associated with arbitrary key/value properties. This enables associating users with application-specific information like balances, badges, SBTs, etc.

Set the user property named "badges" to some JSON representing two SBTs:

```sh
gov4git gov user set --name=petar --key=badges --value='["SBT1", "SBT2"]'
```

Get the user property named "badges":

```sh
gov4git gov user get --name=petar --key=badges
```

## Manage community groups and memberships

Groups are sets of users that support dynamic membership changes.
They can model things like "contributor", "translator", "admin", etc.
Groups can be used in many ways. For instance:
- a ballot targets a specific group
- user's group membership can be used as input to voting algorithms

Create a new empty group (of users):

```sh
gov4git gov group add --name=everyone
```

Add user `petar` to group `everyone`:

```sh
gov4git gov member add --group=everyone --user=petar
```

## Ballots (polls, referendums, etc.)

A ballot is any procedure that draws a conclusion from the inputs of a group of users.
There can be many types of ballots with varying purposes (e.g. polling, deciding) and tallying strategies (e.g. quadratic voting, 1p1v, etc). Furthermore, users can easily implement and plug new ballot strategies.

### Creating a new ballot (e.g. poll)

Create a new poll, named `my_first_ballot`, on the `main` branch of the community repo, which targets all users and uses the `priority-poll` strategy:

```sh
gov4git gov ballot --govern-branch=main --group=everyone --path=my_first_ballot --strategy=priority-poll
```

### Vote on a ballot

A user can cast their vote by identifying the repo, branch and name of the ballot they want to vote on, as well their choice and strength elections. Here the user is voting on ballot named `my_first_ballot`, located on the `main` branch of the community repo (configured in the user's gov4git config file). They are casting a vote for choice `https://github.com/gov4git/gov4git/issues/1` of strength `3`:

```sh
gov4git gov vote --ballot-branch=main --ballot-path=my_first_ballot --choice=https://github.com/gov4git/gov4git/issues/1 --strength=3
```

### Tally the results



```sh
gov4git gov tally --ballot-branch=main --ballot-path=my_first_ballot
```
