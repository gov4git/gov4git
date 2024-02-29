# Waimea: Collective management and compensation for open source communities

Waimea is a new Gov4Git governance mechanism we developed specifically for open-source software development and community management. It is available in [v2.2.0](https://github.com/gov4git/gov4git/releases/tag/v2.2.0) and after.

For context, [Gov4Git](https://gov4git.org) is a platform for augmenting the standard collaborative development workflow — based on issues and pull requests (in GitHub parlance) — with a governance and economics __mechanisms__. From a software perspective, Gov4Git provides an SDK for building mechanisms, as well as a runtime and an app for deploying mechanisms to projects on GitHub, or other similar systems.

## Governance and open source

Every collaborative project has governance. It is the set of rules and processes — be it explicit or implicit, ad-hoc or premeditated, written or unwritten — that determine __prioritization__ of concerns, __decision-making__ on proposals, and __attribution__ of credit to individual collaborators.

At first glance, governance dictates how participants collaborate, and how they are compensated for their effort. 

But, in fact, it is more than that. 

Governance — indirectly — also determines whether participants are willing to collaborate, for how long, and to what degree of involvement. In this regard, collaborators judge the governance system itself to determine if it is able to nurture the type of community that they want to be a member of. They judge if the rules of governance reflect their values and whether their contributions will be acknowledged fairly.

Software development that happens within the context of a business entity — such as startups, corporations, or foundations — inherits the governance (also known as "management") mechanisms and their enforcement from the business.

On the other hand open source — in its purest and most popular form — is practiced outside a business or organization. As a result, open source often defaults to a simple ad-hoc system of __governance, which is typically influenced by the available tools for enforcement, rather than the needs of the community__. For instance, most open source projects on GitHub are essentially forced to adopt a non-nuanced, two-tier hierarchical system of maintainers and contributors.

We believe that this inattention to the governance needs of open source communities is the main cause of well-known problems in open source — lack of __sustainable__ maintainership and support, lack of __funding__, __supply chain vulnerabilities__, as well as __community fragmentation__ upon commercialization or institutionalization of impactful projects.

We think that open source communities should have both tools and standard templates of governance, much like tech companies (from startups to large corporations) have tools and templates for all matters of business — management, finance, compensation, bookkeeping, HR, and so on.

The governance of any community must reflect its **ethos**, as well as the **circumstances** that brought members into the community.

While there is no one concept of an open source ethos, most varieties share a common core of tenets, which includes **transparency**, **inclusivity**, and **fairness**. Waimea is a governance mechanism which presents one possible formalization of these values:

- **Transparency** is provided by the Gov4Git protocol itself wherein, by design, the logic as well as the record of all operations of governance is kept in an always-accessible, immutable ledger

- Waimea embodies **inclusivity** by utilizing a flat organization whereby all collaborators are subject to the same rules and opportunities. (This is not to say that everybody is equal or undifferentiated within the community. Productive strategic behavior may award some with more effective power or stake than others.)

- Waimea embodies **fairness** by rewarding productive behavior that aligns with community preferences (and penalizing the opposite), in a quantified manner, using a market-inspired approach that leverages peer review.

## Principle of operation

Central to the operation of Waimea is a __virtual community currency__, called _credits_. Credits represent stake in the community’s project, and they can be spent within the community to exercise influence on the development trajectory. 

Credits are issued by the governance system, in a manner (described later) that ensures the total supply reflects the amount of work performed by the community since its inception. In other words, the community has an inflationary economy. 

Credits are held by community members, or staked into issues and PRs in the course of the day-to-day collaborative workflow. 

Waimea’s mechanism is centered on the relationship between issues and PRs. As usual, issues describe ideas or tasks or goals (such as bugs to fix, new features, requests for support, etc.) that the community may (or not) be interested in addressing. Whereas pull requests represent proposed solutions to one or more issues. 

### Issues and prioritization

Virtually every healthy project is perpetually in a state where there are more issues than there are community resources to address them. This is why the first key function of governance is to prioritize which issues are addressed first. 

To do this, Waimea associates a __priority poll__ to every project issue. For as long as an issue is open, community members can influence its priority (up or down) by staking a desired amount of their personal credits, expressing the intensity of their opinion. Eventually, when the issue is cancelled or resolved, staked credits are _refunded in full_. Nevertheless, participants are wise to stake their credits judiciously, as staked credits could not be used for anything else while weighing on priorities.  

Besides providing visibility into the community’s preferences, the priority of issues plays a role in incentives to undertake work. The author of a successfully merged pull request that addresses a given issue is awarded a bounty, one of whose components is proportional to the priority of the issues that it resolves. The funds for such bounties are issued by the governance system. The key concept here is that the priority of an issue represents the community's opinion about the value of resolving the issue, while the issuance of new funds ensures that the total supply of credits matches the total value of work that has been merged into the project.

### Pull requests and approval

Pull requests are proposals to make long-lasting changes to a collaborative project. Short of small adjustments that reviewers can request during review, pull requests call for a binary decision whether to adopt (or not) the proposed changes. Arbitrating this decision is the second key role of governance.

As with issues, collaborators can create pull requests freely. Waimea automatically associates an __approval poll__ with every pull request. Approval polls work similarly to priority polls — they provide a venue for community reviewers to express their opinion (for or against the adoption of a pull request) quantitatively and asynchronously, for as long as a pull request is open.

Like priority polls, approval polls require staking of credits, and therefore they involve similar strategic considerations on the part of participants.

Unlike priority polls, when a PR is closed (accepted or rejected) credits staked in the approval poll are awarded to the voters who predicted that outcome. Specifically, _voters who predicted the outcome are refunded in full_. Additionally, the total stake of voters who did not predict the outcome is distributed to the voters who predicted the outcome in proportion of the vote strengths of the latter. E.g. if a PR is accepted, all who voted "for" are refunded, and the stake of all who voted "against" is awarded to the "for" voters. This award mechanism is designed to encourage alignment amongst reviewers, as well as greater care and attention in the cases of controversial PRs.

In the event that a PR is accepted, its approval score also contributes to the bounty awarded to the contributing author — in addition to the contributions (described above) arising from resolved issues. The contribution to the bounty is proportional to the approval score and the funds for it are issued by the governance system.

Like the component of the bounty arising from resolved issues, the component arising from the approval score of a PR reflects the community's opinion of the value of the work associated with reviewing the PR. In the same spirit, the funds issued to fulfill the PR component of the bounty ensure that the main economic invariant of the community holds true — i.e. that the total supply of credits in circulation is a reflection of the total amount of work absorbed by the project.

### Credits as a form of ownership

As the reader has seen by now, our mechanism implements an inflationary economy, where new credits are issued upon successfully merging a PR. Furthermore, the amount issued aspires to represent the community's opinion of the value that this PR brings, along with the associated value of review.

As a result, the total supply of credits in circulation could be viewed as a peer-reviewed measure of the total value of the project in its current state.

Consequently, the credits owned by individual collaborators have a natural interpretation as stake into the project. In this regard, credits are akin to the stock shares that startup companies use to compensate their employees.

While credits are a virtual device, they are not worthless. They provide a mechanism to reward its past and present contributors retro-actively in real terms. In the event that a community project interfaces with the real economy — such as being remunerated for its impact (e.g. through a charitable grant), or is commercialized (as a startup), or is institutionalized (into a foundation) — its contributors could and should be compensated proportionally to their holdings of community credits.

We should note, that the only other source of credit issuance is the capitalization table of the founding members of the community. The first members must elect an initial distribution of credits for themselves, in order to introduce credits in circulation and thereby spark the collaborative economy.

Waimea also leaves open a provision for issuing credits to participants who contribute to the community in ways other than collaboration — such as contributing monetary or other real benefits to the community. It is left to the community to decide how to value such non-collaboration contributions in terms of credits.

## How things work in practice, on the ground

The playground for governed collaboration is the issue and pull request system that the community uses. We currently support GitHub, and are expecting to have integrations with other providers.

Waimea communities are permissionless — everyone is allowed to join a community by sending a join request from the desktop app. Upon joining, а new member begins with zero credits in their account. Members can earn credits by productive participation or by voluntary transfers from other members.

Participation in the community is based around the lifecycle of a code change, which begins with the creation of an issue describing a goal, and ends with the approval and merging of a PR that resolves the issue.

### Issues

Every collaborator can create issues freely. Waimea associates each issue with a __priority poll__, which determines the __priority score__ (a real number) of the issue and affects the bounty for the contributor who eventually resolves the issue through a pull request.

The priority poll is a mechanism which allows members to affect the priority score (up or down) by staking their credits to the issue. Polls are based on a quadratic design, whereby a stake of `P` credits affects the priority score by `SQRT(P)` points — up or down, as desired. Members are allowed to adjust their stakes asynchronously, as long as the issue is open and there is no eligible pull request claiming to resolve it.

The life of an issue terminates in one of two ways — it can be __cancelled__ (triggered by closing it manually on GitHub), or it can be __resolved__ (triggered by merging a PR that addresses it). In both cases, the credits staked by community members to the issue's priority poll are refunded.

In the event that an issue is resolved by a PR, the community governance issues new credits — equal in quantity to the priority score, if it is positive; zero, otherwise — and contributes them to the bounty for the author of the PR.

### Pull requests (PRs)

Similarly to issues, every collaborator can create pull requests freely. The contributor of the PR can make a claim that the PR resolves zero or more issues. To claim an issue, the description of the PR must include a claim statement of the form:

```
claims https://github.com/ORG/REPO/issues/ISSUE_NUMBER
```

Waimea associates each PR with an __approval poll__, which determines the __approval score__ (a real number) of the PR. The approval score affects the bounty awarded to the PR contributor, if the PR is accepted and merged.

The approval poll is a mechanism which allows members to affect the approval score (up or down) by staking credits to the PR. Approval polls are based on the same asynchronous quadratic mechanism used by priority polls (described above).

The life of a PR terminates in one of two ways — it can be __accepted__ (triggered by merging the PR on GitHub), or it can be __rejected__ (triggered by closing the PR on GitHub).

If a PR is accepted:
- The issues it claims are resolved (and automatically closed on GitHub), and their bounties are awarded to the author of the PR
- Reviewers who voted for the PR are refunded their stakes
- The credits staked by reviewers who voted against the PR are distributed to the reviewers who voted for, proportionally to the strength of their votes
- The community governance issues new credits whose quantity equals the approval score of the PR, and awards them to the contributor of the PR

If a PR is rejected:
- Reviewers who voted against the PR are refunded their stakes
- The credits staked by reviewers who voted for the PR are distributed to the reviewers who voted against, proportionally to the strength of their votes

### Frozen issues

A PR is said to make an __eligible__ claim on an issue, if the PR claims the issue and the PR's approval score — a number that can very dynamically — is bigger than zero.

Whenever an issue has eligible claims, it becomes __frozen__. Similarly, when there are no eligible claims it becomes __unfrozen__. Whenever an issue is frozen, votes or vote adjustments on its priority poll are not accepted.

### Maintainer vs autopilot governance

Each community has a distinguished set of __maintainer__ users. In the context of our integration with GitHub, the maintainers are the GitHub maintainers of the project repository.

Maintainers can have anywhere between a minimal and a full, hands-on involvement in the governance of the community.

At a minimum, maintainers are responsible for establishing who the first community members are, and issuing initial sums of credit to them.

On the other end of the spectrum, maintainers can exercise any of the following operations at their discretion, in a manner that is always transparent to the community:

- Accept new members to the community, or remove existing members
- Issue, burn, or transfer credits to, from, or between any accounts
- Cancel issues
- Accept or reject PRs regardless of the popular vote

## Strategic behavior

XXX

## FAQ

### Why do priority and approval polls use Quadratic Voting?

XXX

### When and how is a pull request closed (accepted or rejected)?

XXX

## Contact us and try our solution

We are actively looking for trial and early adopter communities. We think that the current release of Gov4Git, equipped with Waimea, is ideal for community managers and open-source communities of peers. [Get in touch with us](https://docs.google.com/forms/d/e/1FAIpQLSeO9obA-9jFFABMoN0Vjzcsmf9fRDKD5L9OiBq49MExUQ6b4A/viewform) and briefly tell us your use case. We are looking forward to providing our trial users and early adopters with support and attention every step of the way.
