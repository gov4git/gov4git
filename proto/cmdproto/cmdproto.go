package cmdproto

const (
	LocalAgentPath     = ".gov4git"
	LocalAgentTempPath = "gov4git"
)

type Config struct {
	PublicURL       string `json:"public_url"`
	PrivateURL      string `json:"private_url"`
	CommunityURL    string `json:"community_url"`
	CommunityBranch string `json:"community_branch"`
}
