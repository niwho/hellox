package common

import (
	"testing"
	"time"
)

func Test_Dingding(t *testing.T) {
	for i := 0; i < 1; i += 1 {
		//		SendDingDing(fmt.Sprintf("log:%d", i), fmt.Sprintf("message_seq:%d", i))
		time.Sleep(time.Second)
	}
}
