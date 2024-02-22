package process

import (
	"fmt"
	"strconv"
)

type TxType int

const (
	Transfer TxType = iota
	Deposit
)

var stringToTxType = map[string]TxType{
	"TRANSFER": Transfer,
	"DEPOSIT": Deposit,
}

type TxStatus int

const (
	Undeliverable TxStatus = iota
	Deliverable
)

type Tx struct {
	From   string
	To     string
	Amount int
	TT     TxType
	TS     TxStatus
}

func parseTxType(str string) (TxType, error) {
	if tt, exists := stringToTxType[str]; exists {
		return tt, nil
	}
	return -1, fmt.Errorf("invalid tx type: %s", str)
}

func getTransactionString(tx Tx) string {
	if tx.TT == Deposit {
		return "DEPOSIT " + tx.To + " " + strconv.Itoa(tx.Amount)
	}

	return "TRANSFER " + tx.From +  " " + tx.To + " " + strconv.Itoa(tx.Amount)
}