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

// var stringToMsgType = map[string]MsgType{
// 	"0": Normal,
// 	"1": PrpPriority,
// 	"2": AgrPriority,
// }

type Msg struct {
	From     string
	Id       string
	Tx       Tx
//	MT       MsgType
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

	switch parts[0] {
	case "DEPOSIT":
		if (len(parts) != 3) {
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

		msg = &Msg{
			From: node,
			Id:   node + "-" + strconv.Itoa(messageNum),
			Tx: Tx{To: to, Amount: amnt, TT: tt},
		}
	case "TRANSFER":
		if (len(parts) != 5) {
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
		
		msg = &Msg{
			From: node,
			Id:   node + "-" + strconv.Itoa(messageNum),
			Tx: Tx{From: from, To: to, Amount: amnt, TT: tt},
		}
	}

	messageNum++;
	return msg, nil
}
