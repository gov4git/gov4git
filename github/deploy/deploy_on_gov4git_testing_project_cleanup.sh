#!/bin/sh

set -e -x

gov4git -v github remove \
     --token=$GITHUB_GOV4GIT_TESTING_TOKEN \
     --repo=gov4git/testing.project-gov.public

gov4git -v github remove \
     --token=$GITHUB_GOV4GIT_TESTING_TOKEN \
     --repo=gov4git/testing.project-gov.private
