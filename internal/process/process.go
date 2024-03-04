package process

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Process struct {
	self       *Node
	peers      map[string]*Node
	groupSize  int
	ln         net.Listener
	recvd      chan *Msg
	verified   chan *Msg
	send       chan *Msg
	msgs       map[string]*PdMsg
	bank       *Bank
	que        MsgQ
	log_writer *bufio.Writer
}

// var receivedMsg = make(map[string]struct{})

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
	p.recvd = make(chan *Msg, 200)
	p.verified = make(chan *Msg, 200)
	p.send = make(chan *Msg, 200)
	p.msgs = make(map[string]*PdMsg)
	p.bank = NewBank()
	p.que = MsgQ{}
}

func (p *Process) Start() {
	file, err := os.Create("timestamp_log/timestamplog-" + p.self.Id)
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	} else {
		fmt.Println(Cyan, "LOG FILE CREATED!", Reset)
		p.log_writer = bufio.NewWriter(file)
	}
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
	go p.handlePeerConnections()
	time.Sleep(4 * time.Second)
	clearStdin()
	fmt.Println(Cyan, "READY!", Reset)
}

func clearStdin() {
	// Create a new reader for stdin
	reader := bufio.NewReader(os.Stdin)
	// Read and discard bytes until a newline (or EOF) is encountered
	for {
		_, err := reader.ReadString('\n')
		if err != nil {
			break // EOF or an actual error occurred
		}
		if reader.Buffered() == 0 {
			break // No more data to read
		}
	}
}

func (p *Process) Run() {
	var wg sync.WaitGroup

	wg.Add(2)
	//TODO: handle timeout proposal
	// clearStdin()

	go func() {
		defer wg.Done()
		p.ReadInput()
	}()

	go func() {
		defer wg.Done()
		p.MonitorChannel()
	}()

	go func() {
		defer wg.Done()
		p.Ordering()
	}()

	wg.Wait()
}

func (p *Process) Clean() {
	p.ln.Close()
	for _, peer := range p.peers {
		peer.Conn.Close()
	}
}

func (p *Process) MonitorChannel() {
	set := make(Set)
	for {
		select {
		case msg := <-p.recvd:
			if msg.MT == PrpPriority {
				fmt.Printf(Yellow+"[GOT MSG] %s\n"+Reset, msg.toString())
				p.verified <- msg
			} else if !set.Contains(msg.toString()) {
				fmt.Printf(Yellow+"[GOT MSG] %s\n"+Reset, msg.toString())
				set.Add(msg.toString())
				if msg.MT == Normal {
					p.msgs[msg.Id] = &PdMsg{msg, 0}
				}
				p.verified <- msg
				p.multicast(msg)
			}
		case msg := <-p.send:
			msg.From = p.self.Id
			go p.multicast(msg)
		}
	}
}

func (p *Process) handleFailure(peer *Node) {
	p.que.removeDeprecatedMsg(peer)
	p.groupSize--
	delete(p.peers, peer.Id)
}
