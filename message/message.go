package message

import (
	"github.com/ezequielpilcsuk/goMunication/member"
	"github.com/google/uuid"
)

type Message struct {
	Sender member.Member
	Data   []byte
	ID     uuid.UUID
}
