# Governing playbook

All governance operations performed through the GitHub UI will take effect in about two minutes.

Only GitHub users who are collaborators to the project and have maintainer or administrator permissions can perform governance operations.

Once an operation is performed, the system will respond to a GitHub issue with a confirmation comment.

## Managing issues and pull requests

Gov4Git does not manage your existing issues and pull-requests by default. To include an issue or a PR in the management mechanism, apply the label `gov4git:managed` to the GitHub issue/PR.

In a typical workflow, the maintainer will configure the GitHub project to accept issues and PRs from anyone. Project maintainers will briefly skim new issues on a regular basis and apply the label `gov4git:managed` to the issues and PRs they want to manage via Gov4Git.

### Issues

Issues under management are associated with a prioritization ballot. Community members can cast votes for or against the issue by spending credits, which are locked in an escrow account.

When an issue is eventually resolved by a PR, the escrow amount is rewarded to the author of the PR.

An issue can be cancelled by closing it on GitHub. In this event, the escrow amount is refunded to the community voters.

### Pull requests

Gov4Git requires PRs to include a clear indication of which issues they are trying to address. This is accomplished by including zero or more instances of the following pattern anywhere in the PR description:

```
claims https://github.com/ORG/REPO/issues/ISSUE
```

When an issue that has a positive vote is claimed by at least one PR with a positive vote, the issue will be frozen (i.e. it will cease to accept votes) until the PR is accepted or rejected.

If the PR is accepted, the issue will be closed and the credits spent on the issue will be rewarded to the PR author.

If the PR is rejected, the issue will be unfrozen, and it will resume accepting votes from community members.

#### Accepting and rejecting

Gov4Git will not accept or reject a PR automatically. It is the maintainer's responsibility to do so.

To accept a PR, the maintainer would simply merge it using the standard GitHub mechanism. This will trigger Gov4Git to make a record of the decision and disburse rewards according to the Plural Management Protocol.

To reject a PR, the maintainer would simply close the PR without merging. This would trigger Gov4Git to perform refunds according to the management protocol.

## Membership

### Approve a membership request

Respond to the membership request GitHub issue with the comment `approved`.

### Reject a membership request

Close the membership request GitHub issue, or take no action.

## Economics

### Issue credits to a user

Create a GitHub issue, labelled `gov4git:directive`, containing a sentence of the form:

```
issue 30.5 credits to @user
```

### Transfer credits from one user to another

Create a GitHub issue, labelled `gov4git:directive`, containing a sentence of the form:

```
transfer 51 credits from @user1 to @user2
```
