#!/bin/sh
#
# This script imports the issues and pull requests from a GitHub project repository into a governance repository.
#
# The GitHub project variables are inferred from the GitHub action shell environment.
#    GITHUB_PROJECT_OWNER = GitHub owner of the project
#    GITHUB_PROJECT_REPO = GitHub repository name of the project
#
# The governance variables must be set in the GitHub action workflow:
#
#    GOV_PUB_REPO = HTTPS URL of the public governance repository
#    GOV_PRIV_REPO = HTTPS URL of the private governance repository
#
# The authentication variables must be set in the GitHub action workflow:
#
#    GOV_AUTH_USER = authentication user for the project and governance repositories
#    GOV_AUTH_TOKEN = authentication token for the project and governance repositories
#
# The auth token must have permission to write to the governance repositories and
# read the issues and pull requests from the project repository.

GITHUB_PROJECT_OWNER=$(echo $GITHUB_REPOSITORY | cut -d'/' -f1)
GITHUB_PROJECT_REPO=$(echo $GITHUB_REPOSITORY | cut -d'/' -f2)

mkdir -p ~/.gov4git/cache

CACHE_DIR=~/.gov4git/cache

CONFIG_JSON=$(
     jq -n \
          --arg cache_dir "$CACHE_DIR" \
          --arg gov_pub_repo "$GOV_PUB_REPO" \
          --arg gov_priv_repo "$GOV_PRIV_REPO" \
          --arg gov_auth_user "$GOV_AUTH_USER" \
          --arg gov_auth_token "$GOV_AUTH_TOKEN" \
          '{
               "cache_dir": $cache_dir,
               "auth" : {
                    $gov_pub_repo: { "user_password": { "user": $gov_auth_user, "password": $gov_auth_token } },
                    $gov_priv_repo: { "user_password": { "user": $gov_auth_user, "password": $gov_auth_token } }
               },
               "gov_public_url": $gov_pub_repo,
               "gov_public_branch": "main",
               "gov_private_url": $gov_priv_repo,
               "gov_private_branch": "main"
          }'
)
echo $CONFIG_JSON > ~/.gov4git/config.json
cat ~/.gov4git/config.json

./gov4git github --config=~/.gov4git/config.json import \
     --token=$GOV_AUTH_TOKEN \
     --owner=$GITHUB_PROJECT_OWNER \
     --repo=$GITHUB_PROJECT_REPO
