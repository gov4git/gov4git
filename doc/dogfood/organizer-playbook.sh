
# open a ballot for an issue
gov4git ballot open --name issue/6 --title "Issue 6" --desc "https://github.com/gov4git/gov4git/issues/6" --group everybody --choices "issue-6" --use_credits

# cast a vote on an issue
gov4git ballot vote --name issue/6 --choices issue-6 --strengths=+3.0

# tally votes on an issue
gov4git ballot tally --name issue/6

# view user's voting credits balance
gov4git balance get --user=petar --key=voting_credits

# view user's voting credits on hold balance
gov4git balance get --user=petar --key=voting_credits_on_hold

# show current tally for an open issue
gov4git ballot show-open --name=issue/6

# show current tally for a closed issue
gov4git ballot show-closed --name=issue/6
