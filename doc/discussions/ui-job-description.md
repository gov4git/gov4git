# Call for Web UI engineers for gov4git

## Problem

We recently conceived a [new protocol for governance (gov4git)](https://github.com/gov4git/gov4git) of open source communities based on git. From an application standpoint, this project aims to enable rapid experimentation and adoption of community mechanisms, which are accountable and transparent by design. From an engineering standpoint, this project explores the novel architectural idea of building robust decentralized social applications backed by networked git repositories.

We are looking to develop a basic suite of MVP frontend components: a TypeScript client library, a web UI, and a Chrome browser extension.

From a participant perspective, our protocol addresses community organization concerns such as members, groups, polls, votes, referendums, approvals, rewards, policymaking, et cetera. Under the hood, from a backend perspective, a git repository on the Internet accounts for the application's state.

We aim to develop a frontend for polling and voting -- this includes three deliverables:
- _TypeScript client library for the API._ This client will clone public and private git repos and read and write the application state to the repo.
- _Web UI_ with UX for voting, displaying poll tallies, and manipulating participant balances.
- _Chrome/Safari/Mozilla browser extensions_ with UX for voting and viewing rankings of issues and PRs on GitHub pages.

Upon completion, our goal is to deploy the system to a few early adopter open source projects, including the [Plurality Book Project](https://protocol.ai/blog/protocol-labs-and-plurality-book/) as well as some of our projects at Protocol Labs.

## Milestones

The project should take about two months, with a working prototype deliverable midway into the project.

### Milestone 1: JavaScript/TypeScript client for gov4git

This milestone release of the client should include the following:
- Initializing a git repo with newly generated participant credentials (Ed25519 keys)
- Fetching existing poll tallies from a git repo
- Placing an Ed25519-signed vote into a git repo

### Milestone 2: Web UI

The second milestone includes:
- Authentication of the GitHub app
- Securely caching git login credentials in the browser
- Initializing a participant's identity
- Viewing tallies from a git repo
- Casting a vote

### Milestone 3: Chrome browser extension

The Chrome extension includes the same features provided by the web UI, overlaid on top of GitHub pages. In particular, GitHub issue and PR pages should feature current poll tallies and provide UX for up/down voting.

### Milestone 4: Balance management

Balance management operations include two features:
- Indication of participant voice credit balances
- Balance transfers (participants can send vote credits to other participants)

These operations must be added to all frontend components, in particular:
- JavaScript/TypeScript client
- Web UI
- Chrome browser extension
