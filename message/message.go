package message

import (
	"github.com/google/uuid"
	"goMunication/member"
)

type Message struct {
	Sender member.Member
	Data   []byte
	ID     uuid.UUID
}
