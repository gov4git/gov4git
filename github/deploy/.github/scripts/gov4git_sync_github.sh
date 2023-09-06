#!/bin/sh
#
# This script imports the issues and pull requests from a GitHub project repository into a governance repository.
#
# The GitHub project variables are inferred from the GitHub action shell environment.
#    PROJECT_OWNER = GitHub owner of the project
#    PROJECT_REPO = GitHub repository name of the project
#
# The governance variables must be set in the GitHub action workflow:
#
#    GOV_PUBLIC_REPO_URL = HTTPS URL of the public governance repository
#    GOV_PRIVATE_REPO_URL = HTTPS URL of the private governance repository
#
# The authentication variables must be set in the GitHub action workflow:
#
#    ORGANIZER_GITHUB_TOKEN = authentication token for the project and governance repositories
#
# The auth token must have permission to write to the governance repositories and
# read the issues and pull requests from the project repository.

mkdir -p ~/.gov4git/cache

CACHE_DIR=~/.gov4git/cache

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
                    ($gov_priv_repo): { "access_token": $gov_auth_token }
               },
               "gov_public_url": $gov_pub_repo,
               "gov_public_branch": "main",
               "gov_private_url": $gov_priv_repo,
               "gov_private_branch": "main"
          }'
)
echo $CONFIG_JSON > ~/.gov4git/config.json
cat ~/.gov4git/config.json

gov4git --config=$HOME/.gov4git/config.json github import \
     --token=$ORGANIZER_GITHUB_TOKEN \
     --project=$PROJECT_OWNER/$PROJECT_REPO
