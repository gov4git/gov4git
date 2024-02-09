package ballotproto

// Margin captures functions for computing vote marginals:
/*

	{
		"help": {
			"label": "Help",
			"description": "Description of ballot",
			"fn_js": "function() { return "This is a QV ballot." }",
		},
		"cost" : {
			"label": "Cost",
			"description": "Additional cost to reach a desired total impact",
			"fn_js": "function(voteUser, voteChoice, voteImpact) { ... }",
		},
		"impact" : {
			"label": "Impact",
			"description": "Additional impact to reach a desired total cost",
			"fn_js": "function(voteUser, voteChoice, voteCost) { ... }",
		},
	}

*/
type Margin struct {
	Help   *MarginCalculator `json:"help,omitempty"`
	Cost   *MarginCalculator `json:"cost,omitempty"`
	Impact *MarginCalculator `json:"impact,omitempty"`
}

type MarginCalculator struct {
	Label       string `json:"label"`
	Description string `json:"description"`
	FnJS        string `json:"fn_js"`
}
