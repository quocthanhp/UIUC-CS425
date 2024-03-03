package process

import (
	"fmt"
	"os"
)

func (p *Process) deliver(msg *Msg) {
	fmt.Fprintln(os.Stderr, msg.toString())
	if msg.Tx.TT == Deposit {
		p.bank.Deposit(msg.Tx.To, msg.Tx.Amount)
	} else if msg.Tx.TT == Transfer {
		p.bank.Transfer(msg.Tx.From, msg.Tx.To, msg.Tx.Amount)
	} else {
		fmt.Println("Invalid Transaction Type")
	}
}
