package process

import (
	"fmt"
	"time"
)

func GetTimeDiffInMilliSeconds(msg *Msg) int64 {
	parsedTime, err := time.Parse(time.RFC3339, msg.Tx.Timestamp)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return -1
	}
	now := time.Now()
	diff := now.Sub(parsedTime)
	diffInMs := diff.Milliseconds()
	return diffInMs
}
