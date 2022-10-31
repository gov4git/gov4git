package bureauproto

// const xTopic = "bureau"

// func Topic(receiverRepo string) string {
// 	return XXX // comty repo + bureau
// }

type Request struct {
	Identify *IdentifyRequest `json:"identify"`
}

type Response struct {
	Identify *IdentifyResponse `json:"identify"`
}
