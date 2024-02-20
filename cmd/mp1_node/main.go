package main

import (
	"fmt"
	"mp1_node/internal/process"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("./mp1_node <identifier> <configuration file>")
		os.Exit(1)
	}

	var proc process.Process
	proc.Init()
	// defer proc.Clean()
	proc.ReadPeersInfo(os.Args[1], os.Args[2])
	proc.Start()

	fmt.Println("Hello, Go!")
}
