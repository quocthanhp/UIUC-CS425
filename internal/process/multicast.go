package process

func (p *Process) multicast() {
	for msg := range p.send {
		for _, peer := range p.peers {
			peer.Conn.Write([]byte(msg.ToNetworkMsg(p.self.Id)))
		}
	}
}
