package process

import (
	"fmt"
	"mp1_node/internal/util"
	"sort"
)

var maxPriority = 0

func (p *Process) DemoOrdering() {
	for msg := range p.verified {
		fmt.Printf(Green+"[PROCESSING %s] %s\n"+Reset, msg.From, msg.toString())
	}
}

func (p *Process) Ordering() {
	for msg := range p.verified {
		p.Lock()
		if _, ok := p.peers[msg.From]; !ok {
			fmt.Println(Cyan, "Mesage from invalid source.", Reset)
			continue
		}
		fmt.Printf(Green+"[PROCESSING %s] %s\n"+Reset, msg.From, msg.toString())
		if msg.MT == Normal {
			fmt.Println(Blue, "NORMAL MESSAGE", Reset)
			N := p.que.Len()
			if N != 0 {
				maxPriority = p.que[N-1].msg.Priority
			} else {
				maxPriority = 0
			}
			p.unicast(&Msg{From: p.self.Id, Id: msg.Id, MT: PrpPriority, Priority: maxPriority + 1}, p.peers[msg.From])
			p.que = append(p.que, p.msgs[msg.Id])
			sort.Sort(p.que)
			msg.Tx.TS = Undeliverable
		} else if msg.MT == PrpPriority {
			fmt.Println(Blue, "PROPOSED PRIORITY", Reset)
			if !p.contains(msg.Id) {
				fmt.Printf("Proposed priority for an unexisted message\n")
				continue
			}
			p.msgs[msg.Id].msg.Priority = util.Max(p.msgs[msg.Id].msg.Priority, msg.Priority)
			p.msgs[msg.Id].proposed++
			p.que.Print()
			if p.msgs[msg.Id].proposed == p.groupSize {
				p.multicast(&Msg{From: p.self.Id, Id: msg.Id, MT: AgrPriority, Priority: p.msgs[msg.Id].msg.Priority})
			}
		} else if msg.MT == AgrPriority {
			fmt.Println(Blue, "AGREED PRIORITY", Reset)
			p.msgs[msg.Id].msg.Priority = msg.Priority
			p.msgs[msg.Id].msg.Tx.TS = Deliverable
			sort.Sort(p.que)
			// que.Print()
			for p.que.Len() > 0 && p.que[0].msg.Tx.TS == Deliverable {
				fmt.Println("DELIVERING MESSAGE")
				toDeliver := p.que[0]
				p.que = p.que[1:]
				p.deliver(toDeliver.msg)
			}
		} else {
			fmt.Println("Invalid Message Type")
		}
		p.Unlock()
	}
}
