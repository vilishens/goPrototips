package udp

import (
	"fmt"
	"net"
	"strconv"
	vomni "vk/omnibus"
	vparams "vk/params"
	vutils "vk/utils"
)

func Server(chGoOn chan bool, chDone chan int, chErr chan error) {

	lDone := make(chan bool)
	lErr := make(chan error)
	//	rDone := make(chan bool)
	//	rErr := make(chan error)

	chGoOn <- true

	addr := net.UDPAddr{
		Port: vparams.Params.PortUDPInternal,
		IP:   net.ParseIP(vparams.Params.InternalIPv4),
	}

	conn, err := net.ListenUDP("udp", &addr)

	chErr <- fmt.Errorf("kika")

	if err != nil {
		// Something really wrong - let's stop immediately
		addrStr := addr.IP.String() + ":" + strconv.Itoa(addr.Port)

		err = vutils.ErrFuncLine(fmt.Errorf("Couldn't get connection of %s --- %v", addrStr, err))
		chErr <- err
		return
	}

	defer conn.Close()

	//xxx	go runRotate()

	//	go runRotate(rDone)
	//	<-rDone
	//	if err = <-rErr; nil == err {
	//		chErr <- vutils.ErrFuncLine(fmt.Errorf("\nSomething wrong with point rotation -- %v", err))
	//	}

	//xxx	go sendMessages(rDone, rErr)

	//xxx	go waitMsg(conn, lDone, lErr)

	chGoOn <- true
	select {
	case <-lDone:
		break
	case err := <-lErr:
		vomni.RootErr <- err
		break
	}
}

//##################################################
//##################################################
//##################################################

/*
func Server(chGoOn chan bool, chDone chan int, chErr chan error) {

	lDone := make(chan bool)
	lErr := make(chan error)
	rDone := make(chan bool)
	rErr := make(chan error)

	addr := net.UDPAddr{
		Port: vparam.Params.InternalPort,
		IP:   net.ParseIP(vparam.Params.InternalIPv4),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		// Something really wrong - let's stop immediately
		addrStr := addr.IP.String() + ":" + strconv.Itoa(addr.Port)

		err = vutils.ErrFuncLine(fmt.Errorf("Couldn't get connection of %s --- %v", addrStr, err))
		chErr <- err
		return
	}

	defer conn.Close()

	go runRotate()

	//	go runRotate(rDone)
	//	<-rDone
	//	if err = <-rErr; nil == err {
	//		chErr <- vutils.ErrFuncLine(fmt.Errorf("\nSomething wrong with point rotation -- %v", err))
	//	}

	go sendMessages(rDone, rErr)

	go waitMsg(conn, lDone, lErr)

	chGoOn <- true
	select {
	case <-lDone:
		break
	case err := <-lErr:
		vomni.RootErr <- err
		break
	}
}

func runRotate() {

	if 0 >= vparam.Params.RotateRunSecs {
		vomni.RootErr <- vutils.ErrFuncLine(fmt.Errorf("\nWrong point rotation interval %d", vparam.Params.RotateRunSecs))
		return
	}

	for {
		timer := time.NewTimer(time.Duration(vparam.Params.RotateRunSecs) * time.Second)
		//time.Duration(hour*3600+min*60+sec) * time.Second

		// rotate
		if err := setRotateFiles(); nil != err {
			//if err := vutils.RunRotate(vparam.Params.RotateRunCfg); nil != err {
			vomni.RootErr <- vutils.ErrFuncLine(fmt.Errorf("\nRotation command failure -- %v", err))
			return
		}

		vutils.LogStr(vomni.LogInfo, "Rotate check")

		timeStr := time.Now().Format("2006-01-02 15:04:05 -07:00 MST")
		str := fmt.Sprintf("==>>>>>\n==>>>>>\n==>>>>>\n==>>>>>\n %s <<<<<< ROTATE\n==>>>>>\n==>>>>>\n==>>>>>", timeStr)

		fmt.Println(str)

		select {
		case <-timer.C:

			fmt.Println("\n\n\n$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
			fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
			fmt.Printf("$$$$$$$$$$$$$$$$$$$$$$ %q $$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$\n", timeStr)
			fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
			fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$\n\n\n")
			//			timeStr := time.Now().Format("2006-01-02 15:04:05 -07:00 MST")
			//			str := fmt.Sprintf("==>>>>>\n==>>>>>\n==>>>>>\n==>>>>>\n %s <<<<<< ROTATE\n==>>>>>\n==>>>>>\n==>>>>>", timeStr)

			//			fmt.Println(str)
		}
	}
}

func setRotateFiles() (err error) {
	// rotate files if necessary
	if err = vutils.RunRotate(vparam.Params.RotateRunCfg); nil != err {
		vomni.RootErr <- vutils.ErrFuncLine(fmt.Errorf("\nRotation command failure -- %v", err))
		return
	}

	// reassign the main logger files
	if err = reassingMainFile(); nil != err {
		vomni.RootErr <- vutils.ErrFuncLine(fmt.Errorf("\nRotation main file reassign failure -- %v", err))
		return
	}

	// reassign point logger files
	if err = reassignPointFiles(); nil != err {
		vomni.RootErr <- vutils.ErrFuncLine(fmt.Errorf("\nRotation point files reassign failure -- %v", err))
		return
	}

	fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	fmt.Println("@@@@@@@@@@@@@@@@@@@@@@ ZAPAH-ZAPAH-ZAPAH --> RESTART @@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")

	return
}

func reassingMainFile() (err error) {
	if vomni.LogMainFile, err = vutils.LogReAssign(vomni.LogMainFile, vomni.LogMainPath); nil != err {
		vomni.RootErr <- vutils.ErrFuncLine(fmt.Errorf("\nRotation file reaasign failure -- %v", err))
		return
	}

	vomni.LogData.SetOutput(vomni.LogMainFile)
	vomni.LogErr.SetOutput(vomni.LogMainFile)
	vomni.LogFatal.SetOutput(vomni.LogMainFile)
	vomni.LogInfo.SetOutput(vomni.LogMainFile)

	return
}

func reassignPointFiles() (err error) {

	if err = xrun.RotateReAssign(); nil != err {
		return err
	}

	return
}

func waitMsg(conn *net.UDPConn, done chan bool, chErr chan error) {

	for {
		// waiting, waiting, ... UDP
		buff := make([]byte, 2048)
		nn, msgAddr, err := conn.ReadFromUDP(buff)
		if err != nil {
			continue
		}

		if len(buff) == 0 {
			continue
		} else {
			fmt.Printf("KASATONICH ##### %3d Read a message from %v \"%s\" \n", nn, msgAddr, string(buff))
		}

		locErr := make(chan error)
		go xrun.MessageGet(string(buff[:nn]), locErr)

		err = <-locErr

		if err != nil {
			chErr <- vutils.ErrFuncLine(err)
			return
		}

		buff = []byte{}
	}
}

func sendMessages(done chan bool, chErr chan error) {

	for {
		if len(xmsg.SendList) == 0 {
			time.Sleep(msgSendListEmptyDelay)
		}

		for i := 0; i < len(xmsg.SendList); {
			time.Sleep(msgSendListDelay)

			if i >= len(xmsg.SendList) {
				continue
			}

			chDone := make(chan bool)

			if (nil == xmsg.SendList[i]) || ("" == xmsg.SendList[i].Msg) {
				xmsg.SendList.MinusIndex(i, chDone)
				<-chDone
				continue
			}

			if time.Since(xmsg.SendList[i].Last) < msgSendListInterval {
				continue
			}

			xmsg.SendList[i].Repeat++
			if xmsg.SendList[i].Repeat >= MsgSendRepeatLimit {
				go xmsg.SendList.MinusIndex(i, chDone)
				<-chDone

				continue
			}

			xmsg.SendList[i].Last = time.Now()
			if err := SendToAddress(xmsg.SendList[i].UDPAddr, xmsg.SendList[i].Msg); nil != err {
				chErr <- err
				return
			}
			i++
		}

		time.Sleep(msgSendListDelay)
	}
}

func SendToAddress(addr net.UDPAddr, msg string) (err error) {
	addrStr := addr.IP.String() + ":" + strconv.Itoa(addr.Port)
	return sendToAddress(addrStr, msg)
}

func sendToAddress(addr string, msg string) (err error) {

	conn, err := net.Dial("udp", addr)
	if err != nil {
		err = vutils.ErrFuncLine(fmt.Errorf("Connection ERROR: %v", err))
		return
	}
	defer conn.Close()

	if _, err = conn.Write([]byte(msg)); nil != err {
		err = vutils.ErrFuncLine(fmt.Errorf("SentToAddress ERROR: %v", err))
	}

	return
}
*/
