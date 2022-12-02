# gov4git sync

## Next: 2022-12-05

- __Updates from the week:__
  - Discovered need to add an authentication manager to manage user's credentials for different repos (while testing dogfood setup)
    - Upgraded lib4git and gov4git with a new auth manager
  - [ScuttleButt](scuttlebutt.nz) is a nice architectural alternative to building dapps over git. Plan to include a comparison in the whitepaper.
  - Found [git appraise](https://github.com/google/git-appraise) and [git bug](https://github.com/MichaelMure/git-bug) as examples of highly-used applications built on top of git. 
  - We have early adopter access to [GitHub Blocks (beta)](https://blocks.githubnext.com/)
  - First draft of dogfood documentation (and UX flow) is in place
    - [How to join the community](../dogfood/how-to-join.md)
    - [How to participate](../dogfood/how-to-participate.md)
  - Three issues available to dogfood prioritization polling
  - New release [v0.0.2](https://github.com/gov4git/lib4git/releases/tag/v0.0.2) of [lib4git](https://github.com/gov4git/lib4git) including:
    - Support for embedding and tracking repos inside other repos. This is a key functionality for efficient implementation of communication for communities of significant scale (pertains to later milestones in our roadmap).
    - Support for local, on-disk read/write caching of remote repos. This is a key functionality for optimizing the "snappiness" of the gov4git UX. (The functionality will be adopted by the gov4git client at a later opportunity.)
    - Introduced a new cloning abstraction that carries authentication credentials. This was necessary for GitHub automation purposes.
    - Bug fixes.
  - Created first-draft GitHub automation for governance

- __Priorities for this week:__
  - Migrate gov4git governance to dedicated repo
  - Migrate GitHub automation
  - Release dogfood experiment
  - Explore GitHub blocks for gov4git

## Next: 2022-11-28

- __Updates from the week:__
  - Kicked off UI engineering discussions
  - Deployed gov4git dogfooding on gov4git
  - Started looking into GitHub automation for the dogfood deployment:
    - [Use a GitHub form to request joining a community](https://github.com/gov4git/gov4git/issues/new?assignees=petar&labels=community&template=join.yml&title=I%27d+like+to+join+this+project%27s+community)

- __Priorities for this week:__
  - Continue automation and documentation for dogfooding
    - Landing page with instructions for dogfooders
    - Opening and closing polls on issue creation using GitHub actions
  - Plan out a whitepaper

## Next: 2022-11-21

- __Updates from the week:__
  - Rewrite done (merged)
  - Last round of MVP features done (merged)
  - UI project job posting is out (RadX/Microsoft/Trigram/etc)
  - New product [UX walkthrough](../walkthrough.sh)

- __Priorities for this week:__
  - Commence work on UI
  - Dogfood deployment on gov4git
    - GitHub automation
      - Open/close polls for opened issues and PRs
      - Cron for poll tallies
    - Dogfooder documentation
    - Find dogfood cohort (PL/RadX/etc), send invites

## Next: 2022-11-14

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

## Next: 2022-11-07

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
