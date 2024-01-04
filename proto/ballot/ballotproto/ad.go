package ballotproto

import (
	"sort"

	"github.com/gov4git/gov4git/v2/proto/account"
	"github.com/gov4git/gov4git/v2/proto/gov"
	"github.com/gov4git/gov4git/v2/proto/member"
	"github.com/gov4git/gov4git/v2/proto/motion"
	"github.com/gov4git/gov4git/v2/proto/purpose"
	"github.com/gov4git/lib4git/git"
)

type Advertisement struct {
	Gov          gov.Address       `json:"community"`
	ID           BallotID          `json:"id"`
	Owner        account.AccountID `json:"owner"`
	Purpose      purpose.Purpose   `json:"purpose"`
	MotionPolicy motion.PolicyName `json:"motion_policy"`
	//
	Title       string `json:"title"`
	Description string `json:"description"`
	//
	Choices      []string     `json:"choices"`
	Policy       PolicyName   `json:"policy"`
	Participants member.Group `json:"participants_group"`
	//
	Frozen    bool `json:"frozen"` // if frozen, the ballot is not accepting votes
	Closed    bool `json:"closed"` // closed ballots cannot be re-opened
	Cancelled bool `json:"cancelled"`
	//
	ParentCommit git.CommitHash `json:"parent_commit"`
}

type Advertisements []Advertisement

func (x Advertisements) Len() int {
	return len(x)
}

func (x Advertisements) Less(i, j int) bool {
	return x[i].ID.GitPath() < x[j].ID.GitPath()
}

func (x Advertisements) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x Advertisements) Sort() {
	sort.Sort(x)
}
