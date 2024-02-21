package ballotproto

// Margin captures functions for computing vote marginals.
// This is how Margin looks in JSON:
//
//		{
//			"help": {
//				"label": "Help",
//				"description": "Description of ballot",
//				"fn_js": "function() { return "This is a QV ballot." }",
//			},
//			"cost" : {
//				"label": "Cost",
//				"description": "Additional cost to reach a desired total impact",
//				"fn_js": "function(voteUser, voteChoice, voteImpact) { ... }",
//			},
//			"impact" : {
//				"label": "Impact",
//				"description": "Additional impact to reach a desired total cost",
//				"fn_js": "function(voteUser, voteChoice, voteCost) { ... }",
//			},
//		}
//
//	Arguments of the "cost" function:
//		- voteUser: user name of the voter
//		- voteChoice: ballot choice
//		- voteImpact: total vote impact, desired by the voter
//	Result of the "cost" function:
//		- the _additional_ cost to the voter to reach the desired impact
//
//	Arguments of the "impact" function:
//		- voteUser: user name of the voter
//		- voteChoice: ballot choice
//		- voteCost: the total cost (including what has already been charged) the user is willing to spend
//	Result of the "impact" function returns an array with two values:
//		- the _additional_ impact the user can add to their current impact
//		- the _additional_ impact the user can subtract from their current impact
//
//	Arguments of the "reward" function:
//		- voteUser: user name of the voter
//		- voteChoice: ballot choice
//		- voteImpact: total vote impact, desired by the voter
//	Result of the "reward" function:
//		- the _potential_ reward to the voter, assuming a favorable outcome
type Margin struct {
	Help   *MarginCalculator `json:"help,omitempty"`
	Cost   *MarginCalculator `json:"cost,omitempty"`
	Impact *MarginCalculator `json:"impact,omitempty"`
	Reward *MarginCalculator `json:"reward,omitempty"`
}

type MarginCalculator struct {
	Label       string `json:"label"`
	Description string `json:"description"`
	FnJS        string `json:"fn_js"`
}
