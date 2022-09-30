package group

const GovUser = `
SYNOPSIS

Manage community users.

A governed community maintains a general-purpose directory of users.
Each user is associated with:
- a unique community-specific name
- a public (git repo) URL
- a dictionary (key/value store) of general-purpose user-specific information

BASIC OPERATION

File containing required user information (e.g. public user URL):
   /.gov/users/NAME/info

File containing value of KEY, associated with user NAME:
   /.gov/users/NAME/meta/KEY
`
