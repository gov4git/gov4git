## SUMMARY

- function "show ballot spec and tally"
- implemented in [ballot/show.go](../../proto/ballot/ballot/show.go)
- some [context on ballots](list-open-ballots.md#context)

## PSEUDOCODE

Given a ballot name ([which is a filesystem path](list-open-ballots.md#context)), return the ballot advertisement and current tally (if available). Let the ballot name be `BALLOT_PATH`.

1. clone the governance repo of the community

2. parse `ballots/open/BALLOT_PATH/ballot_ad.json`

3. try parsing the current tally
     - parse the file `ballots/open/BALLOT_PATH/ballot_tally.json`, if it exists
     - if the file does not exist, the community organizer has not tallied this ballot yet
     - here is an example tally file

     ```json
     {
          "ballot_advertisement":{
               "community":{
                    "Repo":"git@github.com:gov4git/gov4git.git",
                    "Branch":"gov"
               },
               "path":["issue","8"],
               "title":"Issue 8",
               "description":"https://github.com/gov4git/gov4git/issues/8",
               "choices":["issue-8"],
               "strategy":"priority_poll",
               "participants_group":"everybody","parent_commit":"1b07688fe71d5808b2bf1474f69284e0f817ea8e"
          },
          "ballot_fetched_votes":[
               {
                    "voter_user":"petar",
                    "voter_address":{
                         "Repo":"https://github.com/petar/gov4git.public.git",
                         "Branch":"main"
                    },
                    "voter_elections":[
                         {"vote_choice":"issue-8","vote_strength_change":-1},{"vote_choice":"issue-8","vote_strength_change":-1},{"vote_choice":"issue-8","vote_strength_change":-1}
                    ]
               }
          ],
          "ballot_choice_scores":[
               {"choice":"issue-8","score":-1.7320508075688772}
          ]
     }
     ```

     - for ui purposes, we are interested in the section `ballot_choice_scores`, which lists the current score of every choice in the ballot. scores are derived from users' votes, using the internal logic of the ballot strategy in use.
