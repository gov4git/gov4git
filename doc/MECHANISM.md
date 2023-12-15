# The mechanism of participation

## Plural Management Protocol specification (v0.0.1)

Every managed issue is associated with a _priority poll_ where community members can asynchronously cast votes to prioritize or de-prioritize the issue. Voting is conducted by spending credits which are held in an _escrow account_ associated with the issue.

Every managed pull request (PR) is associated with an _approval poll_ where members cast votes to express their approval (or dis-approval) of the proposed PR, and implicitly their bet on the outcome of the PR. Voting credits spent on the approval poll are held in a _reward account_ associated with the PR.

The author of a managed PR can make a claim that the PR resolves one or more outstanding (i.e. open) issues. To claim that the PR  resolves an issue, the author must include the expression `claims ISSUE_URL` anywhere in the description of the PR.

Whenever there is an _eligible_ PR that claims to resolve an issue, the issue is _frozen_ in that is ceases to accept new votes. An issue will be unfrozen whenever no eligible PRs refer to it.
A PR is _eligible_, whenever it has a positive priority score.

When a PR is accepted (i.e. merged and closed), the governance system will:
- Close all resolved issues.
- Collect the escrows of all resolved issues into a _bounty_ which is given to the author of the PR.
- Refund the escrowed credits of those who voted for the PR.
- Distribute the escrowed credits of those who voted against the PR to those who voted for the PR, using a plural proportioning rule.

When a PR is rejected (i.e. closed and not merged), the governance system will:
- Refund all users who voted on the PR.

Under the _plural proportioning rule_, a sum is distributed to a set of QV voters in proportion to the strength of their quadratic votes (equivalently, the square root of the credits they spent to vote).

## Notes on strategy

