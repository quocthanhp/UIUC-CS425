package process

import (
	"fmt"
	"strconv"
	"strings"
)

type MsgType int

const (
	Normal MsgType = iota
	PrpPriority
	AgrPriority
)

var stringToMsgType = map[string]MsgType{
	"0": Normal,
	"1": PrpPriority,
	"2": AgrPriority,
}

type Msg struct {
	From     string
	Id       string
	Tx       Tx
	MT       MsgType
	Priority int
}

func parseMessageType(str string) (MsgType, error) {
	if mt, exists := stringToMsgType[str]; exists {
		return mt, nil
	}
	return -1, fmt.Errorf("invalid message type: %s", str)
}

func parseRawMessage(str string) (*Msg, error) {
	parts := strings.Split(strings.TrimSpace(str), " ")
	if len(parts) < 4 {
		return nil, fmt.Errorf("message format error")
	}

	mt, err := parseMessageType(parts[1])
	if err != nil {
		return nil, err
	}

	msg := &Msg{
		From: parts[0],
		MT:   mt,
		Id:   parts[2],
	}

	// Handling different message types
	switch msg.MT {
	case PrpPriority, AgrPriority:
		msg.Priority, err = strconv.Atoi(parts[3])
		if err != nil {
			return nil, fmt.Errorf("invalid priority value: %v", err)
		}
	case Normal:
		msg.Tx.TT, err = parseTxType(parts[3])
		if err != nil {
			return nil, fmt.Errorf("invalid tx type value")
		}
		if (msg.Tx.TT == Deposit && len(parts) != 6) || (msg.Tx.TT == Transfer && len(parts) != 7) {
			return nil, fmt.Errorf("not enough fields in the message")
		}
		msg.Tx.Amount, err = strconv.Atoi(parts[4])
		if err != nil {
			return nil, fmt.Errorf("invalid amount value")
		}
		msg.Tx.To = parts[5]
		if msg.Tx.TT == Transfer {
			msg.Tx.From = parts[6]
		}
	}

	return msg, nil
}