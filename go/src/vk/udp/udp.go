package udp

import (
	"fmt"
	"net"
	"strconv"
	"time"
	vmsg "vk/messages"
	vomni "vk/omnibus"
	vparams "vk/params"
	vutils "vk/utils"
)

func Server(chGoOn chan bool, chDone chan int, chErr chan error) {

	addr := net.UDPAddr{
		Port: vparams.Params.PortUDPInternal,
		IP:   net.ParseIP(vparams.Params.InternalIPv4),
	}

	conn, err := net.ListenUDP("udp", &addr)

	if err != nil {
		// Something really wrong - let's stop immediately
		addrStr := addr.IP.String() + ":" + strconv.Itoa(addr.Port)
		err = fmt.Errorf("Couldn't get connection of %s --- %v", addrStr, err)
		vutils.LogErr(err)
		chErr <- vutils.ErrFuncLine(err)
		return
	}

	defer conn.Close()

	sendDone := make(chan bool)
	sendErr := make(chan error)

	go sendMessages(sendDone, sendErr)

	//xxx	go waitMsg(conn, lDone, lErr)

	chGoOn <- true
	select {
	case cd := <-sendDone:
		vutils.LogInfo(fmt.Sprintf("UDP finished with Send RC %d", cd))
		break
	case err := <-sendErr:
		vutils.LogErr(err)
		vomni.RootErr <- vutils.ErrFuncLine(err)
		break
	}
}

func sendMessages(done chan bool, chErr chan error) {

	for {
		time.Sleep(vomni.DelaySendMessage)

		if len(vmsg.MessageList2Send) == 0 {
			// no messages to send
			time.Sleep(vomni.DelaySendMessageListEmpty)
			continue
		}

		for i := 0; i < len(vmsg.MessageList2Send); i++ {
			time.Sleep(vomni.DelaySendMessage)

			if i >= len(vmsg.MessageList2Send) {
				// verify the index isn't out of the list
				continue
			}

			chDone := make(chan bool)

			if "" == vmsg.MessageList2Send[i].Msg {
				// this is the blank message no need to try to send just remove it
				vutils.LogInfo(fmt.Sprintf("Deleted blank message #%d", vmsg.MessageList2Send[i].MessageNbr))
				go vmsg.MessageList2Send.MinusIndex(i, chDone)
				<-chDone
				continue
			}

			if time.Since(vmsg.MessageList2Send[i].Last) < vomni.DelaySendMessageRepeat {
				// this is a repeated message but the repeat interval isn't passed yet
				continue
			}

			vmsg.MessageList2Send[i].Repeat++
			if vmsg.MessageList2Send[i].Repeat >= vomni.MessageSendRepeatLimit {
				go vmsg.MessageList2Send.MinusIndex(i, chDone)
				vutils.LogInfo(fmt.Sprintf("Deleted message #%d due to the exceeded repeat limit"))
				<-chDone
				continue
			}

			vmsg.MessageList2Send[i].Last = time.Now()
			if err := SendToAddress(vmsg.MessageList2Send[i].UDPAddr, vmsg.MessageList2Send[i].Msg); nil != err {
				// write the error in log
				vutils.LogErr(err)
			}
		}
	}
}

func SendToAddress(addr net.UDPAddr, msg string) (err error) {
	addrStr := addr.IP.String() + ":" + strconv.Itoa(addr.Port)

	chErr := make(chan error)
	go sendToAddress(addrStr, msg, chErr)

	return <-chErr
}

func sendToAddress(addr string, msg string, chErr chan error) (err error) {

	conn, err := net.Dial("udp", addr)
	if err != nil {
		err = vutils.ErrFuncLine(fmt.Errorf("Connection ERROR: %v", err))
		chErr <- err
		return
	}

	if nil == err {
		defer conn.Close()

		if _, err = conn.Write([]byte(msg)); nil != err {
			err = vutils.ErrFuncLine(fmt.Errorf("SentToAddress ERROR: %v", err))
		}
	}

	chErr <- err
	return
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



*/
