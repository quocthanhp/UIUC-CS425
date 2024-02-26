package process

import "sort"

var maxPriority = -1

// write code to ingest all the messages
func (p *Process) Ordering() {
	que := MsgQ{}
	for msg := range p.recvd {
		if true {
			// TODO: if is a normal message
			N := len(que)
			maxPriority = que[N-1].Priority
			// TODO: multicast proposed priority = maxPriority + 1
			que = append(que, msg)
			sort.Sort(que)
		} else if true {
			// TODO: if is a proposed priority
			// TODO: need a place to keep track of whether get all the proposed priority
		} else if true {
			// TODO: if is a agreed priority
			p.msgs[msg.Id].Priority = msg.Priority
			p.msgs[msg.Id].Tx.TS = Deliverable
			sort.Sort(que)
		} else {
			//TODO: handle invalid messages
		}
		if que[0].Tx.TS == Deliverable {
			//TODO: deliver the message
		}

	}
}
