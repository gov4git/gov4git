package bureauproto

const bureauTopic = "bureau"

func Topic(topic string) string {
	return bureauTopic + ":" + topic
}

type Request struct {
	Identify *IdentifyRequest `json:"identify"`
}

type Response struct {
	Identify *IdentifyResponse `json:"identify"`
}
