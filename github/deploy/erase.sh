#!/bin/sh
#
# Usage: erase.sh TOKEN GITHUB_REPO

set -e -x

gov4git -v github remove --token=$1 --repo=$2-gov.public
gov4git -v github remove --token=$1 --repo=$2-gov.private
