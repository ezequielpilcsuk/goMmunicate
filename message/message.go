package message

import (
	"goMunication/member"
)

type Message struct {
	Sender member.Member
	Data   []byte
	ID     int
}
