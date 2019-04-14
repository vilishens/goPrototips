package pointrun

import (
	"fmt"
	"strconv"
	"strings"
	"time"
	vmsg "vk/messages"
	vomni "vk/omnibus"
	vutils "vk/utils"
)

func init() {
}

func Run(chGoOn chan bool, chDone chan int, chErr chan error) {
	chGoOn <- true
	for {
		time.Sleep(vomni.DelayStepExec)
	}
}

func MessageReceived(msg string, chErr chan error) {

	fmt.Println("vk-xxx @@@@@@ SITKOVETSKY @@@@@ MSG", msg)

	var err error
	flds := strings.Split(msg, vomni.UDPMessageSeparator)

	msgNbr, err := strconv.Atoi(flds[vomni.MsgIndexNbr])
	if nil != err {
		vutils.LogErr(fmt.Errorf("The Msg Number error of Msg %q", msg))
		chErr <- vutils.ErrFuncLine(err)
		return
	}

	locErr := make(chan error)
	locDelete := make(chan bool)

	go messageReceived(flds, locDelete, locErr)
	select {
	case <-locDelete:
		vmsg.MessageList2Send.MinusNbr(msgNbr)
	case err = <-locErr:
		break
	}

	chErr <- err
}

func messageReceived(flds []string, chDelete chan bool, chErr chan error) {

	var err error
	msgCd := -1

	if msgCd, err = strconv.Atoi(flds[vomni.MsgIndexCd]); nil != err {
		vutils.LogErr(fmt.Errorf("The Msg Code error of Msg %v", flds))
		chErr <- vutils.ErrFuncLine(err)
	}

	_ = msgCd
	/*
		locErr := make(chan error)
		locDone := make(chan bool)
		locDelete := make(chan bool)

		if msgCd == vomni.MsgCdInputHelloFromPoint {

			fmt.Println("RUNNING HELLO")

			go handleHelloFromPoint(flds, locDone, locErr)

			fmt.Println("RUNNING TI KUDA???")

			select {
			case <-locDone:
				//
			case err = <-locErr:
				chErr <- err
				return
			}
		}

		point := flds[vomni.MsgIndexMsgSender]
		item, ok := Points[point]
		if !ok {
			chErr <- vutils.ErrFuncLine(fmt.Errorf("\nThe message received from the unknown point %q", point))
			return
		}

		go item.Response(flds, locDelete, locErr)

		select {
		case <-locDelete:
			chDelete <- true
			return
		case err = <-locErr:
			chErr <- vutils.ErrFuncLine(err)
			return
		}

		//	default:
		//		err = fmt.Errorf("\nNo logic for the received message code 0x%08X\n", msgCd)
		//		break
		//	}
	*/
	chErr <- err
}
