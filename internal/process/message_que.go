package process

// struct for a pending message in the queue
type PdMsg struct {
	msg      *Msg
	proposed int
}

type MsgQ []PdMsg

func (q MsgQ) Len() int           { return len(q) }
func (q MsgQ) Less(i, j int) bool { return q[i].msg.Priority < q[j].msg.Priority }
func (q MsgQ) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }
