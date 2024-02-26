package process

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func (p *Process) ReadInput() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		msg, err := ToNetworkMsg(p.self.Id, line)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("MSG FROM STDIN:", line)
		p.send <- msg
		p.recvd <- msg
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading from stdin:", err)
	}
}

func (p *Process) handlePeerConnections() {
	for _, peer := range p.peers {
		if peer != p.self {
			go p.handleSingleConnection(peer.Conn)
		}
	}
}

func (p *Process) handleSingleConnection(conn net.Conn) {
	for {
		buf, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Client closed the connection")
				break
			}

			fmt.Println(err)
			break
		}

		var msg Msg
		err = json.Unmarshal(buf, &msg)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			continue
		}
		p.recvd <- &msg
	}
}
