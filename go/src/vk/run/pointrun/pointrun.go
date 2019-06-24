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

	// prepare storage for signed in points
	listSigned = map[string]net.UDPAddr{}

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

	go startSigned(locGoOn, locDone, locErr)

	stop = false
	for {
		select {
		case err := <-locErr:
			chErr <- err
			return
		case <-locDone:
			// the done code received
			stop = true
		case <-locGoOn:
			chGoOn <- true
			stop = true
		}
		if stop {
			break
		}
	}
}

func startSigned(chGoOn chan bool, chDone chan int, chErr chan error) {

	for _, cfgType := range startSequence {
		// start all point configuration, sequence set in startSequence
		// Sequence can be important (for instance, to check the point ready state
		for point, addr := range listSigned {
			if pData, ok := Points[point]; !ok {
				err := fmt.Errorf("The point %q (%v) sent SignIn message, but there is no comfiguration of this point", point, addr)
				vutils.LogErr(err)
				chErr <- vutils.ErrFuncLine(err)
				return
			} else {

			if 0 != (pData.Point.Type & cfgType) {
					// the point has configuration of this point
					pt := Points[point]
					ptPt := pt.Point

					logStr := ""
					start := false
					if 0 != (ptPt.State & vomni.PointStateDisconnected) {
						// this point was signed in, but later disconnected
						// need to restart again
						logStr = fmt.Sprintf("Point %q signed in AGAIN", point)
					} else if 0 == (ptPt.State & vomni.PointStateSigned) {
						// the point wasn't signed in, need to start from scratch
						logStr = fmt.Sprintf("Point %q signed in", point)
						start = true
					} else {
						// the point was signed and not disconnected, to update the address is enough
						logStr = fmt.Sprintf("Point %q signed in used the new UDP address %s:%d", point, addr.IP.String(), addr.Port)
					}

					ptPt.UDPAddr = addr
					ptPt.State &^= vomni.PointStateDisconnected
					ptPt.State |= vomni.PointStateSigned

					pt.Point = ptPt
					//					Points[point] = pt

					//Points[point].Point = ptPt

					//					Points[point].Run[cfgType].Logofet(addr)

					//					Points[point].Run[cfgType].SetUDPAddr(addr)
					//############################

					Points[point].Run[cfgType].SetUDPAddr(addr)

					//############################

					fmt.Println("vk-xxx MIZANDARI ", Points[point].Run[cfgType].GetUDPAddr())

					// put messages about signed in into log
					vutils.LogInfo(logStr)
					if !start {
						// rotate files is ready, we can put the message
						Points[point].Run[cfgType].LogStr(vomni.LogFileCdInfo, logStr)
					}

					if start {
						locGoOn := make(chan bool)
						locDone := make(chan int)
						locErr := make(chan error)

						err := Points[point].Run[cfgType].StartRotate()
						if nil != err {
							chErr <- err
							return
						}

						// rotate is now ready, let's put the signed in message into log
						Points[point].Run[cfgType].LogStr(vomni.LogFileCdInfo, logStr)

						go Points[point].Run[cfgType].LetsGo(addr, locGoOn, locDone, locErr)

						select {
						case <-locGoOn:
						case cd := <-locDone:
							chDone <- cd
						case err := <-locErr:
							chErr <- err
						}

						fmt.Println("vk-xxx LOMBARDS ", Points[point].Run[cfgType].GetUDPAddr())

					}

					fmt.Println("SEIT JĀsĀk run ", pData.Point.Point)
				}
			}
		}
	}

	chGoOn <- true
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
		// just ignore it and send delete it signal
		chDelete <- true
		return
	case vomni.MsgCdInputSuccess:
		// don't do anything - just send delete it signal

		fmt.Println(flds[vomni.MsgIndexPrefixSender], "@@@@@@@@@@@@@ vk-xxx -------> SUCCESS received")

		chDelete <- true
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

	// send back the flag to delete this message
	chDelete <- true

	fmt.Printf("vk-xxx PEVICHKA! %+v\nPoint %q UDP %+v\n", flds, Points[point].Point.Point, Points[point].Point.UDPAddr)
}

func SetDisconnectedPoint(addr net.UDPAddr) (point string) {
	for k, v := range Points {
		if vutils.Equal(addr, v.Point.UDPAddr) &&
			(0 != v.Point.State&vomni.PointStateSigned) &&
			(0 == v.Point.State&vomni.PointStateDisconnected) {
			fmt.Printf("vk-xxx >>>>>>>>>>>>>> %s <<<<<<< need 2 disconnect %+v\n", k, addr)
			fmt.Printf("vk-xxx >>>>>>>>>>>>>> %s <<<<<<< need 2 disconnect %+v\n", k, addr)
			fmt.Printf("vk-xxx >>>>>>>>>>>>>> %s <<<<<<< need 2 disconnect %+v\n", k, addr)

			pt := Points[k]
			ptPt := Points[k].Point
			ptPt.State |= vomni.PointStateDisconnected
			pt.Point = ptPt
			Points[k] = pt

			str := fmt.Sprintf("Point %q lost connection", k)

			vutils.LogInfo(str)
			// send disconnection message to all configurations of the point
			for _, v := range Points[k].Run {
				v.LogStr(vomni.LogFileCdErr, str)
				v.GetDone(vomni.DoneDisconnected)
			}
		}
	}

	return
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
