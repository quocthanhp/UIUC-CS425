package process

import (
	"encoding/json"
	"fmt"
)

func (p *Process) multicast(msg *Msg) {
	fmt.Println("start multicasting...")

	for _, peer := range p.peers {
		if peer.Id == p.self.Id {
			continue
		}
		jsonData, _ := json.Marshal(msg)
		toSend := append(jsonData, '\n')
		_, err := peer.Conn.Write(toSend)
		if err != nil {
			fmt.Println("ERROR:")
			fmt.Println(err.Error())
		} else {
			fmt.Println("ERROR FREE!")
		}
	}
}
