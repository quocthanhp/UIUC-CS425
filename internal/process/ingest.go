package process

import "fmt"

// write code to ingest all the messages
func (p *Process) Ingest() {
	for msg := range p.recvd {
		fmt.Println(msg)
	}
}
