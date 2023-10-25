#!/bin/sh
#
# Usage: deploy.sh TOKEN GITHUB_REPO GOV4GIT_RELEASE

gov4git -v github deploy --token=$1 --project=$2 --release=$3
