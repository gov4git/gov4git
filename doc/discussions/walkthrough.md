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

- We could support a CSV printout to facilitate dumping into Excel.

### Set/get user properties

Users can be associated with arbitrary key/value properties. This enables associating users with application-specific information like badges, SBTs, etc.

```sh
gov4git gov user set --name=petar --key=badges --value='{"some":"SBT"}'
```

```sh
gov4git gov user get --name=petar --key=badges
```

## Manage community groups and memberships

Create a new empty group (of users):

```sh
gov4git gov group add --name=everyone
```

Add user `petar` to group `everyone`:

```sh
gov4git gov member add --group=everyone --user=petar
```

## Ballots (polls, referendums, etc.)

### Creating a new ballot (e.g. poll)

```sh
gov4git gov ballot --govern-branch=main --group=everyone --path=my_first_ballot --strategy=priority-poll
```

### Vote on a ballot

```sh
gov4git gov vote --ballot-branch=main --ballot-path=my_first_ballot --choice=https://github.com/gov4git/gov4git/issues/1 --strength=3
```

### Tally the results

```sh
gov4git gov tally --ballot-branch=main --ballot-path=my_first_ballot
```
