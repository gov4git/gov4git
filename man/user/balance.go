package group

const GovUserBalance = `
SYNOPSIS
Manage balances of community users.

Users can be associated with labelled balances.
Balances are stored in the user's key/value store under keys starting with "balance:".

BASIC OPERATION

File containing the balance BAL of user USER:
   /.gov/users/USER/meta/balance:BAL
`

const GovUserBalanceAdd = `
SYNOPSIS
Add (or subtract) an amount from a user balance.
`

const GovUserBalanceMul = `
SYNOPSIS
Multiply (or divide) a user balance by a scaling factor.
`

const GovUserBalanceSet = `
SYNOPSIS
Set a user balance.
`
