package concern

import (
	"context"
	"fmt"

	"github.com/gov4git/gov4git/proto/ballot/ballot"
	"github.com/gov4git/gov4git/proto/ballot/qv"
	"github.com/gov4git/gov4git/proto/docket/policies/pmp"
	"github.com/gov4git/gov4git/proto/docket/policy"
	"github.com/gov4git/gov4git/proto/docket/schema"
	"github.com/gov4git/gov4git/proto/gov"
	"github.com/gov4git/gov4git/proto/member"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/ns"
)

func init() {
	policy.Install(context.Background(), concernPolicy{})
}

const ConcernPolicyName = schema.PolicyName("pmp-concern-policy")

type concernPolicy struct{}

func (x concernPolicy) Name() string {
	return ConcernPolicyName.String()
}

func (x concernPolicy) Open(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	instancePolicyNS ns.NS,

) {

	// open a poll for the motion
	ballotName := pmp.MotionPollBallotName(motion.ID)
	ballot.Open_StageOnly(
		ctx,
		qv.QV{},
		cloned,
		ballotName,
		fmt.Sprintf("Priority poll for motion %v", motion.ID),
		fmt.Sprintf("Up/down vote the priority of motion %v", motion.ID),
		[]string{schema.MotionPollBallotChoice},
		member.Everybody,
	)

}

func (x concernPolicy) Score(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	instancePolicyNS ns.NS,

) schema.Score {

	return schema.Score{}
}

func (x concernPolicy) Close(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	instancePolicyNS ns.NS,

) {

	// close the poll for the motion
	ballotName := pmp.MotionPollBallotName(motion.ID)
	ballot.Close_StageOnly(
		ctx,
		cloned,
		ballotName,
		false,
	)

}

func (x concernPolicy) Cancel(
	ctx context.Context,
	cloned gov.OwnerCloned,
	motion schema.Motion,
	instancePolicyNS ns.NS,

) {

}

func (x concernPolicy) Show(
	ctx context.Context,
	cloned gov.Cloned,
	motion schema.Motion,
	instancePolicyNS ns.NS,

) form.Map {

	return nil
}

func (x concernPolicy) AddRefTo(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
) {

}

func (x concernPolicy) AddRefFrom(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
) {

}

func (x concernPolicy) RemoveRefTo(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
) {

}

func (x concernPolicy) RemoveRefFrom(
	ctx context.Context,
	cloned gov.OwnerCloned,
	refType schema.RefType,
	from schema.Motion,
	to schema.Motion,
	fromPolicyNS ns.NS,
	toPolicyNS ns.NS,
) {

}
