package motionproto

import "github.com/gov4git/gov4git/v2/proto/history/metric"

type PolicyName string

func (x PolicyName) String() string {
	return string(x)
}

type Decision string

func (x Decision) String() string {
	return string(x)
}

func (x Decision) IsAccept() bool {
	return x == Accept
}

func (x Decision) IsReject() bool {
	return x == Reject
}

func (x Decision) MetricDecision() metric.MotionDecision {
	return metric.MotionDecision(x)
}

var (
	Accept Decision = "accept"
	Reject Decision = "reject"
)
