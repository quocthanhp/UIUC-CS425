package process

import (
	"strconv"
)

func (p *Process) deliver(msg *Msg) {
	// fmt.Fprintln(os.Stderr, msg.toString())
	diffInMs := GetTimeDiffInMilliSeconds(msg)

	// Write some content to the file
	_, err := p.log_writer.WriteString(msg.Id + " " + strconv.Itoa(int(diffInMs)) + "\n")
	if err != nil {
		// fmt.Fprintln(os.Stderr, "Error writing timestamp to file:", err)
		return
	}

	err = p.log_writer.Flush()
	if err != nil {
		// fmt.Fprintln(os.Stderr, "Error flushing timestamp buffer to file:", err)
		return
	}
	if msg.Tx.TT == Deposit {
		p.bank.Deposit(msg.Tx.To, msg.Tx.Amount)
	} else if msg.Tx.TT == Transfer {
		p.bank.Transfer(msg.Tx.From, msg.Tx.To, msg.Tx.Amount)
	} else {
		// fmt.Fprintln(os.Stderr, "Invalid Transaction Type")
	}
}
