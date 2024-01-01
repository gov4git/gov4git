package ballotproto

import (
	"testing"
)

func TestBallotTopic(t *testing.T) {
	testBallotID := ParseBallotID("a/b/c")
	expBallotTopic := "ballot:a/b/c"
	if BallotTopic(testBallotID) != expBallotTopic {
		t.Errorf("expecting %v, got %v", expBallotTopic, BallotTopic(testBallotID))
	}
}
