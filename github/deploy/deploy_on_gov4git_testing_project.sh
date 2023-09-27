#!/bin/sh

gov4git -v github deploy \
     --token=$GITHUB_GOV4GIT_TESTING_TOKEN \
     --project=gov4git/testing.project \
     --release=v1.1.10
