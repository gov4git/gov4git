# Call for Web UI engineers for gov4git

## Problem

We recently conceived a [new protocol for governance (gov4git)](https://github.com/gov4git/gov4git) of open source communities based on git. From an application standpoint this project aims to enable rapid experimentation and adoption of community mechanisms, which are accountable and transparent by design. From an engineering standpoint this project explores the novel architectural idea of building robust decentralized social applications entirely backed by networked git repositories. 

We are looking to develop a basic suite of MVP frontend components: a TypeScript client library, a web UI, a Chrome browser extension.

From a user perspective, our protocol addresses community organization concerns such as: members, groups, polls, votes, referendums, approvals, rewards, policymaking, and so on. From a backend perspective, the state of the application is embedded in a git repository on the Internet.

We aim to develop a frontend for polling and voting. This includes three deliverables:
- _TypeScript client library for the API._ This client will clone public and private git repos, and read and write the application state to the repo.
- _Web UI_ with UX for voting, displaying poll tallies and manipulation of user balances.
- _Chrome/Safari/Mozilla browser extensions_ with UX for voting and viewing rankings of issues and PRs on GitHub pages.

Upon completion, our goal is to deploy the system to a few early adopter open source projects, including the [Plurality Book Project](https://protocol.ai/blog/protocol-labs-and-plurality-book/) as well as some of our own projects at Protocol Labs.

## Milestones

We expect that the project should comfortably take about 2 months to completion, with a working prototype deliverable midway into the project.

### Milestone 1: JavaScript/TypeScript client for gov4git

This milestone release of the client should include:
- Initializing a git repo with newly generated user credentials (Ed25519 keys)
- Fetching existing poll tallies from a git repo
- Placing an Ed25519-signed signed vote into a git repo

### Milestone 2: Web UI

This includes:
- Authentication of GitHub app
- Securely caching git login credentials in the browser
- Initializing a user identity
- Viewing tallies from a git repo
- Casting a vote

### Milestone 3: Chrome browser extension

This includes the same features provided by the web UI, overlaid on top of GitHub pages. In particular, GitHub issue and PR pages should feature current poll tallies and provide UX for up/down-voting.

### Milestone 4: Balance management

Balance management operations includes two features:
- Viewing user voting credit balances
- Balance transfers (users can send vote credits to others users)

For this milestone, these operations must be added to:
- JavaScript/TypeScript client
- Web UI
- Chrome browser extension
