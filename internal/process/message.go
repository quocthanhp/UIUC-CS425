package process

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type MsgType string

const (
	Normal      MsgType = "NML"
	PrpPriority MsgType = "PP"
	AgrPriority MsgType = "AP"
)

type Msg struct {
	From     string
	Id       string
	Tx       Tx
	MT       MsgType
	Priority int
}

var messageNum int

// func parseMessageType(str string) (MsgType, error) {
// 	if mt, exists := stringToMsgType[str]; exists {
// 		return mt, nil
// 	}
// 	return -1, fmt.Errorf("invalid message type: %s", str)
// }

func ToNetworkMsg(node string, rawMessage string) (*Msg, error) {
	parts := strings.Split(strings.TrimSpace(rawMessage), " ")
	if len(parts) < 3 {
		return nil, fmt.Errorf("message format error")
	}

	var msg *Msg
	var tx Tx

	switch parts[0] {
	case "DEPOSIT":
		if len(parts) != 3 {
			return nil, fmt.Errorf("not enough fields in the message")
		}

		tt, err := parseTxType("DEPOSIT")
		if err != nil {
			return nil, err
		}

		to, amount := parts[1], parts[2]

		amnt, err := strconv.Atoi(amount)
		if err != nil {
			return nil, err
		}

		tx = Tx{To: to, Amount: amnt, TT: tt}
	case "TRANSFER":
		if len(parts) != 5 {
			return nil, fmt.Errorf("not enough fields in the message")
		}

		tt, err := parseTxType("TRANSFER")
		if err != nil {
			return nil, err
		}

		from, to, amount := parts[1], parts[3], parts[4]

		amnt, err := strconv.Atoi(amount)
		if err != nil {
			return nil, err
		}

		tx = Tx{From: from, To: to, Amount: amnt, TT: tt}
	}
	msg = &Msg{
		Id: uuid.New().String(),
		MT: Normal,
		Tx: tx,
	}

	messageNum++
	return msg, nil
}

func (p *Process) contains(Id string) bool {
	_, ok := p.msgs[Id]
	return ok
}

func (msg *Msg) toString() string {
	bytes, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("JSON marshalling failed: %v\n", err)
	}
	return string(bytes)
}



