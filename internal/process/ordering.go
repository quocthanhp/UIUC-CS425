package process

import "sort"

var maxPriority = -1

// write code to ingest all the messages
//TODO: search for all the places called function multicast
// and add the message type parameter
func (p *Process) Ordering() {
	que := MsgQ{}
	for msg := range p.recvd {
		if msg.MT == Normal {
			N := len(que)
			maxPriority = que[N-1].Priority
			// TODO: multicast proposed priority = maxPriority + 1
			que = append(que, msg)
			sort.Sort(que)
			msg.Tx.TS = Undeliverable
		} else if msg.MT == PrpPriority {
			// TODO: need a place to keep track of whether get all the proposed priority
			// TODO: if get all the proposed priority from all the processes in teh group
			// choose the agreed priority adn send it out
		} else if msg.MT == AgrPriority {
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
