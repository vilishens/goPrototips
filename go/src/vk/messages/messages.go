package messages

import (
	"sync"
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

func (d SendMsgArray) MinusIndex(ind int, chDone chan bool) {

	if ind < len(d) {
		lock := new(sync.Mutex)
		lock.Lock()
		defer lock.Unlock()

		MessageList2Send = append(MessageList2Send[:ind], MessageList2Send[ind+1:]...)
	}

	chDone <- true
}
