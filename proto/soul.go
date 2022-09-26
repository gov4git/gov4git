package proto

type SoulID string // hash of soul public key

type PublicSoulAddress struct {
	PublicURL string
}

type PrivateSoulAddress struct {
	PublicURL  string
	PrivateURL string
}
