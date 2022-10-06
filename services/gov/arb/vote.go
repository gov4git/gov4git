package arb

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/petar/gov4git/lib/files"
	"github.com/petar/gov4git/lib/form"
	"github.com/petar/gov4git/lib/git"
	"github.com/petar/gov4git/proto"
	"github.com/petar/gov4git/services/identity"
)

type VoteIn struct {
	ReferendumBranch string  `json:"referendum_branch"`
	VoteChoice       string  `json:"vote_choice"`
	VoteStrength     float64 `json:"vote_strength"`
}

type VoteOut struct {
	VoteRepo         string `json:"vote_repo"`
	VoteBranch       string `json:"vote_branch"`
	ReferendumRepo   string `json:"referendum_repo"`
	ReferendumBranch string `json:"referendum_branch"`
}

func (x VoteOut) Human(context.Context) string {
	return fmt.Sprintf("Vote placed in repo %v at branch %v.\n"+
		"Regarding referendum in repo %v and branch %v",
		x.VoteRepo, x.VoteBranch,
		x.ReferendumRepo, x.ReferendumBranch,
	)
}

func (x GovArbService) Vote(ctx context.Context, in *VoteIn) (*VoteOut, error) {
	// find poll advertisement in community repo

	// clone community repo at referendum branch locally
	community := git.LocalFromDir(files.WorkDir(ctx).Subdir("community"))
	if err := community.CloneBranch(ctx, x.GovConfig.CommunityURL, in.ReferendumBranch); err != nil {
		return nil, err
	}

	// find the genesis commit of the referendum
	findGenesis, err := x.FindPollGenesisLocal(ctx, community, &FindPollGenesisIn{PollBranch: in.ReferendumBranch})
	if err != nil {
		return nil, err
	}

	// read the poll advertisement
	pollPath, err := proto.PollPathFromBranch(in.ReferendumBranch)
	if err != nil {
		return nil, err
	}
	pollAdPath := filepath.Join(pollPath, proto.GovPollAdFilebase)
	if err := community.CheckoutBranch(ctx, findGenesis.PollGenesisCommit); err != nil {
		return nil, err
	}
	var pollAd proto.GovPollAd
	if _, err := community.Dir().ReadFormFile(ctx, pollAdPath, &pollAd); err != nil {
		return nil, err
	}

	// cast vote in voters public identity repo

	// clone voter identity public repo at the main identity branch
	voter := git.LocalFromDir(files.WorkDir(ctx).Subdir("voter"))
	if err := voter.CloneBranch(ctx, x.IdentityConfig.PublicURL, proto.IdentityBranch); err != nil {
		return nil, err
	}

	// retrieve the voter's private keys from the private identity repo
	idService := identity.IdentityService{IdentityConfig: x.IdentityConfig}
	voterCredentials, err := idService.GetPrivateCredentials(ctx, &identity.GetPrivateCredentialsIn{})
	if err != nil {
		return nil, err
	}

	// compute the name of the vote branch
	voteBranch, err := proto.PollVoteBranch(ctx, pollAd)
	if err != nil {
		return nil, err
	}

	// checkout an existing voting branch or create an orphan one
	if err := voter.CheckoutBranch(ctx, voteBranch); err != nil {
		if err := voter.CheckoutNewOrphan(ctx, voteBranch); err != nil {
			return nil, err
		}
		if err := voter.ResetHard(ctx); err != nil {
			return nil, err
		}
	}

	// add vote to vote branch
	vote := proto.GovPollVote{
		PollAd:   pollAd,
		Choice:   in.VoteChoice,
		Strength: in.VoteStrength,
	}
	// sign vote
	voteData, err := form.EncodeForm(ctx, vote)
	if err != nil {
		return nil, err
	}
	signature, err := proto.SignPlaintext(ctx, &voterCredentials.PrivateCredentials, voteData)
	if err != nil {
		return nil, err
	}
	signatureData, err := form.EncodeForm(ctx, signature)
	if err != nil {
		return nil, err
	}

	// write vote and signature
	stage := files.ByteFiles{
		files.ByteFile{Path: proto.GovPollVoteFilepath, Bytes: voteData},
		files.ByteFile{Path: proto.GovPollVoteSignatureFilepath, Bytes: signatureData},
	}
	if err := voter.Dir().WriteByteFiles(stage); err != nil {
		return nil, err
	}
	if err := voter.Add(ctx, stage.Paths()); err != nil {
		return nil, err
	}
	msg := proto.PollVoteCommitHeader(x.GovConfig.CommunityURL, in.ReferendumBranch)
	if err := voter.Commit(ctx, msg); err != nil {
		return nil, err
	}

	// push identity repo to origin
	if err := voter.PushUpstream(ctx); err != nil {
		return nil, err
	}

	return &VoteOut{
		VoteRepo:         x.IdentityConfig.PublicURL,
		VoteBranch:       voteBranch,
		ReferendumRepo:   x.GovConfig.CommunityURL,
		ReferendumBranch: in.ReferendumBranch,
	}, nil
}
