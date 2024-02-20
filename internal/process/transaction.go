package process

import (
	"fmt"
)

type TxType int

const (
	Transfer TxType = iota
	Deposit
)

var stringToTxType = map[string]TxType{
	"0": Transfer,
	"1": Deposit,
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
}

func parseTxType(str string) (TxType, error) {
	if tt, exists := stringToTxType[str]; exists {
		return tt, nil
	}
	return -1, fmt.Errorf("invalid tx type: %s", str)
}
