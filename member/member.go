package member

import (
	"github.com/google/uuid"
	"goMunication/group"
	"goMunication/message"
	"goMunication/utils"
	"net"
	"time"
)

const (
	ConnType = "tcp4"
)

// Member is part of the group
type Member struct {
	Id    int          `json:"id"`
	Port  int          `json:"port"`
	Group *group.Group `json:"group"`
	Role  string       `json:"role"`
}

// Receive a message from the group
func (me *Member) Receive() (message message.Message) {
	tcpAddr, err := net.ResolveTCPAddr(ConnType, string(me.Port))
	utils.CheckErr(err)
	listener, err := net.ListenTCP(ConnType, tcpAddr)
	utils.CheckErr(err)
	conn, err := listener.Accept()
	conn.SetDeadline(time.Now().Add(time.Minute))

	_, err = conn.Read(message.Data)
	utils.CheckErr(err)
	defer utils.CheckErr(conn.Close())

	return message
}

// Send a message to a specific member of the group
func Send(member Member, message message.Message) {
	tcpAddr, err := net.ResolveTCPAddr(ConnType, string(member.Port))
	utils.CheckErr(err)

	conn, err := net.DialTCP(ConnType, nil, tcpAddr)
	utils.CheckErr(err)
	defer utils.CheckErr(conn.Close())

	conn.SetDeadline(time.Now().Add(time.Minute))

	_, err = conn.Write(message.Data)
	utils.CheckErr(err)
}

// bMulticast sends a message to the whole group
func (member *Member) bMulticast(data []byte) {
	for i := 0; i < member.Group.NMembers; i++ {
		message := member.WrapMessage(data)
		Send(member.Group.Members[i], message)
	}
}

func (me Member) WrapMessage(data []byte) (message message.Message) {
	message.Sender = me
	message.Data = data
	message.ID = uuid.New()
	return message
}
