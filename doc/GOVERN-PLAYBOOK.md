# Governing playbook

All governance operations performed through the GitHub UI will take effect in about two minutes.

Only GitHub users who are collaborators to the project and have maintainer or administrator permissions can perform governance operations.

Once an operation is performed, the system will respond to a GitHub issue with a confirmation comment.

## Issue prioritization

### Include any issue in prioritization polling

Apply the label `gov4git:prioritize` to the GitHub issue.

### Freeze polling on an issue

Lock the GitHub issue.

### Unfreeze polling on an issue

Unlock the GitHub issue.

## Membership

### Approve a membership request

Respond to the membership request GitHub issue with the comment `Approved`.

## Economics

### Issue credits to a user

Create a GitHub issue, labelled `gov4git:directive`, containing a sentence of the form:

```
Issue 30.5 credits to @user.
```

### Transfer credits from one user to another

Create a GitHub issue, labelled `gov4git:directive`, containing a sentence of the form:

```
Transfer 51 credits from @user1 to @user2.
```
