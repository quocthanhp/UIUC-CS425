package process

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

func (p *Process) matchConnToPeer(conn net.Conn) {
	reader := bufio.NewReader(conn)
	buf, err := reader.ReadString('\n')
	buf = strings.TrimSpace(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}
	if _, ok := p.peers[buf]; !ok {
		fmt.Println("Cannot identify the connected peer!")
		os.Exit(1)
	}
	p.peers[buf].Conn = conn
	fmt.Printf("[SELF] Established connection with peer <%s>\n", buf)
}

func (p *Process) createConnection(conn net.Conn) {
	p.matchConnToPeer(conn)
}

func (p *Process) startListen() {
	var wg sync.WaitGroup
	ln, err := net.Listen("tcp", ":"+p.self.Port)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	p.ln = ln
	fmt.Printf("Listening on port %s....\n", p.self.Port)

	// SHOULD IT BE groupsize - 1 ????
	for i := 0; i < p.groupSize-1; i++ {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %s\n", err)
			continue
		}
	
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.createConnection(conn)
		}()
	}
	wg.Wait()
}

func (p *Process) connectToSinglePeer(node *Node, wg *sync.WaitGroup) {
	defer wg.Done()
	connected := false
	for !connected {
		conn, err := net.Dial("tcp", node.GetIPAddr())
		if err != nil {
			fmt.Printf("Cannot connect to peer %s, retrying after 1 second...\n", node.Id)
			time.Sleep(time.Second)
			continue
		}
		node.Conn = conn
		fmt.Fprintf(conn, "%s\n", p.self.Id)
		connected = true
	}
}

func (p *Process) connectToPeers() {
	var wg sync.WaitGroup

	for _, peer := range p.peers {
		if peer != p.self {
			wg.Add(1)
			go p.connectToSinglePeer(peer, &wg)
		}
	}
	wg.Wait()
	fmt.Println("Connected to All Peers!")
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
