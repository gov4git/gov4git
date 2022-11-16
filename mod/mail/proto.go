package mail

import (
	"path/filepath"

	"github.com/gov4git/gov4git/mod/id"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/ns"
)

type SendBoxInfo struct {
	ReceiverID id.ID  `json:"receiver_id"`
	Topic      string `json:"topic"`
}

type ReceiveBoxInfo struct {
	SenderID id.ID  `json:"sender_id"`
	Topic    string `json:"topic"`
}

const (
	BoxInfoFilebase = "box_info.json"
	NextFilebase    = "next.json"
)

var SendNS = id.PublicNS.Sub("mail/sent")
var ReceiveNS = id.PublicNS.Sub("mail/received")

func ReceiveTopicNS(senderID id.ID, topic string) ns.NS {
	return ReceiveNS.Sub(
		filepath.Join(
			form.StringHashForFilename(string(senderID)),
			form.StringHashForFilename(topic),
		),
	)
}

func SendTopicNS(receiverID id.ID, topic string) ns.NS {
	return SendNS.Sub(
		filepath.Join(
			form.StringHashForFilename(string(receiverID)),
			form.StringHashForFilename(topic),
		),
	)
}
