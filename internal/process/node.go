package process

import "net"

type Node struct {
	Id       string
	Hostname string
	Port     string
	Conn     net.Conn
}

func (node *Node) GetIPAddr() string {
	return node.Hostname + ":" + node.Port
}
