
function calcMargin(currentTally, voteUser, voteChoice, targetVote) {

     var currentVote = currentTally.scores_by_user[voteUser][voteChoice].score;
     var currentCharge = currentTally.scores_by_user[voteUser][voteChoice].strength;
     var targetCharge = targetVote * targetVote;
     var cost = targetCharge - currentCharge;

     return {
          "currentVote": {
               "label": "Current vote",
               "description": "Your current vote",
               "value": currentVote,
          },
          "targetVote": {
               "label": "Target vote",
               "description": "Your target vote",
               "value": targetVote,
          },
          "cost": {
               "label": "Cost",
               "description": "Cost of changing your vote",
               "value": cost,
          },
     }
}
