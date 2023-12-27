- update GOVERN.md with known subtleties:
  - vulnerability: community members must not be allowed to close issues to prevent them from closing before a pr is submitted
- update GOVERN-PLAYBOOK.md

- [ ] test github integration of metrics

- split history to metric and trace
- VoteEvent add policy

- voting on the client must support different versions of the motion strategy
  - compute reward in margin calc
    - motion writes bounty to ballot strategy instance state
    - custom js to compute reward

- add software revision of cron software to every commit

- upgrade ballot
  - flat string namespace (not path)
  - rename strategy to policy

- dynamic ownership graph over account IDs
  - objects (like ballot, motion, etc) register themselves with the ownership system?
  - ballot can query about its owner (accounts already have owners)

- configuration
  - management strategy