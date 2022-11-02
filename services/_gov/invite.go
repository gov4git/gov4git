package gov

import (
	"bytes"
	"context"
	"fmt"
	"net/smtp"
	"text/template"

	"github.com/gov4git/gov4git/proto/cmdproto"
)

type InviteIn struct {
	SMTP            cmdproto.SMTPPlainAuth `json:"smtp_plain_auth"`
	From            cmdproto.EmailAddress  `json:"from"`
	To              cmdproto.EmailAddress  `json:"to"`
	CommunityURL    string                 `json:"community_url"`
	CommunityBranch string                 `json:"community_branch"`
}

type InviteOut struct {
	In *InviteIn `json:"in"`
}

func (x GovService) Invite(ctx context.Context, in *InviteIn) (*InviteOut, error) {
	auth := smtp.PlainAuth(in.SMTP.Identity, in.SMTP.Username, in.SMTP.Password, in.SMTP.Host)

	var msg bytes.Buffer

	fmt.Fprintf(&msg, "From: %v\n", in.From.Address)
	fmt.Fprintf(&msg, "To: %v\n", in.To.Address)
	fmt.Fprintf(&msg, "Subject: You are invited to join the git community %v\n\n", in.CommunityURL)

	if err := template.Must(template.New("").Parse(inviteBody)).Execute(&msg, in); err != nil {
		return nil, err
	}

	err := smtp.SendMail(in.SMTP.Host+":"+in.SMTP.Port, auth, in.From.String(), []string{in.To.String()}, []byte(msg.String()))
	if err != nil {
		return nil, err
	}

	return &InviteOut{In: in}, nil
}

const inviteBody = `
Hi {{.To.Name}},

You are invited to join the online git community at {{.CommunityURL}}.

To participate in community governance (polls, referendums, and so on),
you need to install the open-source tool gov4git from https://github.com/gov4git/gov4git

Gov4git is a protocol for transparent and accountable, decentralized community governance, based entirely on git.

Governance is viewed as a state machine, whose state is persisted in
the git repo of the community it applies to, in a human-readable form.

Governance operations (liking polling, calling referendums, making proposals, voting, issuing rewards to members)
are persisted on git such that their authenticity can be verified.

Cheers!
`
