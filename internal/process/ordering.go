package process

import (
	"encoding/json"
	"fmt"
	"log"
	"mp1_node/internal/util"
	"sort"
)

var maxPriority = 0

func (p *Process) Ordering() {
	que := MsgQ{}
	for msg := range p.verified {
		rawbytes, err := json.Marshal(msg)
		if err != nil {
			log.Fatalf("JSON marshalling failed: %v\n", err)
		}
		fmt.Printf(Green+"[PROCESSING %s] %s\n"+Reset, msg.From, string(rawbytes))
		if msg.MT == Normal {
			fmt.Println(Blue, "NORMAL MESSAGE", Reset)
			N := len(que)
			if N != 0 {
				maxPriority = que[N-1].msg.Priority
			} else {
				maxPriority = 0
			}
			p.unicast(&Msg{From: p.self.Id, Id: msg.Id, MT: PrpPriority, Priority: maxPriority + 1}, p.peers[msg.From])
			que = append(que, p.msgs[msg.Id])
			sort.Sort(que)
			msg.Tx.TS = Undeliverable
		} else if msg.MT == PrpPriority {
			fmt.Println(Blue, "PROPOSED PRIORITY", Reset)
			if !p.contains(msg.Id) {
				fmt.Printf("Proposed priority for an unexisted message\n")
				continue
			}
			p.msgs[msg.Id].msg.Priority = util.Max(p.msgs[msg.Id].msg.Priority, msg.Priority)
			p.msgs[msg.Id].proposed++
			que.Print()
			if p.msgs[msg.Id].proposed == p.groupSize {
				p.multicast(&Msg{From: p.self.Id, Id: msg.Id, MT: AgrPriority, Priority: p.msgs[msg.Id].msg.Priority})
			}
		} else if msg.MT == AgrPriority {
			fmt.Println(Blue, "AGREED PRIORITY", Reset)
			p.msgs[msg.Id].msg.Priority = msg.Priority
			p.msgs[msg.Id].msg.Tx.TS = Deliverable
			sort.Sort(que)
			// que.Print()
			for que.Len() > 0 && que[0].msg.Tx.TS == Deliverable {
				fmt.Println("DELIVERING MESSAGE")
				toDeliver := que[0]
				que = que[1:]
				p.deliver(toDeliver.msg)
			}
		} else {
			fmt.Println("Invalid Message Type")
		}
	}
}
