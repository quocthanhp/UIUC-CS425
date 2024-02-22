package process

import (
	"encoding/json"
	"fmt"
)

// CANNOT MULTICAST TO SELF ????

func (p *Process) multicast(msg *Msg) {
	fmt.Println("start multicasting...")
		
	for _, peer := range p.peers {
		if peer.Id == p.self.Id {
			continue
		}
		fmt.Printf("send data to %s\n", peer.Id)
		jsonData, _ := json.Marshal(msg)
		peer.Conn.Write(append(jsonData, '\n'))
	}
}
