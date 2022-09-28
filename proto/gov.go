package proto

type GovConfig struct {
	CommunityURL string `json:"community_url"`
	AdminURL     string `json:"admin_url"`
}

type GovUserInfo struct {
	URL string `json:"url"` // url of user's public soul repository
}
