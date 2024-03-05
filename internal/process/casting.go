package process

import (
	"encoding/json"
)

func (p *Process) multicast(msg *Msg) {
	// fmt.Fprintln(os.Stderr, "start multicasting...", Reset)

	for _, peer := range p.peers {
		p.unicast(msg, peer)
	}
	// fmt.Fprintln(os.Stderr, "Finished mulitcast.", Reset)
}

func (p *Process) unicast(msg *Msg, peer *Node) {
	// fmt.Fprintf(os.Stderr, Purple+"Unicasting to %s...."+Reset+"\n", peer.Id)
	if peer.Id == p.self.Id {
		p.recvd <- msg
		// fmt.Fprintf(os.Stderr, Purple+"Unicasted to SELF"+Reset+"\n")
		return
	}

	jsonData, _ := json.Marshal(msg)
	toSend := append(jsonData, '\n')
	// bytes, err := peer.Conn.Write(toSend)
	_, err := peer.Conn.Write(toSend)
	if err != nil {
		// fmt.Fprintln(os.Stderr, "ERROR:")
		// fmt.Fprintln(os.Stderr, err.Error())
	} else {
		// fmt.Fprintf(os.Stderr, Purple+"Unicasted %d bytes to %s"+Reset+"\n", bytes, peer.Id)
	}
}
