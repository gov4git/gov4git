- update GOVERN.md with known subtleties:
  - vulnerability: community members must not be allowed to close issues to prevent them from closing before a pr is submitted
- update GOVERN-PLAYBOOK.md

- metrics
  - test integration of metrics

- voting on the client must support different versions of the motion strategy
  - [ ] record policy/strategy in ballot/motion
  - compute reward in margin calc
    - motion writes bounty to ballot strategy instance state
    - custom js to compute reward

- populate history with ballot, motion, account, user IDs
- populate history with policy/strategy info

- add software revision of cron software to every commit

- dynamic ownership graph over account IDs
  - objects (like ballot, motion, etc) register themselves with the ownership system?
  - ballot can query about its owner (accounts already have owners)


- configuration
  - management strategy