package notice

import (
	"fmt"
	"time"
)

type Notice struct {
	Body string `json:"body"`
}

func Noticef(format string, args ...any) Notices {
	return Notices{
		Notice{Body: fmt.Sprintf(format, args...)},
	}
}

type Notices []Notice

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

func (x *NoticeQueue) Append(notices ...Notice) {
	for _, notice := range notices {
		s := NoticeState{
			Stamp:     time.Now(),
			Notice:    notice,
			Displayed: false,
		}
		x.NoticeStates = append(x.NoticeStates, s)
	}
}
