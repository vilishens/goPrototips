package pointrun

import (
	"fmt"
	"net"
	"strconv"
	"time"
	vmsg "vk/messages"
	vnetscan "vk/net/netscan"
	vomni "vk/omnibus"
	vrunrelayinterval "vk/run/relayinterval"
	vutils "vk/utils"
)

var Points map[string]PointRun
var listSigned map[string]net.UDPAddr
var startSequence []int

func init() {
	Points = make(map[string]PointRun)
	listSigned = make(map[string]net.UDPAddr)
	startSequence = []int{vomni.CfgTypeRelayInterval}
}

func Runners() {
	relayIntervalRunners()
}

func relayIntervalRunners() {
	for k, v := range vrunrelayinterval.RunningPoints {
		tPoint := PointData{}
		tRun := PointRunners{}

		tPoint.Point = v.Point
		tPoint.UDPAddr = net.UDPAddr{}
		tPoint.Type |= v.Type //   vomni.CfgTypeRelayInterval
		tPoint.State = v.State

		tRun[v.Type] = v

		Points[k] = PointRun{Point: tPoint, Run: tRun}
	}
}

func RunStart(chGoOn chan bool, chDone chan int, chErr chan error) {

	locGoOn := make(chan bool)
	locDone := make(chan int)
	locErr := make(chan error)

	listSigned = map[string]net.UDPAddr{}

	go scanStationNet(locGoOn, locDone, locErr)

	stop := false
	for {
		select {
		case err := <-locErr:
			chErr <- err
			return
		case done := <-chDone:
			chDone <- done
			return
		case <-locGoOn:
			fmt.Println("### Kurtenkov ###")
			stop = true

		}

		if stop {
			break
		}

	}

	fmt.Println("Alex Sitkovetsky ")
	chGoOn <- true

	fmt.Println("TAGAD starts", len(listSigned))
	fmt.Printf("TAGAD oooooo %+v\n", listSigned)

	for {
		time.Sleep(vomni.DelayStepExec)
	}
}

func scanStationNet(chGoOn chan bool, chDone chan int, chErr chan error) {

	locGoOn := make(chan bool)
	locDone := make(chan int)
	locErr := make(chan error)

	go scanNet(locGoOn, locDone, locErr)

	stop := false
	for {
		select {
		case err := <-locErr:
			chErr <- err
			return
		case <-chDone:
			// the done code received
			stop = true
		case <-locGoOn:
			stop = true
		}
		if stop {
			break
		}
	}

	go startSigned(locGoOn, locErr)

	stop = false
	for {
		select {
		case err := <-locErr:
			chErr <- err
			return
		case <-locGoOn:
			chGoOn <- true
			stop = true
		}
		if stop {
			break
		}
	}
}

func startSigned(chGoOn chan bool, chErr chan error) {

	for _, cfgType := range startSequence {

		for point, addr := range listSigned {
			if pData, ok := Points[point]; !ok {
				err := fmt.Errorf("Unknown point %q (%v) sent SignIn message", point, addr)
				vutils.LogErr(err)
				chErr <- vutils.ErrFuncLine(err)
				return
			} else {

				if 0 != (pData.Point.Type & cfgType) {
					fmt.Println("SEIT JĀsĀk run ", pData.Point.Point)
				}
			}
		}
	}
}

func scanNet(chGoOn chan bool, chDone chan int, chErr chan error) {

	locGoOn := make(chan bool)
	locDone := make(chan int)
	locErr := make(chan error)

	go vnetscan.ScanOctet(locGoOn, locDone, locErr)

	for {
		select {
		case err := <-locErr:
			chErr <- err
			return
		case done := <-chDone:
			chDone <- done
			return
		case <-locGoOn:
			chGoOn <- true
			return
		}
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
		err = fmt.Errorf("The Msg Code error of Msg %v", flds)
		vutils.LogErr(err)
		chErr <- vutils.ErrFuncLine(err)
	}

	locErr := make(chan error)
	locDone := make(chan bool)
	locDelete := make(chan bool)

	switch msgCd {
	case vomni.MsgCdInputHelloFromPoint:
		fmt.Println("..............................................................")
		fmt.Printf("........................ RUNNING HELLO! %s\n", flds[vomni.MsgIndexPrefixSender])
		fmt.Println("..............................................................")

		//go handleHelloFromPoint(flds, locDone, locErr)

		go addSignIn(flds, locDelete, locErr)

	case vomni.MsgCdOutputHelloFromStation:
		// this is the hello message from another station
		// just ignore it
		chErr <- nil
		return
	default:
		chErr <- vutils.ErrFuncLine(fmt.Errorf("RECEIVED->RECEIVED->RECEIVED unknowm CMD %d", msgCd))
		fmt.Println("Eduards")
		return
	}

	select {
	case <-locDone:
		// the done code received
	case <-locDelete:
		chDelete <- true
	case err = <-locErr:
		// the error received
		vomni.RootErr <- err
		return
	}

	chErr <- err
}

func addSignIn(flds []string, chDelete chan bool, chErr chan error) {

	fmt.Println(".................................................>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	point := flds[vomni.MsgIndexPrefixSender]
	addr, ok := getUDPAddr(flds, vomni.MsgPrefixLen+vomni.MsgIndexHelloFromPointIP, vomni.MsgPrefixLen+vomni.MsgIndexHelloFromPointPort)

	if ok {
		listSigned[point] = addr
	}

	chDelete <- true

	fmt.Printf("PEVICHKA! %+v\nPoint %q UDP %+v\n", flds, Points[point].Point.Point, Points[point].Point.UDPAddr)
}

func FindDisconnectedPoint(addr net.UDPAddr) (point string) {

	fmt.Println("SITAS VEL JAIZVEIDO!!!! ", addr)

	return
}

func handleHelloFromPoint(flds []string, chDone chan bool, chErr chan error) {
	point := flds[vomni.MsgIndexPrefixSender]

	fmt.Println("#### SLUCHAJ #####")

	item, ok := Points[point]

	if ok {
		if addr, ok := getUDPAddr(flds, vomni.MsgPrefixLen+vomni.MsgIndexHelloFromPointIP, vomni.MsgPrefixLen+vomni.MsgIndexHelloFromPointPort); ok {
			for k, v := range item.Run {

				_ = k

				locGoOn := make(chan bool)
				locDone := make(chan int)
				locErr := make(chan error)

				go v.LetsGo(addr, flds, locGoOn, locDone, locErr)

				select {
				case <-locGoOn:
					// all done return flag to go on
					chDone <- true
				//case rc := <-locDone:
				case <-locDone:
					// the done code received
				case err := <-locErr:
					chErr <- vutils.ErrFuncLine(fmt.Errorf("Couldn't handle Starter of point %q - %s",
						flds[vomni.MsgIndexPrefixSender], err.Error()))
					return
				}
			}
		}

	} else {
		err := vutils.ErrFuncLine(fmt.Errorf("Received message from the unknown point %q", point))
		vutils.LogErr(err)
	}
}

func getUDPAddr(flds []string, ipInd int, portInd int) (addr net.UDPAddr, ok bool) {

	intPort, err := strconv.Atoi(flds[portInd])
	if nil != err {
		err = vutils.ErrFuncLine(fmt.Errorf("A message received (%v) with the wrong Port format %q - %s",
			flds,
			flds[portInd],
			err.Error()))
		vutils.LogErr(err)
	}

	netIP := net.ParseIP(flds[ipInd])
	if nil == netIP {
		err = vutils.ErrFuncLine(fmt.Errorf("A message received (%v) with the invalid IP %q",
			flds,
			flds[ipInd]))
		vutils.LogErr(err)
	}

	if nil != err {
		return
	}

	addr = net.UDPAddr{IP: netIP, Port: intPort}

	return addr, true
}
