package process

import (
	"fmt"
	"mp1_node/internal/util"
	"sort"
)

var maxPriority = -1

// write code to ingest all the messages
//TODO: search for all the places called function multicast
// and add the message type parameter
func (p *Process) Ordering() {
	que := MsgQ{}
	for msg := range p.recvd {
		if msg.MT == Normal {
			N := len(que)
			maxPriority = que[N-1].msg.Priority
			p.unicast(&Msg{From: p.self.Id, Id: msg.Id, MT: PrpPriority, Priority: maxPriority + 1}, p.peers[msg.From])
			que = append(que, PdMsg{msg, 0})
			sort.Sort(que)
			msg.Tx.TS = Undeliverable
		} else if msg.MT == PrpPriority {
			if !p.contains(msg.Id) {
				fmt.Printf("Proposed priority for an unexisted message\n")
				continue
			}
			p.msgs[msg.Id].msg.Priority = util.Max(p.msgs[msg.Id].msg.Priority, msg.Priority)
			p.msgs[msg.Id].proposed++
			if p.msgs[msg.Id].proposed == p.groupSize {
				// TODO: multicast agreed priority
				p.multicast(&Msg{From: p.self.Id, Id: msg.Id, MT: AgrPriority, Priority: p.msgs[msg.Id].msg.Priority})
			}
		} else if msg.MT == AgrPriority {
			p.msgs[msg.Id].msg.Priority = msg.Priority
			p.msgs[msg.Id].msg.Tx.TS = Deliverable
			sort.Sort(que)
			for que[0].msg.Tx.TS == Deliverable {
				toDeliver := que[0]
				que = que[1:]
				p.deliver(toDeliver.msg)
			}
		} else {
			fmt.Println("Invalid Message Type")
		}
	}
}
