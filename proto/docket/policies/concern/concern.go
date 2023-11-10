package concern

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/proto/ballot/qv"
	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/ns"
)

func init() {
	policy.Install(context.Background(), concernPolicy{})
}

const ConcernPolicyName = schema.PolicyName("concern-policy")

type concernPolicy struct{}

func (x concernPolicy) Name() string {
	return ConcernPolicyName.String()
}

func (x concernPolicy) Open(
	ctx context.Context,
	govAddr gov.GovOwnerAddress,
	govCloned gov.GovOwnerCloned,
	motion schema.Motion,
	instancePolicyNS ns.NS,

) {

	// open a poll for the motion
	ballotName := schema.MotionPollBallotName(motion.ID)
	ballot.Open_StageOnly(
		ctx,
		qv.QV{},
		gov.GovAddress(govAddr.Public),
		govCloned.Public,
		ballotName,
		fmt.Sprintf("Priority poll for motion %v", motion.ID),
		fmt.Sprintf("Up/down vote the priority of motion %v", motion.ID),
		[]string{schema.MotionPollBallotChoice},
		member.Everybody,
	)

}

func (x concernPolicy) Score(
	ctx context.Context,
	govAddr gov.GovOwnerAddress,
	govCloned gov.GovOwnerCloned,
	motion schema.Motion,
	instancePolicyNS ns.NS,

) schema.Score {

	return schema.Score{}
}

func (x concernPolicy) Close(
	ctx context.Context,
	govAddr gov.GovOwnerAddress,
	govCloned gov.GovOwnerCloned,
	motion schema.Motion,
	instancePolicyNS ns.NS,

) {

	// close the poll for the motion
	ballotName := schema.MotionPollBallotName(motion.ID)
	ballot.Close_StageOnly(
		ctx,
		govAddr,
		govCloned,
		ballotName,
		false,
	)

}
