package schema

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

var (
	Accept Decision = "accept"
	Reject Decision = "reject"
)
