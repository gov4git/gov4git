package notice

import "time"

type Notice struct {
	Body string `json:"body"`
}

type NoticeState struct {
	Stamp     time.Time `json:"stamp"`
	Notice    Notice    `json:"notice"`
	Displayed bool      `json:"shown"`
}

func (x *NoticeState) IsDisplayed() bool {
	return x.Displayed
}

func (x *NoticeState) SetDisplayed() {
	x.Displayed = true
}

type NoticeQueue struct {
	NoticeStates []NoticeState `json:"notices"`
}

func NewNoticeQueue() *NoticeQueue {
	return &NoticeQueue{}
}

func (x *NoticeQueue) Push(notice Notice) {
	s := NoticeState{
		Stamp:     time.Now(),
		Notice:    notice,
		Displayed: false,
	}
	x.NoticeStates = append(x.NoticeStates, s)
}
