package member

import (
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

func (me *Member) Receive() (message []byte) {
	tcpAddr, err := net.ResolveTCPAddr(ConnType, string(me.Port))
	utils.CheckErr(err)
	listener, err := net.ListenTCP(ConnType, tcpAddr)
	utils.CheckErr(err)
	conn, err := listener.Accept()
	conn.SetDeadline(time.Now().Add(time.Minute))

	_, err = conn.Read(message)
	utils.CheckErr(err)
	defer utils.CheckErr(conn.Close())

	return message
}

func Send(member Member, message []byte) {
	tcpAddr, err := net.ResolveTCPAddr(ConnType, string(member.Port))
	utils.CheckErr(err)

	conn, err := net.DialTCP(ConnType, nil, tcpAddr)
	utils.CheckErr(err)
	defer utils.CheckErr(conn.Close())

	conn.SetDeadline(time.Now().Add(time.Minute))

	_, err = conn.Write(message)
	utils.CheckErr(err)
}

func (me *Member) WrapMessage(data []byte) (message message.Message) {

	uuid
	return message
}
