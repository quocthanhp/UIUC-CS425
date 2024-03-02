package process

import "fmt"

func (p *Process) deliver(msg *Msg) {
	if msg.Tx.TT == Deposit {
		p.bank.Deposit(msg.Tx.To, msg.Tx.Amount)
	} else if msg.Tx.TT == Transfer {
		p.bank.Transfer(msg.Tx.From, msg.Tx.To, msg.Tx.Amount)
	} else {
		fmt.Println("Invalid Transaction Type")
	}
}
