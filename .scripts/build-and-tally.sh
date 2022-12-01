#!/bin/sh
# run this script from the root of the gov4git repo

cd gov4git
go build

mkdir -p ~/.gov4git

cat <<EOF >> ~/.gov4git/config.json
{
     "auth" : {
          "https://github.com/gov4git/gov4git.git": { "access_token": "$GITHUB_TOKEN" },
          "https://github.com/gov4git/gov4git.private.git": { "access_token": "$GITHUB_TOKEN" }
     },
     "gov_public_url": "https://github.com/gov4git/gov4git.git",
	"gov_public_branch": "gov",
	"gov_private_url": "https://github.com/gov4git/gov4git.private.git",
	"gov_private_branch": "gov"
}
EOF
cat ~/.gov4git/config.json

./gov4git ballot list-open --only_names | xargs -L 1 ./gov4git -v ballot tally --name
