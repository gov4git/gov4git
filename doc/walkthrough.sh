#!/bin/sh

# This shell script walks through the end-to-end UX workflow that involves:
# - Initializing participant identies
# - Initializing governance
# - Managing users and groups
# - Administering a ballot, casting votes and tallying

# There are two types of participants in a governed community: the community organizer and the community members.
# The community organizer is the one who owns the community's home git repository.

# __Step 1__ Every participant begins by installing gov4git.
go install github.com/gov4git/gov4git/gov4git

# __Step 2__ Participants prepare git repos for their identities and the community.

# __Step 2.1__ Every participant creates a pair of home and vault git repos that will represent their identity in the system.
# For instance, I use GitHub for git hosting. I would create two repos:
# - github.com/petar/gov4git.home (a public GitHub repo)
# - github.com/petar/gov4git.vault (a private GitHub repo)

# __Step 2.2__ The organizer creates a pair of home and vault git repos for the community.
# If a home community repo already exists, it can be reused.
# For the example here, the organizer will create two GitHub repos:
# - github.com/petar/community.home (a public GitHub repo)
# - github.com/petar/community.vault (a private GitHub repo)

# __Step 3__ Every participant prepares their local gov4git client configuration.
# The configuration describes the participant's own identity (home and vault) repos,
# as well as the community's that they intend to interact with.
# In this example, a participant might use:

cat <<EOF >> ~/.gov4git/config
{
	"ssh_private_keys_file": "/Users/petar/.ssh/id_rsa",
	"community_home_url": "git@github.com:petar/community.home.git",
	"community_home_branch": "gov",
	"community_vault_url": "git@github.com:petar/community.vault.git",
	"community_vault_branch": "main",
	"member_home_url": "git@github.com:petar/gov4git.home.git",
	"member_home_branch": "main",
	"member_vault_url": "git@github.com:petar/gov4git.vault.git",
	"member_vault_branch": "main"
}
EOF

# `ssh_private_keys_file` points to a local file containing your SSH credentials for cloning the vault repos in the config.
# `community_home_url` is the git URL of the home community repo
# `community_home_branch` is the branch in the home community repo where the governance state will reside
# `community_vault_url` is the git URL of the vault community repo
# `community_vault_branch` is the branch in the vault community repo where the governance private keys will reside
# `member_home_url` is the git URL of your home identity repo
# `member_home_branch` is the branch in your home identity repo where home keys will reside
# `member_vault_url` is the git URL of your vault identity repo
# `member_vault_branch` is the branch in your vault identity repo where private keys will reside

# Note that only the organizer of the community has to fill in `community_vault_url` and `community_vault_branch`.

# __Step 4__ Initialize identities

# __Step 4.1__ Every participant initializes their own identity.
# This results in generating new public and private keys and
# populating the participant's home and vault identity repos.
gov4git init-id

# __Step 4.2__ The community organizer initializes the governance application
# This creates a new branch in the community's home repo, dedicated to tracking the state of governance.
# Also, public and private keys are generated for the community itself and
# the respective home and vault community repos are populated.
gov4git init-gov

# __Step 5__ The community organizer adds some users to the community
gov4git user add --name petar --repo git@github.com:petar/gov4git.home.git --branch main

# All users are automatically made members of group `everybody`
gov4git group list --name everybody

# The community organizer can add additional groups of users ...
gov4git group add --name contributors

# ... and associate community users with a group.
gov4git member add --user petar --group contributors

# __Step 6__ The community organizer opens a new ballot
# A ballot is a mechanism for soliciting community votes on a set of choices.
# Ballots can be configured to use a variety of voting and tallying strategies, and
# users can define their own.
#
# This example creates a new ballot using a simple default QV strategy.
# The ballot has a name `issue/1`, a title `Issue ` and
# an arbitrary description which points to a GitHub issue in this case.
# This ballot has only one choice, named `i1` which voters can up/down vote.
# The ballot is open to all users in the group `contributors`.
gov4git ballot open --name issue/1 --title "Issue 1" --desc "https://github.com/petar/community.home/issues/1" --group contributors --choices "i1"

# __Step 7__ Participants in the ballot group can cast up/down votes asynchronously until the ballot is closed
gov4git ballot vote --name issue/1 --choices i1 --strengths -1.0

# __Step 8__ Occasionally, the community organizer fetches votes and updates the running tally
gov4git ballot tally --name issue/1

# __Step 9__ The community organizer closes the ballot, when its time to conclude and call the outcome
gov4git ballot close --name issue/1
