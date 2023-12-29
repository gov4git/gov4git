package ballotproto

// MarginCalcJS must be a JS function of this form:
//
//	function calcMargin(currentTally, voteUser, voteChoice, voteTarget) {
//		...
//		return {
//			"currentVote": {
//				"label": "Current vote",
//				"description": "Your current vote",
//				"value": currentVote,
//			},
//			"targetVote" : {
//				"label": "Target vote",
//				"description": "Your target vote",
//				"value": targetVote,
//			},
//			"cost" : {
//				"label": "Cost",
//				"description": "Cost of changing your vote",
//				"value": cost,
//			},
//			...
//		}
//	}
type MarginCalcJS string

type Margin struct {
	CalcJS MarginCalcJS `json:"calc_js"`
}
