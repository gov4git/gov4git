
function calcVote(tally, voterUser, voteChoice, targetVote) {

     var currentVote = tally.scores_by_user[voterUser][voteChoice].score;
     var currentCharge = tally.scores_by_user[voterUser][voteChoice].strength;
     var targetCharge = targetVote * targetVote;
     var cost = targetCharge - currentCharge;

     return {
          "current_vote": {
               "description": "Current vote by user on this ballot choice",
               "value": currentVote,
          },
          "target_vote" : {
               "description": "Target vote by user on this ballot choice",
               "value": targetVote,
          },
          "cost" : {
               "description": "Cost of changing the vote",
               "value": cost,
          },
          "reward" : {
               "description": "Possible reward to user under best outcome",
               "value": 0.0,
          },
     }
}
