package process

import (
	"bufio"
	"fmt"
	"log"
	"mp1_node/internal/node"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Process struct {
	self      *node.Node
	peers     []*node.Node
	groupSize int
	ln        net.Listener
}

func (p *Process) ReadPeersInfo(self_id string, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening configuration file:", err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		size, err := strconv.Atoi(strings.TrimSpace(line))
		if err != nil {
			fmt.Println("Err convert group size to int:", err)
			os.Exit(1)
		}
		p.groupSize = size
	}

	for i := 0; i < p.groupSize; i++ {
		if scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
			if err != nil {
				fmt.Println("Err convert group size to int:", err)
				os.Exit(1)
			}
			line = strings.TrimSpace(line)
			parts := strings.Split(line, " ")
			if len(parts) != 3 {
				fmt.Println("Wrong formatting in config file:", line)
				continue
			}
			peer := new(node.Node)
			peer.Id = parts[0]
			peer.Hostname = parts[1]
			peer.Port = parts[2]
			p.peers = append(p.peers, peer)
			if peer.Id == self_id {
				p.self = peer
			}
		} else {
			fmt.Println("Wrong formatting in config file.")
			os.Exit(1)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading configuration file:", err)
	}
}

func (p *Process) handleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)

	buf, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		return
	}
	fmt.Printf("Received: %s", buf)
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

func connectToSinglePeer(node *node.Node, wg *sync.WaitGroup) {
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

func (p *Process) Start() {
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		p.connectToPeers()
	}()

	go func() {
		defer wg.Done()
		p.startListen()
	}()

	wg.Wait()
}

func (p *Process) Clean() {
	p.ln.Close()
	for _, peer := range p.peers {
		peer.Conn.Close()
	}
}
