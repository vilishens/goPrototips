package pointrun

import (
	"fmt"
	"net"
	"strconv"
	"time"
	vmsg "vk/messages"
	vomni "vk/omnibus"
	vrunrelayinterval "vk/run/relayinterval"
	vutils "vk/utils"
)

var Points map[string]PointRun
var ScanDone bool

func init() {
	Points = make(map[string]PointRun)
	ScanDone = false
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

func Run(chGoOn chan bool, chDone chan int, chErr chan error) {
	chGoOn <- true

	waitStart()

	fmt.Println("Alex Sitkovetsky")

	for {
		time.Sleep(vomni.DelayStepExec)
	}
}

func waitStart() {

	ScanDone = false

	for !ScanDone {
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

	switch msgCd {
	case vomni.MsgCdInputHelloFromPoint:
		fmt.Println("..............................................................")
		fmt.Printf("........................ RUNNING HELLO! %s\n", flds[vomni.MsgIndexPrefixSender])
		fmt.Println("..............................................................")

		//go handleHelloFromPoint(flds, locDone, locErr)

		signIn(flds)
		chErr <- nil

	case vomni.MsgCdOutputHelloFromStation:
		// this is the hello message from another station
		// just ignore it
		chErr <- nil
	default:
		fmt.Println("Eduards")
	}

	select {
	case <-locDone:
		// the done code received
	case err = <-locErr:
		// the error received
		vomni.RootErr <- err
		return
	}

	chErr <- err
}

func signIn(flds []string) {

	fmt.Println(".................................................>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	point := flds[vomni.MsgIndexPrefixSender]

	addr, ok := getUDPAddr(flds, vomni.MsgPrefixLen+vomni.MsgIndexHelloFromPointIP, vomni.MsgPrefixLen+vomni.MsgIndexHelloFromPointPort)
	if ok {

		newP := PointRun{}

		if _, has := Points[point]; !has {

			newP.Point.Point = point
			newP.Point.UDPAddr = addr

			Points[point] = newP
		}
	}

	//	ip := flds[vomni.MsgPrefixLen+vomni.MsgIndexHelloFromPointIP]
	//	port := flds[vomni.MsgPrefixLen+vomni.MsgIndexHelloFromPointPort]

	fmt.Printf("PEVICHKA! %+v\nPoint %q UDP %+v\n", flds, Points[point].Point.Point, Points[point].Point.UDPAddr)

	intNbr, _ := strconv.Atoi(flds[vomni.MsgIndexPrefixNbr])

	vmsg.MessageMinusByNbr(intNbr)
}

func FindDisconnectedPoint(addr net.UDPAddr) (point string) {
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
