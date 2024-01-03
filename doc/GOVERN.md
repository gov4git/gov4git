# Governance manual

This manual is for the community organizer.

There is a short [playbook](GOVERN-PLAYBOOK.md) summarizing the main operations from this document.

## Managing membership

### Approving requests to join the community

Non-members must join your project community before they can participate in governance.

Users can request to join directly from the [desktop app](https://github.com/gov4git/desktop-application/).

Join requests will appear as new issues in your GitHub project repository, assigned to the community organizer's GitHub user.

Membership requests can be granted by any repository collaborator who has _maintainer_ or _admin_ permissions on GitHub. To grant a request, reply to the issue on GitHub with the comment:

```
approved
```

In a couple of minutes Gov4Git automation will process the issue. It will reply with a comment indicating success — or a reason for failure — and will close the issue. In the event of a failure, you can re-open the issue to prompt the system to retry.

If you choose to deny a request, you do not need to take any action, however it is nice to explain your reasoning in the form of a comment.

## Managing economics

The economics of collaborative governance is based on an internal community currency called _plural credits_, or _credits_ for short.

Every community member has an account holding credits.

The community organizer has the power to issue (or withdraw) credits to any member.

### Issuing credits

The organizer can issue new credits and deposit them into the account of any community member.

To issue credits, the organizer creates a GitHub issue labelled `gov4git:directive`. The body of the issue includes a directive of the form:

```
issue 30.5 credits to @user
```

### Transferring credits

The organizer can transfer credits from one community member to another.

To transfer credits, the organizer creates a GitHub issue labelled `gov4git:directive`. The body of the issue includes a directive of the form:

```
transfer 51 credits from @user1 to @user2
```

## Managing collaboration

### Key concepts: Concerns and Proposals

Collaborative projects, whether open-source or enterprise-specific, often operate based on a set of core principles for managing collaboration. These principles primarily revolve around two key devices, which we call "concerns" and "proposals".

#### Concerns

A _concern_ is a device for initiating, discussing and tracking project improvement tasks and their dependencies. 

##### Examples of concerns across platforms

- On GitHub, GitLab, and Gitea: Concerns are referred to as "issues"
- On Jira: They are known as "tickets"

#### Proposals

A _proposal_ encapsulates a possible, sometimes partial, solution to address one or more concerns. This is typically in the form of a suggested change to the project's repository content.

##### Examples of proposals across platforms

- On GitHub and Gitea: Proposals are known as "pull requests"
- On GitLab: They are known as "merge requests"
- Other platforms may use "change requests" or similar phrases

#### Interoperability

gov4git understands and manipulates concerns and proposals natively. Moreover, it can interoperate with any source management system that employs similar mechanisms.

We adhere to the terminology of concerns and proposals. You can consider them synonymous with GitHub issues and GitHub pull requests.

### Governing collaboration

<!-- 
The role of governance in collaborative communities is to:

- Facilitate productive collaboration towards impactful outcomes
- Minimize regret in the event of subpar outcomes (of past decisions)

In this regard, governance must support the community in addressing three core day-to-day operational questions:

- _Direction_: What are the objectives of the community?

- _Allocation_: How should community resources be allocated towards solving for the objectives?

- _Adoption_: Which solutions should be adopted?

TBD -->

#### Concern prioritization

_Concern prioritization_ is a governance primitive, whose purpose is to derive a relative ordering of relevant concerns by deferring to the opinion of the community stakeholders in a pluralistic manner.

gov4git implements concern prioritization using Quadratic Polling. Each relevant concern is associated with a prioritization "ballot box". Community members can make positive or negative contributions to a ballot box by spending a desired amount voting credits from their personal balances.

The impact of a member's contributions to a concern is proportional to the square root of the amount of voting credits they spend into the respective ballot box.

The _(priority) score_ of a concern is the sum of impacts from all members.

##### Prioritization on GitHub

On GitHub, a collaborator that has permissions to assign labels to project issues can include an issue in prioritization, simply by applying the label `gov4git:prioritize` to the issue.

This mechanism allows the organizer — who is the GitHub owner of the project repository — to control which GitHub collaborators can administer matters of prioritization. Collaborators with GitHub triage permissions (to the project repository) can include issues in prioritization.
