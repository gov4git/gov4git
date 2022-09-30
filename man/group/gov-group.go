package group

const GovGroup = `
SYNOPSIS

Manage community groups.

A governed community maintains a general-purpose directory of groups.
Each group is associated with:
- a unique community-specific name
- a dictionary (key/value store) of general-purpose group-specific information

BASIC OPERATION

File containing required group information (e.g. public group URL):
   /.gov/groups/NAME/info

File containing value of KEY, associated with group NAME:
   /.gov/groups/NAME/meta/KEY
`
