package process

import (
	"bufio"
	"encoding/json"
	"fmt"
	"mp1_node/internal/util"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Process struct {
	self      *Node
	peers     map[string]*Node
	groupSize int
	ln        net.Listener
	recvd     chan *Msg
	verified  chan *Msg
	send      chan *Msg
	msgs      map[string]*PdMsg
	bank      *Bank
	hMsg      map[[32]byte]bool
}

// type HashMsg struct {
// 	hashMsg map[[32]byte]bool
// 	hLock   sync.Mutex
// }

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
	p.recvd = make(chan *Msg, 200)
	p.verified = make(chan *Msg, 200)
	p.send = make(chan *Msg, 200)
	p.msgs = make(map[string]*PdMsg)
	p.hMsg = make(map[[32]byte]bool, 200)
	p.bank = NewBank()
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
	go p.handlePeerConnections()
	time.Sleep(2 * time.Second)
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

	wg.Add(3)
	//TODO: handle timeout proposal
	clearStdin()

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
	for {
		select {
		case msg := <-p.recvd:
			fmt.Printf(Yellow+"[GOT MSG] %s\n"+Reset, msg.toString())

			// Calculate hash to detect duplication
			data, _ := json.Marshal(*msg)
			hashVal := util.GetHash(data)
			if _, ok := p.hMsg[hashVal]; ok {
				//fmt.Println("duplicated message!")
				continue
			}
			p.hMsg[hashVal] = true
			//fmt.Printf("[HASH] of %s is %s\n", msg.Id, hashVal)

			if msg.MT == Normal {
				p.msgs[msg.Id] = &PdMsg{msg, 0}

				if msg.From != p.self.Id {
					//fmt.Println("Reliably multicast NORMAL")
					p.multicast(msg)
				}
			}

			p.verified <- msg

			// key := string(rawbytes)
			// fmt.Printf(Yellow+"[GOT MSG] %s\n"+Reset, string(rawbytes))
			// if msg.MT == PrpPriority {
			// 	fmt.Println("First branch")
			// 	p.verified <- msg
			// } else if _, ok := receivedMsg[key]; !ok {
			// 	fmt.Println("Second branch")
			// 	receivedMsg[key] = struct{}{}
			// 	fmt.Println("Second branch")
			// 	if msg.MT == Normal {
			// 		fmt.Println("Second branch")
			// 		p.msgs[msg.Id] = &PdMsg{msg, 0}
			// 	}
			// 	fmt.Println("Second branch")
			// 	// fmt.Printf(Green+"[FROM %s] %s\n"+Reset, msg.From, string(rawbytes))
			// 	p.verified <- msg
			// 	if msg.From != p.self.Id {
			// 		fmt.Println("Second branch")
			// 		go p.multicast(msg)
			// 	}
			// }
		case msg := <-p.send:
			msg.From = p.self.Id
			go p.multicast(msg)
		}
	}
}

// func (p *Process) AppendHash(hash [32]byte) {
// 	p.hMsg.hLock.Lock()
// 	defer p.hMsg.hLock.Unlock()
	
// 	p.hMsg.hashMsg[hash] = true
// }
