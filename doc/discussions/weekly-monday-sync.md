# gov4git sync

## Next: 2022-10-14

- __Updates from the week:__
  - Rewrite almost done:
    - Identity management
    - Mailing
    - Member management
    - Balloting
    - Tests and testing framework
  - Staffing for UI:
    - Job [blurb](ui-job-description.md)
    - Socializing in RadX channels

- __Priorities for this week:__
  - Rewrite remaining:
    - End-to-end ballot tests
  - Implement the new features:
    - Balance transfers
    - Balance refunds dependent on ballot outcomes

## Next: 2022-10-07

- __Context:__ October high-level goals
  - Develop ground framework for decentralized apps over git:
    - Identity
    - Signed mail between identities (repos)
  - Develop first prototype of gov4git, focusing on community management features:
    - User and group management
    - User balance management
    - General ballot/referendum flow
    - Baseline QV-based prioritization polling for open issues/PRs on GitHub

- __Updates so far:__
  - Gave demo to RadX (Glen/Leon/Alex). Feedback:
    - Aligned on voting model and players' workflow.
      - Glen's use case:
        - Community organizer distributes voting credits
        - Users can up/down vote open PRs and issues asynchronously
        - PR/issue closure credits/debits user voting accounts
    - Web UI requested for live deployment with Plurality Book users (for MVP)
    - User-initiated balance transfers requested (for MVP)
      - Accepted this feature into roadmap. Rationale:
        - Introduces a significant but important demand on the underlying framework:
          - An abstraction for invoking/proposing governance operations
          - Analog of smart contract method invocation
        - Opportunities:
          - Prioritization applies to libp2p/IPFS/Filecoin community management
    - Captured details in [roadmap](../roadmap.md)

- __Priorities for this week:__
  - Push through rewrite as far as possible. Remaining:

- __Anything else:__
  - Glen made intro to govrn.io
  - The ground framework (decentralized apps over git) enables some curious apps:
    - Twitter over git (already mentioned)
    - DHT over git (very curious)
