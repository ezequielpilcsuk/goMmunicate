package group

import (
	"goMunication/member"
	"net"
)

const (
	ConnType = "tcp4"
)

// Group represents the group being communicated to
type Group struct {
	NMembers  int
	Members   []member.Member
	Address   net.Addr
	BasePort  int
	Sequencer member.Member
}

var thisGroup Group

// Start initializes a group
func (group *Group) Start(addr net.Addr, basePort int) {
	thisGroup.NMembers = 0
	thisGroup.Members = []member.Member{}

	if basePort == 0 {
		thisGroup.BasePort = 8180
	} else {
		thisGroup.BasePort = basePort
	}

	thisGroup.Address = addr
}

// Join adds a member to the group
func (group *Group) Join(member *member.Member) {
	member.Id = group.NMembers
	member.Port = group.BasePort + member.Id
	member.Group = group

	group.NMembers++
	group.Members = append(group.Members, *member)
}

// BDeliver is a basic
func (group *Group) BDeliver() {

}

/*
On initialization
	Received := {}

For process p to R-multicast message m to group g
	B-multicast(g,m)		p in g


On B-deliver(m) at process q with g = group(m)
	if(m not in Received)
		append(Received, m)
		if q not p
			B-multicast(g,m)
		R-deliver(m);
*/
