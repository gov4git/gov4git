// Package pmp implements the Plural Management Protocol.
package pmp

import (
	"github.com/gov4git/gov4git/proto/ballot/common"
	"github.com/gov4git/gov4git/proto/docket/schema"
)

func MotionPollBallotName(id schema.MotionID) common.BallotName {
	return common.BallotName{"pmp", "motion", "poll", id.String()}
}
