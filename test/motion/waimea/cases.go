package waimea

type testCase struct {
	Voter0Credits          float64
	Voter1Credits          float64
	Voter0ConcernStrength  float64
	Voter0ProposalStrength float64
	Voter1ConcernStrength  float64
	Voter1ProposalStrength float64
	//
	Voter0EndBalance float64
	Voter1EndBalance float64
	AuthorEndBalance float64
}

var (
	/*
		   	                | user0 | user1 |
			issued     | 101     | 103    |
			concern  |  30     |  -20   |   cost of priority = 50
			proposal |  70     |  -10    |   cost of review = 80
			concern  | 5.477 | -4.472 |   priority score = 1.005
			proposal | 8.366 | -3.162 |   review score = 5.204

			projected priority bounty = 2 * 1.005 = 2.010
			projected review bounty = 2 * 5.204 = 10.408

			user0 end balance = 101-30-70+30+80 = 111
			user1 end balance = 103-20-10+20 = 93
			user2 end balance = 10.408+2.010 = 12.418
	*/
	testAcceptProposal = &testCase{
		Voter0Credits:          101.0,
		Voter1Credits:          103.0,
		Voter0ConcernStrength:  30.0,
		Voter0ProposalStrength: 70.0,
		Voter1ConcernStrength:  -20,
		Voter1ProposalStrength: -10.0,
		//
		Voter0EndBalance: 111.0,
		Voter1EndBalance: 93,
		AuthorEndBalance: 12.418,
	}

	/*
					 | user0 | user1 |
		issued     | 101     | 103    |
		concern  |  30     |  20   |   cost of priority = 50
		proposal |  70     |  -10    |   cost of review = 80
		concern  | 5.477 | -4.472 |   priority score = 1.005
		proposal | 8.366 | -3.162 |   review score = 5.204

		user0 end balance = 101-30-70 = 1
		user1 end balance = 103-20-10+80 = 153
		user2 end balance = 0
	*/
	testRejectProposal = &testCase{
		Voter0Credits:          101.0,
		Voter1Credits:          103.0,
		Voter0ConcernStrength:  30.0,
		Voter0ProposalStrength: 70.0,
		Voter1ConcernStrength:  20,
		Voter1ProposalStrength: -10.0,
		//
		Voter0EndBalance: 1.0,
		Voter1EndBalance: 153.0,
		AuthorEndBalance: 0.0,
	}

	/*
					 | user0 | user1 |
		issued     | 101     | 103    |
		concern  |  30     |  20   |   cost of priority = 50
		proposal |  70     |  -10    |   cost of review = 80
		concern  | 5.477 | -4.472 |   priority score = 1.005
		proposal | 8.366 | -3.162 |   review score = 5.204

		user0 end balance = 101-30-70+30+80 = 111
		user1 end balance = 103-20-10+20 = 93
		user2 end balance = 2*5.204 = 10.408
	*/
	testCancelConcernAcceptProposal = &testCase{
		Voter0Credits:          101.0,
		Voter1Credits:          103.0,
		Voter0ConcernStrength:  30.0,
		Voter0ProposalStrength: 70.0,
		Voter1ConcernStrength:  20,
		Voter1ProposalStrength: -10.0,
		//
		Voter0EndBalance: 111.0,
		Voter1EndBalance: 93.0,
		AuthorEndBalance: 10.408,
	}

	/*
					 | user0 | user1 |
		issued     | 101     | 103    |
		concern  |  30     |  20   |   cost of priority = 50
		proposal |  70     |  -10    |   cost of review = 80
		concern  | 5.477 | -4.472 |   priority score = 1.005
		proposal | 8.366 | -3.162 |   review score = 5.204

		user0 end balance = 101
		user1 end balance = 103
		user2 end balance = 0
	*/
	testCancelProposalCancelConcern = &testCase{
		Voter0Credits:          101.0,
		Voter1Credits:          103.0,
		Voter0ConcernStrength:  30.0,
		Voter0ProposalStrength: 70.0,
		Voter1ConcernStrength:  20,
		Voter1ProposalStrength: -10.0,
		//
		Voter0EndBalance: 101.0,
		Voter1EndBalance: 103.0,
		AuthorEndBalance: 0.0,
	}

	/*
					 | user0 | user1 |
		issued     | 101     | 103    |
		concern  |  30     |  20   |   cost of priority = 50
		proposal |  70     |  -10    |   cost of review = 80
		concern  | 5.477 | -4.472 |   priority score = 1.005
		proposal | 8.366 | -3.162 |   review score = 5.204

		user0 end balance = 101
		user1 end balance = 103
		user2 end balance = 0
	*/
	testCancelConcernCancelProposal = &testCase{
		Voter0Credits:          101.0,
		Voter1Credits:          103.0,
		Voter0ConcernStrength:  30.0,
		Voter0ProposalStrength: 70.0,
		Voter1ConcernStrength:  20,
		Voter1ProposalStrength: -10.0,
		//
		Voter0EndBalance: 101.0,
		Voter1EndBalance: 103.0,
		AuthorEndBalance: 0.0,
	}
)
