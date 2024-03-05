package process

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"strings"
)

func (p *Process) ReadInput() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		msg, err := ToNetworkMsg(p.self.Id, line)
		if err != nil {
			// fmt.Fprintln(os.Stderr, err)
			return
		}

		// fmt.Fprintln(os.Stderr, "MSG FROM STDIN:", line)
		p.send <- msg
	}

	if err := scanner.Err(); err != nil {
		// fmt.Fprintln(os.Stderr, "Error reading from stdin:", err)
	}
}

func (p *Process) handlePeerConnections() {
	for _, peer := range p.peers {
		if peer != p.self {
			go p.handleSingleConnection(peer)
		}
	}
}

func (p *Process) handleSingleConnection(peer *Node) {
	reader := bufio.NewReader(peer.Conn)
	for {
		buf, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				// fmt.Fprintln(os.Stderr, "Client closed the connection")
				p.handleFailure(peer)
				break
			}

			// fmt.Fprintln(os.Stderr, err)
			break
		}

		var msg Msg
		err = json.Unmarshal(buf, &msg)
		if err != nil {
			// fmt.Fprintln(os.Stderr, "Error reading:", err.Error())
			continue
		}
		p.recvd <- &msg
	}
}
