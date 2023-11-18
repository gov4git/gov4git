- [x] prohibit reopening
- [x] integration test github sync
- [x] list all motions with corresponding motion and policy state (which includes ballot state)
- [x] re-score managed issues and prs after tally in cron
- [x] remove freeze triggers (eligible proposals) if referencing proposals drop in score (become ineligible)
  - [x] Policy.Update, called after rescoring
- [x] make a zero policy for derek + testing

———————

- [ ] designate a reference type for eligible proposals

- [ ] produce a report when freezing a motion (generically on all policy interactions)

- [ ] add a global threshold parameter for considering proposals as eligible solutions
  - fixed number? number relative to community capitalization?

———————

- [ ] implement pmp policy for motions
- [ ] capture reward reports when prs are approved/rejected

- [ ] unit tests for docket
