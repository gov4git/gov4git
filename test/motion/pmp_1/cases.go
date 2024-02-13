package pmp

type testCase struct {
	MatchCredits          float64
	User0Credits          float64
	User1Credits          float64
	User0ConcernStrength  float64
	User0ProposalStrength float64
	User1ConcernStrength  float64
	User1ProposalStrength float64
	//
	User0EndBalance float64
	User1EndBalance float64
	User2EndBalance float64
}

var (
	/*
	   With matching funds = 0:

	   	              | user0 | user1 |
	   issued     | 101     | 103    |
	   concern  |  30     |  -20   |   funding = 50
	   proposal |  70     |  -10    |   funding = 80
	   concern  |  5.477 |  -4.472 |   concern vote = sqrt(spend)
	   proposal |  21.602 |  -8.164 | proposal vote = sqrt(spend * inverse_multiplier)

	   cost of priority = 50

	   cost = multiplier * vote^2
	   vote = sqrt( cost * inverse_multiplier )
	   inverse_multiplier = vote^2 / cost = 6.66

	   reviewer rewards:
	   - user0(positive voter) receives 2*21.60 = 43.20
	   - user1(negative voter) receives 0
	   contributor bounty:
	   - user2(author) receives remaining priority escrow = 50

	   end balance:
	   - user0 end balance = 101 - 30 - 70 + 43.20 = 44.20
	   - user1 end balance = 103 - 20 - 10 = 73
	   - user2 end balance = 50
	*/
	testCaseWithoutMatch = &testCase{
		MatchCredits:          0,
		User0Credits:          101.0,
		User1Credits:          103.0,
		User0ConcernStrength:  30.0,
		User0ProposalStrength: 70.0,
		User1ConcernStrength:  -20,
		User1ProposalStrength: -10.0,
		//
		User0EndBalance: 44.20,
		User1EndBalance: 73,
		User2EndBalance: 50,
	}

	/*
		With matching funds = 40:

					| user0 | user1 |
		issued     | 101     | 103    |
		concern  |  30     |  20   |   funding = 50
		proposal |  70     |  -10    |   funding = 80
		concern  |  5.477 |  4.472 |   concern vote = sqrt(spend)
		proposal |  28.982 |  -10.95 | proposal vote = sqrt(spend * inverse_multiplier)

		cost of priority = 50
		ideal deficit = 48.989
		priority score = 90
		match funds = 40
		match deficit = 48.989
		match ratio = 0.816


		cost = multiplier * vote^2
		vote = sqrt( cost * inverse_multiplier )
		inverse_multiplier = vote^2 / cost = 12.0

		reviewer rewards:
		- user0(positive voter) receives 57.965
		- user1(negative voter) receives 0
		contributor bounty:
		- user2(author) receives remaining escrow = 90

		end balance:
		- user0 end balance = 101 - 30 - 70 + 57.965 = 58.965
		- user1 end balance = 103 - 20 - 10 = 73
		- user2 end balance = 50
	*/
	testCaseWithMatch = &testCase{
		MatchCredits:          40,
		User0Credits:          101.0,
		User1Credits:          103.0,
		User0ConcernStrength:  30.0,
		User0ProposalStrength: 70.0,
		User1ConcernStrength:  20,
		User1ProposalStrength: -10.0,
		//
		User0EndBalance: 58.965,
		User1EndBalance: 73,
		User2EndBalance: 90,
	}
)
