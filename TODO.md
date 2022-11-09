REWRITE
- [ ] clone private repos
- [ ] API versions of methods
- [ ] command-line hooks

PRODUCT

QV-based polling with balance updates
     - [ ] transfer user balance
       - [ ] mechanism for proposing operation executions
     - [ ] vote sanity checks (e.g. allowed elections, enough balance, XXX, TODO)
     - [ ] rename gov command to ballot: ballot open, ballot close, etc.

     - [ ] polling strategy decrements balances

     - [ ] move voting interactions to identity main branch

Merge approval workflow
     - [ ] ...
     - [ ] verify command

FRAMEWORK
     -> [ ] identity and gov configs should include branch spec
       - [ ] every service should have a config that reflects what it needs
       - [ ] name services based on role
  
     - [ ] Dir->interface, Local.Sub() returns DirWithinLocal with method RelPath (or AbsPath)

DOC
- [ ] terminal-cast type video tutorials
- [ ] documentation with mdBook https://github.com/rust-lang/mdBook

FEATURES
- [ ] semantic human-readable stack traces
