package bureauproto

type Request struct {
	Identify *IdentifyRequest `json:"identify"`
}

type Response struct {
	Identify *IdentifyResponse `json:"identify"`
}
