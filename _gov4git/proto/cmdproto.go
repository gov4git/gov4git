package cmdproto

import "fmt"

const (
	LocalAgentPath     = ".gov4git"
	LocalAgentTempPath = "gov4git"
)

type Config struct {
	PublicURL       string         `json:"public_url"`
	PrivateURL      string         `json:"private_url"`
	CommunityURL    string         `json:"community_url"`
	CommunityBranch string         `json:"community_branch"`
	SMTPPlainAuth   *SMTPPlainAuth `json:"smtp_plain_auth"` // for sending invitations
	InviteEmailFrom *EmailAddress  `json:"invite_email_from"`
}

type SMTPPlainAuth struct {
	Identity string `json:"identity"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

type EmailAddress struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func (x EmailAddress) String() string {
	return fmt.Sprintf("%v <%v>", x.Name, x.Address)
}
