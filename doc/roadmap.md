
## Roadmap

### Milestone 1: Barebone QV prioritization polling (finished, end of Oct 2022)

This milestone completed at end of Oct 2022.

#### Develop
- Decentralized apps over git framework:
  - [x] Identity over git and DNS
  - [x] Signed mail over git
- Community management:
  - [x] User and group management
  - [x] User balances
- Governance:
  - [x] Verifiable ballots: voting, tallying
  - [x] QV-based tallying

#### Validation
- [x] Demonstration to RadX (Glen/Leon/Alex)
- Feedback:
  - Aligned on voting model. 
    - Glen's use case:
      - Community organizer distributes voting credits
      - Users can up/down vote open PRs and issues
      - PR/issue closure credits/debits user voting accounts
  - Web UI required for live deployment with Plurality Book users

### Milestone 2: Framework rewrite + community features (targeting end of Nov 2022)

#### Develop
- Decentralized apps over git framework:
     - [x] Rewrite to enable easy extensibility and rapid feature development
       - [x] Identity
       - [x] User and group management
       - [x] Signed mail
       - [x] RPC over signed git mail
     - [x] Generic data structures over git (key-value, etc)
- Community services:
  - [x] Bureau: governance operational proposals by users
    - [x] User-initiated balance transfers
- Governance:
  - [x] User balances
  - [x] Balance holds during ballots
  - [x] Balance deductions on ballot closure/clearance

#### Validation
- [ ] Dogfood deployment on gov4git repo
  - [ ] Dogfooder docs
  - [ ] Invite dogfooders (PL/RadX)

### Milestone 3a: Document and socialize

### Milestone 3b: Web app and browser extension for voting

#### Develop
- Stand-alone TypeScript client library:
  - View open ballots and current tallies
  - View user balances
  - Voting
  - Balance transfers
- Web app
  - Dashboards for PR/issue prioritization based on current tallies
- Chrome extension for GitHub

#### Validation
- Dogfood on gov4git repo
- Dogfood on libp2p/IPFS/Filecoin?
- Dogfood on Plurality Book project

### Milestone 4: Verifiable change arbitration and approval

#### Develop
- Governance:
  - Directory-level change approval policies
  - Change approval arbitration
  - Verification of compliance
  - Arbitration library:
    - Quorum (at least N approvals out of M members)
    - Quadratic vote

### Milestone 5: QV mechanism research

#### Develop
- Governance:
  - Facilities for programmable voting logics and data analysis
  - Cluster-based QV tallying for promoting vote diversity
