package process

import "fmt"

// struct for a pending message in the queue
type PdMsg struct {
	msg      *Msg
	proposed int
}

type MsgQ []*PdMsg

func (q MsgQ) Len() int { return len(q) }
func (q MsgQ) Less(i, j int) bool {
	if q[i].msg.Priority == q[j].msg.Priority {
		return q[i].msg.Id < q[j].msg.Id
	}
	return q[i].msg.Priority < q[j].msg.Priority
}
func (q MsgQ) Swap(i, j int) { q[i], q[j] = q[j], q[i] }

func (q MsgQ) Print() {
	fmt.Print(Red)
	fmt.Println("PRINTING QUEUE:")
	for _, pdmsg := range q {
		fmt.Printf("ID:%s\tFROM:%s\tPRP#:%d\tPRIORITY:%d\tTS:%d\t%d\n", pdmsg.msg.Id, pdmsg.msg.From, pdmsg.proposed, pdmsg.msg.Priority, pdmsg.msg.Tx.TS, pdmsg.msg.Tx.TS)
	}
	fmt.Print(Reset)
}
