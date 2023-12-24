package ballotproto

import (
	"testing"
)

func TestBallotTopic(t *testing.T) {
	testBallotName := BallotName{"a", "b", "c"}
	expBallotTopic := "ballot:a/b/c"
	if BallotTopic(testBallotName) != expBallotTopic {
		t.Errorf("expecting %v, got %v", expBallotTopic, BallotTopic(testBallotName))
	}
}
