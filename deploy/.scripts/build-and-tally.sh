#!/bin/sh
# run this script from the root of the gov4git repo

cd gov4git
go build

mkdir -p ~/.gov4git

cat <<EOF >> ~/.gov4git/config.json
{
     "auth" : {
          "https://github.com/gov4git/governance.git": {
               "user_password": {
                    "user": "$GOVERNANCE_ACCESS_USER",
                    "password": "$GOVERNANCE_ACCESS_TOKEN"
               }
          },
          "https://github.com/gov4git/gov4git.private.git": {
               "user_password": {
                    "user": "$GOVERNANCE_ACCESS_USER",
                    "password": "$GOVERNANCE_ACCESS_TOKEN"
               }
          }
     },
     "gov_public_url": "https://github.com/gov4git/governance.git",
	"gov_public_branch": "main",
	"gov_private_url": "https://github.com/gov4git/gov4git.private.git",
	"gov_private_branch": "main"
}
EOF
cat ~/.gov4git/config.json

./gov4git ballot list-open --only_names > open-ballots
echo Tallying open ballots:
cat open-ballots

cat open-ballots | xargs -t -L 1 ./gov4git -v ballot tally --name
