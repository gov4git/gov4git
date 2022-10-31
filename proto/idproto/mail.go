package idproto

import (
	"crypto/sha256"
	"encoding/base64"
	"path/filepath"
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
	return filepath.Join(IdentityRoot, "mail", "response", stringHash(string(senderID)), stringHash(topic))
}

func SendMailTopicDirpath(receiverID ID, topic string) string {
	return filepath.Join(IdentityRoot, "mail", "request", stringHash(string(receiverID)), stringHash(topic))
}

func stringHash(s string) string {
	h := sha256.New()
	if _, err := h.Write([]byte(s)); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
