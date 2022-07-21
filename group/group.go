package group

// Group represents the group being communicated to
type Group struct {
	NMembers    int `json:"n_members"`
	MembersIDs  []int
	Address     string `json:"adress"`
	BasePort    int    `json:"base_port"`
	SequencerID int    `json:"sequencer_id"`
}
