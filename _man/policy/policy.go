package policy

const GovPolicy = `
SYNOPSIS

Manage directory policy for accepting changes.

Each directory can specify a policy that governs how changes inside the directory can be accepted.
Directory changes are acceptable if the policy of any ancestor directory applies.

Currently supported policies are:
- Quorum: approval by a minimum number of members of a given group

BASIC OPERATION

File containing policies for a directory DIR are in:
   /DIR/.gov/policy
`
