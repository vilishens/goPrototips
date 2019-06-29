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

var Points map[string]*PointRun
var listSigned map[string]net.UDPAddr
var startSequence []int

func init() {
	Points = make(map[string]*PointRun)
	listSigned = make(map[string]net.UDPAddr)
	startSequence = []int{vomni.CfgTypeRelayInterval}
}

func Runners() {
	relayIntervalRunners()
}

func relayIntervalRunners() {
	for k, v := range vrunrelayinterval.RunningPoints {

		if _, has := Points[k]; !has {
			// it is required to create a new point running object from the template
			addNewPointRun(k)
		}

		// set the type of the Point
		tPoint := Points[k].Point
		tPoint.Type |= v.Type

		// save the the configuration data
		tRun := Points[k].Run
		tRun[v.Type] = v

		// put all data into the point running object
		Points[k].Point = tPoint
		Points[k].Run = tRun
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

	listHandled := make(map[string]bool) // list of signed already handled points
	listStart := make(map[string]bool)   // list of points in start state

	for _, cfgType := range startSequence {
		// start all point configuration, sequence set in startSequence
		// Sequence can be important some times (for instance, to check the point ready state)

		fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!! Nothing will stop them !!!!!!!!!!!!!!!!!!")

		for point, addr := range listSigned {
			logStr := ""
			pData := new(PointRun)
			ok := false

			if pData, ok = Points[point]; !ok {
				err := fmt.Errorf("The point %q (%v) sent SignIn message, but there is no configuration of this point", point, addr)
				vutils.LogErr(err)
				chErr <- vutils.ErrFuncLine(err)

				return
			}

			fmt.Printf("!!!!!!!!!!!!!!!!!!!!!!!! %s !!!!!!!!!!!!!!!!!!\n", point)

			// Point handling
			if _, ok := listHandled[point]; !ok {
				// the point isn't handled yet

				// save the point address
				pData.setUDPAddr(addr)

				if 0 != (pData.Point.State & vomni.PointStateDisconnected) {
					// this point was signed in, but later disconnected
					// need to restart again
					logStr = fmt.Sprintf("START SIGNED *** Point %q signed in AGAIN", point)
					listStart[point] = true // need to restart
				} else if 0 == (pData.Point.State & vomni.PointStateSigned) {
					// the point wasn't signed in, need to start from scratch
					logStr = fmt.Sprintf("START SIGNED *** Point %q signed in", point)
					listStart[point] = true
				} else {
					// the point was signed and not disconnected, to update the address is enough
					logStr = fmt.Sprintf("START SIGNED *** Point %q (signed in already) saves the new UDP address %s:%d", point, addr.IP.String(), addr.Port)
				}

				// put messages about signed in into log
				vutils.LogInfo(logStr)

				// set the clean signed state
				pData.setState(vomni.PointStateDisconnected, false)
				pData.setState(vomni.PointStateSigned, true)

				listHandled[point] = true
			}

			//#####################################

			if 0 != (pData.Point.Type & cfgType) {
				pCfg := pData.Run[cfgType]
				state := pCfg.GetState()
				startCfg := false // start configuration

				if !pCfg.Ready() {
					// the configuration of this point is not ready
					strState := ""
					if 0 != (vomni.PointCfgStateUnavailable & state) {
						// this configuration has been unavailable already
						// no log message required
					} else if vomni.PointCfgStateUnknown == state {
						strState = "not started yet"
					} else if 0 != (vomni.PointCfgStateReady & state) {
						strState = "was ready"
					}

					if strState != "" {
						logStr = fmt.Sprintf("The point %q (%s) configuration %q is not ready",
							point,
							strState,
							vomni.PointCfgData[cfgType].CfgStr)

						vutils.LogErr(fmt.Errorf("%s", logStr))

						// send log to the point configuration
						// (it succeeds only if the point configuration was ready (rotate files were started))
						pCfg.LogStr(vomni.LogFileCdErr, logStr)
					}

					pCfg.SetState(vomni.PointCfgStateReady, false)
					pCfg.SetState(vomni.PointCfgStateUnavailable, false)

					continue
				} else {
					strState := ""

					// the point configuration is ready
					if vomni.PointCfgStateUnknown == state {
						// this the very 1st start of the configuration
						startCfg = true
						strState = "wasn't started yet"

						// start rotation of the log files
						if err := pCfg.StartRotate(); nil != err {
							chErr <- err
							return
						}
					}

					if 0 != (vomni.PointCfgStateUnavailable & state) {
						startCfg = true
						strState = "was unavailable"
					}

					if strState != "" {
						logStr = fmt.Sprintf("The point %q (%s) configuration %q is ready",
							point,
							strState,
							vomni.PointCfgData[cfgType].CfgStr)

						vutils.LogInfo(logStr)

						// send log to the point configuration
						// (it succeeds only if the point configuration was ready (rotate files were started))
						pCfg.LogStr(vomni.LogFileCdInfo, logStr)
					}

					// remember this configuration state
					pCfg.SetState(vomni.PointCfgStateUnavailable, false)
					pCfg.SetState(vomni.PointCfgStateReady, true)
				}

				fmt.Printf("ooooooooooooooooo Point %q has type 0x%06x ooooooooooooooooooooooo\n", point, cfgType)

				if listStart[point] || startCfg {
					locGoOn := make(chan bool)
					locDone := make(chan int)
					locErr := make(chan error)

					pCfg.SetUDPAddr(addr)

					go pCfg.LetsGo(locGoOn, locDone, locErr)

					select {
					case <-locGoOn:
					case cd := <-locDone:
						chDone <- cd
					case err := <-locErr:
						chErr <- err
					}

					fmt.Println("vk-xxx LOMBARDS ", Points[point].Run[cfgType].GetUDPAddr())

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

		fmt.Println("vk-xxx Colombus")

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

			Points[k].setState(vomni.PointStateDisconnected, true)
			str := fmt.Sprintf("Point %q lost connection", k)

			vutils.LogInfo(str)
			// send disconnection code to all configurations of the point
			for _, v := range Points[k].Run {
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

func (d *PointRun) setState(state int, on bool) {
	if on {
		d.Point.State |= state
	} else {
		d.Point.State &^= state
	}
}

func (d *PointRun) setUDPAddr(addr net.UDPAddr) {
	d.Point.UDPAddr = addr
}

func addNewPointRun(point string) {
	newP := new(PointRun)

	newP.Point.Point = point
	newP.Point.State = vomni.PointStateUnknown
	newP.Point.Type = vomni.CfgTypeUnknown
	newP.Point.UDPAddr = net.UDPAddr{}

	newP.Run = make(map[int]Runner)

	Points[point] = newP
}
