package bureau

import (
	"context"

	"github.com/gov4git/gov4git/lib/base"
	"github.com/gov4git/gov4git/lib/git"
	"github.com/gov4git/gov4git/proto/bureauproto"
	"github.com/gov4git/gov4git/proto/govproto"
	"github.com/gov4git/gov4git/proto/idproto"
	"github.com/gov4git/gov4git/services/gov/member"
	"github.com/gov4git/gov4git/services/gov/user"
	"github.com/gov4git/gov4git/services/id"
)

type BureauServer struct {
	GovConfig      govproto.GovConfig
	IdentityConfig idproto.IdentityConfig
}

func (x BureauServer) IdentityService() id.IdentityService {
	return id.IdentityService{IdentityConfig: x.IdentityConfig}
}

type RespondFunc func(
	ctx context.Context,
	community git.Local,
	user user.UserInfo,
	req []byte,
) (response []byte, err error)

func (x BureauServer) RespondGroup(
	ctx context.Context,
	branch git.Branch,
	group string,
	respond RespondFunc,
	topic string,
) error {
	community, err := git.CloneBranch(ctx, x.GovConfig.CommunityURL, string(branch))
	if err != nil {
		return err
	}
	err = x.RespondGroupLocal(ctx, community, group, respond, topic)
	if err != nil {
		return err
	}
	if err := community.PushUpstream(ctx); err != nil {
		return err
	}
	return nil
}

func (x BureauServer) RespondGroupLocal(
	ctx context.Context,
	community git.Local,
	group string,
	respond RespondFunc,
	topic string,
) error {
	// list users in the group
	memberService := member.GovMemberService{GovConfig: x.GovConfig}
	participants, err := memberService.ListLocal(ctx, community, "", group)
	if err != nil {
		return err
	}
	userInfos, err := user.GetInfos(ctx, community, member.ExtractUsersFromMembership(participants))
	if err != nil {
		return err
	}

	// respond to each user's requests
	for _, userInfo := range userInfos {
		err := x.RespondUserLocal(ctx, community, userInfo, respond, topic)
		if err != nil {
			base.Infof("responding to user %v (%v)", userInfo.UserName, err)
		}
	}

	return nil
}

func (x BureauServer) RespondUserLocal(
	ctx context.Context,
	community git.Local,
	user user.UserInfo,
	respond RespondFunc,
	topic string,
) error {

	// fetch user's repo
	userRepo, err := git.CloneBranch(ctx, user.UserInfo.PublicURL, idproto.IdentityBranch)
	if err != nil {
		return err
	}

	// respond to user's requests
	reply := func(ctx context.Context, req []byte) ([]byte, error) {
		return respond(ctx, community, user, req)
	}
	_, err = id.ReceiveRespondMailLocalStageOnly(ctx, community, userRepo, bureauproto.Topic(topic), reply)
	if err != nil {
		return err
	}

	// commit
	if err := community.Commitf(ctx, "Responded to user %v on topiv %v", user.UserName, bureauproto.Topic(topic)); err != nil {
		return err
	}

	return nil
}
