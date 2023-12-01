- implement pmp policy for motions
  - add to qv strategy: cost function of the form A*V^2+B*V+C
  - compute quadratic priority as a field either in Score or in the policy structure
  - capture reward reports when prs are approved/rejected

+ [√] when PR closes, extract approved/denied, and pass through to policy close method
+ PR approved
  + flow credits
    + PR → yes voters (bettors)
    + PR → matching pool
    + ISSUE → author

- unit tests for docket
- address XXX
