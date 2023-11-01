package mail

import (
	"sort"

	"github.com/gov4git/gov4git/proto/id"
	"github.com/gov4git/lib4git/form"
	"github.com/gov4git/lib4git/ns"
)

type SendBoxInfo struct {
	// ReceiverAddr id.PublicAddress     `json:"receiver_address"`
	ReceiverCred id.PublicCredentials `json:"receiver_credentials"`
	Topic        string               `json:"topic"`
}

type ReceiveBoxInfo struct {
	// SenderAddr id.PublicAddress     `json:"sender_address"`
	SenderCred id.PublicCredentials `json:"sender_credentials"`
	Topic      string               `json:"topic"`
}

type SeqNo int64

type SentMsg[Msg form.Form] struct {
	SeqNo SeqNo `json:"seqno"`
	Msg   Msg   `json:"msg"`
}

type SentMsgs[Msg form.Form] []SentMsg[Msg]

func (x SentMsgs[Msg]) Sort()              { sort.Sort(x) }
func (x SentMsgs[Msg]) Len() int           { return len(x) }
func (x SentMsgs[Msg]) Less(i, j int) bool { return x[i].SeqNo < x[j].SeqNo }
func (x SentMsgs[Msg]) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

type RequestEnvelope[Req form.Form] struct {
	SeqNo   SeqNo `json:"seqno"`
	Request Req   `json:"request"`
}

type ResponseEnvelope[Resp form.Form] struct {
	SeqNo    SeqNo `json:"seqno"`
	Response Resp  `json:"response"`
}

const (
	BoxInfoFilebase = "box_info.json"
	NextFilebase    = "next.json"
)

var SendNS = id.PublicNS.Append("mail/sent")
var ReceiveNS = id.PublicNS.Append("mail/received")

func ReceiveTopicNS(senderID id.ID, topic string) ns.NS {
	return ReceiveNS.Append(
		form.StringHashForFilename(string(senderID)),
		form.StringHashForFilename(topic),
	)
}

func SendTopicNS(receiverID id.ID, topic string) ns.NS {
	return SendNS.Append(
		form.StringHashForFilename(string(receiverID)),
		form.StringHashForFilename(topic),
	)
}
