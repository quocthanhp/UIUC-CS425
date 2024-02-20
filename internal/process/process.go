package process

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Process struct {
	self      *Node
	peers     map[string]*Node
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
			peer := new(Node)
			peer.Id = parts[0]
			peer.Hostname = parts[1]
			peer.Port = parts[2]
			p.peers[peer.Id] = peer
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

func (p *Process) Init() {
	p.peers = make(map[string]*Node)
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
