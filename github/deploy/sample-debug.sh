#!/bin/sh

set -x -e

go install ../../gov4git
# ./sample-erase.sh
./sample-deploy.sh

gov4git --config sample-config.json user add --name petar --repo https://github.com/petar/gov4git-identity-public.git --branch main
gov4git --config sample-config.json account issue --to user:petar --asset plural -q 11000
gov4git --config sample-config.json account issue --to pmp+matching --asset plural -q 1000

# gov4git --config sample-config.json ballot vote --name pmp/motion/priority_poll/14 --choices rank --strengths 10.0
# gov4git --config sample-config.json ballot vote --name pmp/motion/approval_poll/13 --choices rank --strengths 20.0
# gov4git --config sample-config.json ballot show --name pmp/motion/approval_poll/13
# gov4git --config sample-config.json motion show --name 10
