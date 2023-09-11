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

## Governing issues and prioritization

