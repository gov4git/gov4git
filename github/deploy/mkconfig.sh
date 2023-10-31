#!/bin/sh
#
# Usage: mkconfig.sh GITHUB_REPO ORGANIZER_GITHUB_TOKEN

GITHUB_REPO=$1
ORGANIZER_GITHUB_TOKEN=$2

CACHE_DIR=""
GOV_PUBLIC_REPO_URL="https://github.com/${GITHUB_REPO}-gov.public.git"
GOV_PRIVATE_REPO_URL="https://github.com/${GITHUB_REPO}-gov.private.git"

CONFIG_JSON=$(
     jq -n \
          --arg cache_dir "$CACHE_DIR" \
          --arg gov_pub_repo "$GOV_PUBLIC_REPO_URL" \
          --arg gov_priv_repo "$GOV_PRIVATE_REPO_URL" \
          --arg gov_auth_token "$ORGANIZER_GITHUB_TOKEN" \
          '{
               "cache_dir": $cache_dir,
               "auth" : {
                    ($gov_pub_repo): { "access_token": $gov_auth_token },
                    ($gov_priv_repo): { "access_token": $gov_auth_token },
                    "git@github.com:petar/gov4git.public.git": { "ssh_private_keys_file": "/Users/petar/.ssh/id_rsa" },
                    "git@github.com:petar/gov4git.private.git": { "ssh_private_keys_file": "/Users/petar/.ssh/id_rsa" }
               },
               "gov_public_url": $gov_pub_repo,
               "gov_public_branch": "main",
               "gov_private_url": $gov_priv_repo,
               "gov_private_branch": "main",
               "member_public_url": "git@github.com:petar/gov4git.public.git",
               "member_public_branch": "main",
               "member_private_url": "git@github.com:petar/gov4git.private.git",
               "member_private_branch": "main"
          }'
)
echo $CONFIG_JSON
