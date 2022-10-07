package arb

const Tally = `
SYNOPSIS
Collect latest votes from all users in the referendum group.
Summarize results.

BASIC OPERATION
1. Clone the community repo locally at the referendum branch
2. Contact the public repos of all users in the referendum and fetch their latest votes
3. Store all votes as well as a tally summary in a new commit on the referendum branch
4. Push the referendum branch back to the community repo
`
