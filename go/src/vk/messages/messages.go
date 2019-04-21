package messages

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
	vomni "vk/omnibus"
	vparams "vk/params"
)

var MessageList2Send SendMsgArray

func init() {
	MessageList2Send = SendMsgArray{}

	vomni.MessageNumber = 0
}

func Message2SendPlus(addr net.UDPAddr, msgCd int, data []string) {

	msg := message2SendNew(msgCd, data)
	message2SendAdd(addr, msg)
}

func message2SendAdd(addr net.UDPAddr, msg string) {

	d := SendMsg{}

	d.UDPAddr = addr
	d.MessageNbr = vomni.MessageNumber
	d.Msg = msg

	MessageList2Send = append(MessageList2Send, d)
}

func message2SendNew(msgCd int, data []string) (msg string) {

	vomni.MessageNumberNext()

	msg = ""
	msg += vparams.Params.StationName + vomni.UDPMessageSeparator
	msg += strconv.Itoa(msgCd) + vomni.UDPMessageSeparator
	msg += strconv.Itoa(vomni.MessageNumber)

	for _, v := range data {
		msg += vomni.UDPMessageSeparator + v
	}

	return msg
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

func TryHello(dst net.UDPAddr, chDone chan bool) {

	msgData := msgStationHello()

	Message2SendPlus(dst, vomni.MsgCdOutputHelloFromStation, msgData)

	fmt.Println("================== try ================================> ", dst, " #", vomni.MessageNumber)

	chDone <- true

}

func msgStationHello() (d []string) {

	_, tzSecs := time.Now().Zone()

	d = make([]string, vomni.MsgHelloFromStationLen)

	d[vomni.MsgIndexHelloFromStationTime] = strconv.Itoa(int(time.Now().Unix()))
	d[vomni.MsgIndexHelloFromStationOffset] = strconv.Itoa(tzSecs)
	d[vomni.MsgIndexHelloFromStationIP] = vparams.Params.IPAddressInternal
	d[vomni.MsgIndexHelloFromStationPort] = strconv.Itoa(vparams.Params.PortUDPInternal)

	return
}
