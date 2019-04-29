package pointrun

import (
	"fmt"
	"strconv"
	"time"
	vmsg "vk/messages"
	vomni "vk/omnibus"
	vrunrelayinterval "vk/run/relayinterval"
	vutils "vk/utils"
)

var Points map[string]PointRunners

func init() {
	Points = make(map[string]PointRunners)
}

func Runners() {
	relayIntervalRunners()
}

func relayIntervalRunners() {
	for k, v := range vrunrelayinterval.RunningPoints {
		if _, has := Points[k]; !has {
			Points[k] = make(map[int]Runner)
		}

		Points[k][v.Type] = v
	}
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
	var flds []string
	if flds, err = vmsg.MessageFields(msg); nil != err {
		chErr <- vutils.ErrFuncLine(err)
		return
	}

	fmt.Printf("SITKOVETSKY %+q\n", flds)

	msgNbr, err := strconv.Atoi(flds[vomni.MsgIndexPrefixNbr])
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

	if msgCd, err = strconv.Atoi(flds[vomni.MsgIndexPrefixCd]); nil != err {
		vutils.LogErr(fmt.Errorf("The Msg Code error of Msg %v", flds))
		chErr <- vutils.ErrFuncLine(err)
	}

	locErr := make(chan error)
	locDone := make(chan bool)
	locDelete := make(chan bool)

	switch msgCd {
	case vomni.MsgCdInputHelloFromPoint:
		fmt.Println("RUNNING HELLO!")
		go handleHelloFromPoint(flds, locDone, locErr)
	default:
		fmt.Println("Eduards")
	}

	select {
	case <-locDone:

	case err = <-locErr:
		chErr <- err
		return
	case <-locDelete:
		chDelete <- true
	}

	/*

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
}

func handleHelloFromPoint(flds []string, chDone chan bool, chErr chan error) {
	point := flds[vomni.MsgIndexPrefixSender]

	item, ok := Points[point]

	if ok {
		for k, v := range item {
			fmt.Println("=== KOVuktorska ", k)

			locGoOn := make(chan bool)
			locErr := make(chan error)

			v.Starter(locGoOn, locErr)

		}

		/*
			newItem := false
			if item.GetUDPAddr().IP == nil {
				newItem = true
			}

			portStr := flds[indexHelloFromPointPort]
			ipStr := flds[indexHelloFromPointIP]
			port, err := strconv.Atoi(portStr)
			if ip := net.ParseIP(ipStr); nil == ip {
				//		if nil != err {
				err = fmt.Errorf("Wrong port '%s' in the message --- %s", portStr, err.Error())
				chErr <- vutils.ErrFuncLine(err)
				return
			}

			if newItem {
				item.SetUDPAddr(ipStr, port)

				//			chMsg := make(chan string)

				fmt.Println("@@@@@@@@@@@@@@@@@@@ VILODJA MOSSABIT @@@@@@@@@@@@@@@@@@@@@@@@@@@")

				chGoOn := make(chan bool)

				go item.LetsGo(chGoOn, vomni.RootErr)
				<-chGoOn

			} else {
				item.SetUDPAddr(ipStr, port)
			}
		*/
		chDone <- true
	} else {
		err := vutils.ErrFuncLine(fmt.Errorf("Received message from the unknown point '%s'", point))
		vutils.LogErr(err)
	}
}
