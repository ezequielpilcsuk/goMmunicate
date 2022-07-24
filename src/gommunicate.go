package src

import (
	"encoding/json"
	"github.com/google/uuid"
	"goMunication/utils"
	"io/ioutil"
	"log"
	"net"
	"time"
)

const (
	ConnType = "tcp4"
)

type Message struct {
	SenderID int
	Data     []byte
	ID       uuid.UUID
}

// Member is part of the group
type Member struct {
	Id        int `json:"id"`
	Port      int `json:"port"`
	NMembers  int `json:"n_members"`
	LClocks   []int
	MReceived map[uuid.UUID]bool
}

// Group represents the group being communicated to
type Group struct {
	NMembers    int `json:"n_members"`
	MembersIDs  []int
	Address     string `json:"adress"`
	BasePort    int    `json:"base_port"`
	SequencerID int    `json:"sequencer_id"`
}

var ThisMember Member
var OurGroup Group

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

	//TODO: set this dynamically

	OurGroup.NMembers = ThisMember.NMembers
	OurGroup.MembersIDs = []int{0, 1, 2}
	OurGroup.Address = "localhost"
	OurGroup.BasePort = 8080

	//TODO: fix when adding total order
	OurGroup.SequencerID = -1
}

// Send a message to a specific member of the group
func Send(memberID int, data []byte) {
	finalAddr := string(OurGroup.BasePort + memberID)
	tcpAddr, err := net.ResolveTCPAddr(ConnType, finalAddr)
	utils.CheckErr(err)

	conn, err := net.DialTCP(ConnType, nil, tcpAddr)
	utils.CheckErr(err)
	defer utils.CheckErr(conn.Close())

	utils.CheckErr(conn.SetDeadline(time.Now().Add(time.Minute)))

	var msg = Message{
		SenderID: ThisMember.Id,
		Data:     data,
		ID:       uuid.New(),
	}

	//tmp := make([]byte, 5000)

	_, err = conn.Write([]byte(message))
	utils.CheckErr(err)
}

// Receive a message from the group
func (m *Member) Receive() (message Message) {
	tcpAddr, err := net.ResolveTCPAddr(ConnType, string(m.Port))
	utils.CheckErr(err)
	listener, err := net.ListenTCP(ConnType, tcpAddr)
	utils.CheckErr(err)
	conn, err := listener.Accept()

	utils.CheckErr(conn.SetDeadline(time.Now().Add(time.Minute)))

	tmp := make([]byte, 5000)
	_, err = conn.Read(tmp)

	tmpstruct := new(Message)

	utils.CheckErr(err)
	defer utils.CheckErr(conn.Close())

	return message
}

// bMulticast sends a message to the whole group
func (m Member) bMulticast(message Message) {
	for i := 0; i < OurGroup.NMembers; i++ {
		Send(OurGroup.MembersIDs[i], message)
	}
}

// BDeliver receives an unreceived message from the group
func BDeliver(m Message) {
	if !ThisMember.MReceived[m.ID] {
		ThisMember.MReceived[m.ID] = true
		if m.SenderID != ThisMember.Id {
			ThisMember.bMulticast(m)
		}
	}
}
