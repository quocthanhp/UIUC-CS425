package process

import (
	"encoding/json"
	"fmt"
)

func (p *Process) multicast(msg *Msg) {
	fmt.Println(Purple, "start multicasting...", Reset)

	for _, peer := range p.peers {
		p.unicast(msg, peer)
	}
	fmt.Println(Purple, "Finished mulitcast.", Reset)
}

func (p *Process) unicast(msg *Msg, peer *Node) {
	fmt.Printf(Purple+"Unicasting to %s...."+Reset+"\n", peer.Id)
	if peer.Id == p.self.Id {
		p.recvd <- msg	
		return
	}

	jsonData, _ := json.Marshal(msg)
	toSend := append(jsonData, '\n')
	_, err := peer.Conn.Write(toSend)
	if err != nil {
		fmt.Println("ERROR:")
		fmt.Println(err.Error())
	}
}
