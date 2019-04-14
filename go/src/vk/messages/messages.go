package messages

import (
	"fmt"
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

func (d SendMsgArray) MinusNbr(nbr int) {
	ind := -1
	for key, val := range d {
		if val.MessageNbr == nbr {
			ind = key
			break
		}
	}

	if 0 > ind {
		fmt.Printf("Received MSG #%d without record\n", nbr)
		return
	}

	chDone := make(chan bool)
	go d.MinusIndex(ind, chDone)
	<-chDone
}
