package process

type MsgQ []*Msg

func (q MsgQ) Len() int           { return len(q) }
func (q MsgQ) Less(i, j int) bool { return q[i].Priority < q[j].Priority }
func (q MsgQ) Swap(i, j int)      { q[i], q[j] = q[j], q[i] }
