package arb

const Tally = `
SYNOPSIS
Collect latest votes from all users in the ballot group.
Summarize results.

BASIC OPERATION
1. Clone the community repo locally at the ballot branch
2. Contact the public repos of all users in the ballot and fetch their latest votes
3. Store all votes as well as a tally summary in a new commit on the ballot branch
4. Push the ballot branch back to the community repo
`
