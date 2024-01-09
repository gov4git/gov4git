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

### Concerns and proposals

Collaborative projects, whether open-source or enterprise, often operate based on a set of core principles for managing collaboration. These principles primarily revolve around two key devices, which we call "concerns" and "proposals". They correspond to GitHub issues and pull requests, respectively, when using Gov4Git with GitHub.

### Goals of management

Managing a project largely revolves around repeatedly solving two problems: prioritization (of concerns) and decision-making (on proposals).
