#!/bin/sh

# This shell script walks through the end-to-end UX workflow that involves:
# - Initializing participant identies
# - Initializing governance
# - Managing users and groups
# - Administering a ballot, casting votes and tallying

# There are two types of participants in a governed community: the community organizer and the community members.
# The community organizer is the one who owns the community's public git repository.

# __Step 1__ Every participant begins by installing gov4git.
go install github.com/gov4git/gov4git/gov4git@latest

# __Step 2__ Participants prepare git repos for their identities and the community.

# __Step 2.1__ Every participant creates a pair of public and private git repos that will represent their identity in the system.
# For instance, I use GitHub for git hosting. I would create two repos:
# - github.com/petar/gov4git.public (a public GitHub repo)
# - github.com/petar/gov4git.private (a private GitHub repo)

# __Step 2.2__ The organizer creates a pair of public and private git repos for the community.
# If a public community repo already exists, it can be reused.
# For the example here, the organizer will create two GitHub repos:
# - github.com/petar/community.public (a public GitHub repo)
# - github.com/petar/community.private (a private GitHub repo)

# __Step 3__ Every participant prepares their local gov4git client configuration.
# The configuration describes the participant's own identity (public and private) repos,
# as well as the community's that they intend to interact with.
# In this example, a participant might use:

cat <<EOF >> ~/.gov4git/config
{
     "cache_dir": "/Users/petar/.gov4git/cache",
     "cache_ttl_seconds": 120,
     "auth" : {
          "git@github.com:petar/community.public.git": { "ssh_private_keys_file": "/Users/petar/.ssh/id_rsa" },
          "git@github.com:petar/community.private.git": { "ssh_private_keys_file": "/Users/petar/.ssh/id_rsa" },
          "git@github.com:petar/gov4git.public.git": { "ssh_private_keys_file": "/Users/petar/.ssh/id_rsa" },
          "git@github.com:petar/gov4git.private.git": { "ssh_private_keys_file": "/Users/petar/.ssh/id_rsa" }
     },
     "gov_public_url": "git@github.com:petar/community.public.git",
	"gov_public_branch": "gov",
	"gov_private_url": "git@github.com:petar/community.private.git",
	"gov_private_branch": "gov",
	"member_public_url": "git@github.com:petar/gov4git.public.git",
	"member_public_branch": "main",
	"member_private_url": "git@github.com:petar/gov4git.private.git",
	"member_private_branch": "main"
}
EOF

# `cache_dir`, if specified, sets the on-disk location for local caches of remote repos
# `cache_ttl_seconds` sets the ttl for cached replicas of community and member repos
# `ssh_private_keys_file` points to a local file containing your SSH credentials for cloning the private repos in the config.
# `community_public_url` is the git URL of the public community repo
# `community_public_branch` is the branch in the public community repo where the governance state will reside
# `community_private_url` is the git URL of the private community repo
# `community_private_branch` is the branch in the private community repo where the governance private keys will reside
# `member_public_url` is the git URL of your public identity repo
# `member_public_branch` is the branch in your public identity repo where public keys will reside
# `member_private_url` is the git URL of your private identity repo
# `member_private_branch` is the branch in your private identity repo where private keys will reside

# Note that only the organizer of the community has to fill in `community_private_url` and `community_private_branch`.

# Also note that gov4git supports SSH as well as HTTPS urls for git repos.
# If you were to rewrite the above config using HTTPS, it would look like this:

cat <<EOF >> ~/.gov4git/config
{
     "cache_dir": "/Users/petar/.gov4git/cache",
     "cache_ttl_seconds": 120,
     "auth" : {
          "https://github.com/petar/community.public.git": { "access_token": "YOUR_ACCESS_TOKEN" },
          "https://github.com/petar/community.private.git": { "access_token": "YOUR_ACCESS_TOKEN" },
          "https://github.com/petar/gov4git.public.git": { "access_token": "YOUR_ACCESS_TOKEN" },
          "https://github.com/petar/gov4git.private.git": { "access_token": "YOUR_ACCESS_TOKEN" }
     },
     "gov_public_url": "https://github.com/petar/community.public.git",
	"gov_public_branch": "gov",
	"gov_private_url": "https://github.com/petar/community.private.git",
	"gov_private_branch": "gov",
	"member_public_url": "https://github.com/petar/gov4git.public.git",
	"member_public_branch": "main",
	"member_private_url": "https://github.com/petar/gov4git.private.git",
	"member_private_branch": "main"
}
EOF

# By default gov4git looks for this config file in `~/.gov4git/config`, however you can specify a custom location
# using the `--config` flag, e.g. `gov4git --config /path/to/config.json ...`

# __Step 4__ Initialize identities

# __Step 4.1__ Every participant initializes their own identity.
# This results in generating new public and private keys and
# populating the participant's public and private identity repos.
gov4git init-id

# __Step 4.2__ The community organizer initializes the governance application
# This creates a new branch in the community's public repo, dedicated to tracking the state of governance.
# Also, public and private keys are generated for the community itself and
# the respective public and private community repos are populated.
gov4git init-gov

# __Step 5__ The community organizer adds some users to the community
gov4git user add --name petar --repo git@github.com:petar/gov4git.public.git --branch main

# All users are automatically made members of group `everybody`
gov4git group list --name everybody

# The community organizer can add additional groups of users ...
gov4git group add --name contributors

# ... and associate community users with a group.
gov4git member add --user petar --group contributors

# grant a few voting_credits tokens to user petar
# voting_credits can be spent on ballots.
gov4git balance add --user petar --key voting_credits --value 30.00

# __Step 6__ The community organizer opens a new ballot
# A ballot is a mechanism for soliciting community votes on a set of choices.
# Ballots can be configured to use a variety of voting and tallying strategies, and
# users can define their own.
#
# This example creates a new ballot using a simple default QV policy.
# The ballot has a name `issue/1`, a title `Issue ` and
# an arbitrary description which points to a GitHub issue in this case.
# This ballot has only one choice, named `i1` which voters can up/down vote.
# The ballot is open to all users in the group `contributors`.
gov4git ballot open --name issue/1 --title "Issue 1" --desc "https://github.com/petar/community.public/issues/1" --group contributors --choices "i1"

# __Step 7__ Participants in the ballot group can cast up/down votes asynchronously until the ballot is closed
gov4git ballot vote --name issue/1 --choices i1 --strengths -1.0

# __Step 8__ Occasionally, the community organizer fetches votes and updates the running tally
gov4git ballot tally --name issue/1

# __Step 9__ The community organizer closes the ballot, when its time to conclude and call the outcome
gov4git ballot close --name issue/1
