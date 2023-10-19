package common

import (
	"testing"

	"github.com/gov4git/lib4git/ns"
)

func TestBallotTopic(t *testing.T) {
	testBallotName := ns.NS{"a", "b", "c"}
	expBallotTopic := "ballot:a/b/c"
	if BallotTopic(testBallotName) != expBallotTopic {
		t.Errorf("expecting %v, got %v", expBallotTopic, BallotTopic(testBallotName))
	}
}
