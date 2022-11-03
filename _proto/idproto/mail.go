package idproto

import (
	"path/filepath"

	"github.com/gov4git/gov4git/lib/form"
)

type SendBoxInfo struct {
	ReceiverID ID     `json:"receiver_id"`
	Topic      string `json:"topic"`
}

type ReceiveBoxInfo struct {
	SenderID ID     `json:"sender_id"`
	Topic    string `json:"topic"`
}

const (
	BoxInfoFilebase = "info.json"
	NextFilebase    = "next"
)

func ReceiveMailTopicDirpath(senderID ID, topic string) string {
	return filepath.Join(IdentityRoot, "mail", "received", form.StringHashForFilename(string(senderID)), form.StringHashForFilename(topic))
}

func SendMailTopicDirpath(receiverID ID, topic string) string {
	return filepath.Join(IdentityRoot, "mail", "sent", form.StringHashForFilename(string(receiverID)), form.StringHashForFilename(topic))
}
