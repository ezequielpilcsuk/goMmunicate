package member

import (
	"encoding/json"
	"github.com/google/uuid"
	"goMunication/group"
	"goMunication/message"
	"goMunication/utils"
	"io/ioutil"
	"log"
	"net"
	"time"
)

const (
	ConnType = "tcp4"
)

// Member is part of the group
type Member struct {
	Id       int `json:"id"`
	Port     int `json:"port"`
	NMembers int `json:"n_members"`
}

var ThisMember Member
var OurGroup group.Group

// Start should initialize ThisMember and OurGroup loading data from dist/config.json
func Start() {
	var err error
	var byteValue []byte

	localConfigFilePath := "/dist/config.json"
	byteValue, err = ioutil.ReadFile(localConfigFilePath)
	if err != nil {
		log.Println("init fail: no local config file found: ", err)
	}

	if err := json.Unmarshal(byteValue, &ThisMember); err != nil {
		log.Println("init fail: wrong json format: ", err)
	}

	//TODO: set this dinamically

	OurGroup.NMembers = ThisMember.NMembers
	OurGroup.MembersIDs = []int{0, 1, 2}
	OurGroup.Address = "localhost"
	OurGroup.BasePort = 8080

	//TODO: fix when adding total order
	OurGroup.SequencerID = -1
}

// Send a message to a specific member of the group
func Send(memberID int, message message.Message) {
	finalAddr := string(OurGroup.BasePort + memberID)
	tcpAddr, err := net.ResolveTCPAddr(ConnType, finalAddr)
	utils.CheckErr(err)

	conn, err := net.DialTCP(ConnType, nil, tcpAddr)
	utils.CheckErr(err)
	defer utils.CheckErr(conn.Close())

	conn.SetDeadline(time.Now().Add(time.Minute))

	_, err = conn.Write(message.Data)
	utils.CheckErr(err)
}

// Receive a message from the group
func (m *Member) Receive() (message message.Message) {
	tcpAddr, err := net.ResolveTCPAddr(ConnType, string(m.Port))
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

// bMulticast sends a message to the whole group
func (m Member) bMulticast(message message.Message) {
	for i := 0; i < OurGroup.NMembers; i++ {
		Send(OurGroup.MembersIDs[i], message)
	}
}

// BDeliver is a basic
func (m Member) BDeliver() (message message.Message) {
	return m.Receive()
}

//
func bDeliver(m message.Message) {
	received := map[uuid.UUID]bool{}
	if !received[m.ID] {
		received[m.ID] = true
		if m.Sender.Id != ThisMember.Id {
			ThisMember.bMulticast(m)
		}
	}
}