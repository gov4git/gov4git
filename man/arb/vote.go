package arb

const Vote = `
SYNOPSIS
Update vote on a ballot. A referendum can be a ballot or an approval proposal.

BASIC OPERATION
1. Clone your (voter) identity public repo locally
2. Create a new branch designated to hold your vote
	- the branch name uniquely identifies the referendum the vote is for:
		"vote#" + hash(ballot_advertisement)
3. Write and commit your vote in the form of two files:
	- vote: your vote ballot form
	- vote.signed.ed25519: signature form of signed vote using your identity's public signing key
4. Push branch back to your identity public repo
`
