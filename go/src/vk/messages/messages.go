package messages

import (
	"time"
	vomni "vk/omnibus"
)

var MessageList2Send SendMsgArray

func init() {
	MessageList2Send = SendMsgArray{}

	vomni.MessageNumber = 0
}

func Run(chGoOn chan bool, chDone chan int, chErr chan error) {

	chGoOn <- true
	for {
		time.Sleep(vomni.DelayStepExec)
	}
}
