package process

import (
	"bufio"
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

func (p *Process) handleConnection(conn net.Conn) {
	p.matchConnToPeer(conn)

	reader := bufio.NewReader(conn)
	for {
		buf, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Println("Client closed the connection")
			} else {
				fmt.Println("Error reading:", err.Error())
			}
			break
		}
		msg, err := parseRawMessage(buf)
		if err != nil {
			fmt.Println("[ERROR] Invalid Message")
			continue
		}
		p.recvd <- msg
	}
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

	for i := 0; i < p.groupSize; i++ {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %s\n", err)
			continue
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			p.handleConnection(conn)
		}()
	}
	wg.Wait()
}

func connectToSinglePeer(node *Node, wg *sync.WaitGroup) {
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
		fmt.Fprintf(conn, "%s\n", node.Id)
		connected = true
	}
}

func (p *Process) connectToPeers() {
	var wg sync.WaitGroup

	for _, peer := range p.peers {
		if peer != p.self {
			wg.Add(1)
			go connectToSinglePeer(peer, &wg)
		}
	}
	wg.Wait()
	fmt.Println("Connected to All Peers!")
}
