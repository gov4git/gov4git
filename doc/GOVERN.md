# Governance manual

This manual is for you — the community organizer.

## Managing membership

You must provide a way for non-members to request membership into the community.

On GitHub this can be arranged conveniently by inviting your users to fill out an issue on your project repository, requesting membership. We have included a [sample GitHub issue template](../github/deploy/.github/ISSUE_TEMPLATE/join.yml), which instructs your prospective members how to install the governance desktop application and asks them for the relevant information:

- the HTTPS URL of their public gov4git identity repository (required), and
- their email, optionally (not used by gov4git)

The issue template automatically assigns these issues to the organizer's GitHub user to ensure you receive timely notifications upon new requests.

When an issue is filed, it will contain three relevant pieces of information:

- the applicant's GitHub user, which we denote `MEMBER_GITHUB_USER`
- the applicant's gov4git public repository URL, denoted `MEMBER_PUBLIC_REPO_URL`
- the branch within the applicant's public repository used for governance, denoted `MEMBER_PUBLIC_REPO_BRANCH`

If you choose to grant their request, you can add them to the community members, using the gov4git command-line tool:

```bash
gov4git user add \
     --name MEMBER_GITHUB_USER \
     --repo MEMBER_PUBLIC_REPO_URL \
     --branch MEMBER_PUBLIC_REPO_BRANCH
```

If you choose to deny their request, you can explain your reasoning in the form of a response comment on the GitHub issue itself.

## Concerns and proposals

Most collaborative projects — be it public and open-source, or private and enterprise-specific — share a few core workflow management principles.

They are invariably centered around two mechanisms: One for managing _concerns_, and one for managing _proposals_.

_Concerns_ are a device for nucleating, discussing and tracking project improvement tasks and their dependencies. 

In the context of GitHub, GitLab or Gitea, for instance, concerns are embodied by "issues". Other platforms may attribute different names to the analogous concept. For instance, Jira embodies concerns in the form of "tickets".

_Proposals_ encapsulate a solution (possibly partial) to one or more concerns in the form of a proposed change to the contents of the project repository.

GitHub and Gitea call their proposals "pull requests". GitLab uses the term "merge requests". And a litter of other terms, such as "change requests", can be found in other systems.

gov4git incorporates native support for concerns and proposals, and it can interoperate with any source management system that meets a few basic criteria. In particular, their analog of concerns and proposals supports:

- incorporating free-form text
- including references to other concerns and proposals
- assigning opaque labels for categorization

We adhere to the terminology of concerns and proposals, but the reader can treat them synonymously with GitHub issues and GitHub pull requests.

## Prioritization



XXX
