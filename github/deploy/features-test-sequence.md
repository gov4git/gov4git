### set tally frequency to 60 seconds

### add user to community

```sh
gov4git --config testing-config.json user add --name petar --repo git@github.com:petar/gov4git-identity-public.git --branch main
```

### issue tokens to @petar via github UI

Create an issued, labelled `gov4git:directive` with body:

```
issue 3000 credits to @petar
```

### cast a vote on issue #1

```sh
gov4git --config testing-config.json ballot vote --name pmp/motion/priority_poll/1 --choices rank --strengths 10.0
```

After 1 minute, force cron.

Check voter account was debited:

```sh
gov4git --config testing-config.json account balance --id user:petar --asset plural
```

Check vote is reflected in issue #1:

```sh
gov4git --config testing-config.json ballot show --name pmp/motion/priority_poll/1
```

Check for an updated tally in the GitHub UI: https://github.com/gov4git/testing.project/issues/1
