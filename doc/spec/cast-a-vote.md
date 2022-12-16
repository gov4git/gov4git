## SUMMARY

- function "cast a vote on a ballot"
- implemented in [ballot/vote.go](../../proto/ballot/ballot/vote.go)
- some [context on ballots](list-open-ballots.md#context)

## CONTEXT

Casting a vote entails submitting a set of "elections" for a given open ballot in the community.
An [election](https://github.com/gov4git/gov4git/blob/main/proto/ballot/common/schema.go#L49) associates some voting strength (a number) to a ballot choice.

## PSEUDOCODE

To cast a vote, the user must specify the ballot name and a set of elections.

1. fetch the ballot ad

     - fetch the ballot advertisement, using the [method for showing a ballot](show-ballot.md)

2. prepare the vote envelope

     - the vote envelope is a [JSON data structure](https://github.com/gov4git/gov4git/blob/main/proto/ballot/common/schema.go#L56), as in this example

     ```json
     {
          "ballot_ad_commit": "1b07688fe71d5808b2bf1474f69284e0f817ea8e",
          "ballot_ad": { /* the ballot advertisement */ },
          "ballot_elections": [
               {
                    "vote_choice": "issue-8",
                    "vote_strength_change": 3.0
               },
          ]
     }
     ```

     `ballot_ad_commit` must be the commit hash of the governance branch used to retrieve the ballot in step (1).

     `ballot_ad` equals the ballot advertisement retrieved in step (1)

     `ballot_elections` holds the user's elections. Valid elections must refer to choices that are specified in the ballot ad. There is one exception to this rule. If no choices are listed in the ballot ad, then any string is a valid election choice.

3. send the vote envelope to the governance

     - use the [send](send.md) operation to send the vote envelope from the voter to the community governance, using topic `"ballot:BALLOT_NAME"`, where `BALLOT_NAME` is the path string corresponding to the ballot name

          - note that the ballot name is a list of strings, e.g. `["a", "b", "c"]`. the path string corresponding to this ballot name is `"a/b/c"`
          
          - for reference, the computation of the topic is [implemented here](https://github.com/gov4git/gov4git/blob/main/proto/ballot/common/schema.go#L29)
