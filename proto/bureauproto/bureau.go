package bureauproto

const bureauTopic = "bureau"

func Topic(topic string) string {
	return bureauTopic + ":" + topic
}
