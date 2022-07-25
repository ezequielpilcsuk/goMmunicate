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
	SeqNum   int
	LClocks  []int
	ID       uuid.UUID
}

// Member is part of the group
type Member struct {
	Id   int `json:"id"`
	Port int `json:"port"`

	SeqNum      int       // SeqNum is message sequence counter, sequencer exclusive
	NextDeliver int       // NextDeliver is the sequence number for the next message to be received
	Pending     []Message // Pending is a list of messages to be received

	LClocks   []int              //LClocks is a vector of logical clocks for Casual Order
	MReceived map[uuid.UUID]bool // MReceived is a vector of received messages
}

// Group represents the group being communicated to
type Group struct {
	NMembers    int `json:"n_members"`
	MembersIDs  []int
	Address     string `json:"address"`
	BasePort    int    `json:"base_port"`
	SequencerID int    `json:"sequencer_id"`
}

var ThisMember Member
var OurGroup Group

// TODO: review structure, maybe Start is not needed and should be in another program

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
	OurGroup.NMembers = 3
	OurGroup.MembersIDs = []int{0, 1, 2}
	OurGroup.Address = "localhost"
	OurGroup.BasePort = 8080

	//TODO: fix when adding total order
	OurGroup.SequencerID = 0
	if OurGroup.SequencerID == ThisMember.Id {
		ThisMember.SeqNum = 1
	}

	ThisMember.NextDeliver = 1

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

	_, err = conn.Write(utils.UnwrapMessage(msg))
	utils.CheckErr(err)

	ThisMember.LClocks[ThisMember.Id]++
}

// Receive a message from the group
func Receive() (messages []Message) {
	tcpAddr, err := net.ResolveTCPAddr(ConnType, string(ThisMember.Port))
	utils.CheckErr(err)
	listener, err := net.ListenTCP(ConnType, tcpAddr)
	utils.CheckErr(err)
	conn, err := listener.Accept()
	defer utils.CheckErr(conn.Close())
	utils.CheckErr(conn.SetDeadline(time.Now().Add(time.Minute)))

	tmp := make([]byte, 5000)
	_, err = conn.Read(tmp)
	utils.CheckErr(err)

	message := utils.WrapMessage(tmp)

	// If this member is the sequencer and not the sender
	if ThisMember.Id == OurGroup.SequencerID && message.SenderID != ThisMember.Id {
		SequencerReceive(message)
		return
	}

	ThisMember.Pending = append(ThisMember.Pending, message)

	//TODO: fix for later

	// While there are pending messages to be received
	for i := 0; i < len(ThisMember.Pending); i++ {
		if ThisMember.Pending[i].SeqNum == ThisMember.NextDeliver {
			currMsg := ThisMember.Pending[i]

			messages = append(messages, currMsg)
			ThisMember.NextDeliver++
			ThisMember.LClocks[ThisMember.Id]++

			// Updating Logical clocks on receive
			ThisMember.LClocks = utils.UpdateLClocks(ThisMember.LClocks, currMsg.LClocks)
		}
	}

	return messages
}

// SequencerReceive assigns a sequence number to a message and broadcasts it to the group
func SequencerReceive(message Message) {
	message.SeqNum = ThisMember.SeqNum
	BMulticast(message)
	ThisMember.SeqNum++
}

// Broadcast Sends a message to the sequencer to be broadcasted
func Broadcast(message Message) {
	finalAddr := string(OurGroup.BasePort + OurGroup.SequencerID)
	tcpAddr, err := net.ResolveTCPAddr(ConnType, finalAddr)
	utils.CheckErr(err)

	conn, err := net.DialTCP(ConnType, nil, tcpAddr)
	utils.CheckErr(err)
	defer utils.CheckErr(conn.Close())

	utils.CheckErr(conn.SetDeadline(time.Now().Add(time.Minute)))

	_, err = conn.Write(utils.UnwrapMessage(message))
	utils.CheckErr(err)
	ThisMember.LClocks[ThisMember.Id]++
}

// BMulticast sends a message to the whole group
func BMulticast(message Message) {
	for i := 0; i < OurGroup.NMembers; i++ {
		if ThisMember.Id != OurGroup.MembersIDs[i] {
			Send(OurGroup.MembersIDs[i], utils.UnwrapMessage(message))
		}
	}
}

// BDeliver receives an unreceived message from the group
func BDeliver(m Message) {
	if !ThisMember.MReceived[m.ID] {
		ThisMember.MReceived[m.ID] = true
		if m.SenderID != ThisMember.Id {
			Broadcast(m)
		}
	}
}
