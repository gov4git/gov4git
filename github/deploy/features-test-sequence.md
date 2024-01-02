### set tally frequency to 60 seconds

### add user to community

```sh
gov4git --config testing-config.json user add --name petar --repo https://github.com/petar/gov4git-identity-public.git --branch main
```

### issue tokens to @petar via github UI

Create an issued, labelled `gov4git:directive` with body:

```
issue 3000 credits to @petar
```

### cast a vote on issue #1 and pr #5

```sh
gov4git --config testing-config.json ballot vote --name pmp/motion/priority_poll/1 --choices rank --strengths 10.0

gov4git --config testing-config.json ballot vote --name pmp/motion/approval_poll/5 --choices rank --strengths 10.0
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

Check for an updated tally in the GitHub UI:
- https://github.com/gov4git/testing.project/issues/1
- https://github.com/gov4git/testing.project/pull/5

Check in the GitHub UI that issue #1 and pr #5 show one eligible reference to each other.

### cancel pr #5

Cancel pr #5 through the GitHub UI.

Check that issues #1 and #2:
- were unfrozen
- the set of eligible references was updated to empty

[ Reopen pr #5 through the UI, for the next experiment. ]

### vote for pr #30

```sh
gov4git --config testing-config.json ballot vote --name pmp/motion/approval_poll/30 --choices rank --strengths 10.0
```

Verify it shows up as eligible for issues #1 and #2, and vice-versa.

### merge pr #30

Verify issues #1 and #2 closed.

Verify bounty and reward disbursed.

[ Revert pull request and re-open. ]

### create a dashboard issue
