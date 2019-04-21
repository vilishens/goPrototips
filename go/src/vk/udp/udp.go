package udp

import (
	"fmt"
	"net"
	"strconv"
	"time"
	vmsg "vk/messages"
	vomni "vk/omnibus"
	vparams "vk/params"
	vpointrun "vk/run/pointrun"
	vutils "vk/utils"
)

func Server(chGoOn chan bool, chDone chan int, chErr chan error) {

	addr := net.UDPAddr{
		Port: vparams.Params.PortUDPInternal,
		IP:   net.ParseIP(vparams.Params.IPAddressInternal),
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

	sendDone := make(chan int)
	sendErr := make(chan error)

	go sendMessages(sendDone, sendErr)

	waitDone := make(chan int)
	waitErr := make(chan error)

	go waitMsg(conn, waitDone, waitErr)

	chGoOn <- true
	select {
	case cd := <-sendDone:
		vutils.LogInfo(fmt.Sprintf("UDP finished with Send RC %d", cd))
	case cd := <-waitDone:
		vutils.LogInfo(fmt.Sprintf("UDP finished with Send RC %d", cd))
	case err := <-sendErr:
		vutils.LogErr(err)
		vomni.RootErr <- vutils.ErrFuncLine(err)
	case err := <-waitErr:
		vutils.LogErr(err)
		vomni.RootErr <- vutils.ErrFuncLine(err)
	}
}

func sendMessages(done chan int, chErr chan error) {

	for {
		time.Sleep(vomni.DelaySendMessage)

		// vk-xxx šitas jāaizvāc
		continue

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

			if !vmsg.MessageList2Send[i].Last.IsZero() && time.Since(vmsg.MessageList2Send[i].Last) < vomni.DelaySendMessageRepeat {
				// this is a repeated message but the repeat interval isn't passed yet
				continue
			}

			vmsg.MessageList2Send[i].Repeat++

			if vmsg.MessageList2Send[i].Repeat > vomni.MessageSendRepeatLimit {

				vutils.LogInfo(fmt.Sprintf("Deleted message #%d due to the exceeded send repeat limit", vmsg.MessageList2Send[i].MessageNbr))

				go vmsg.MessageList2Send.MinusIndex(i, chDone)
				<-chDone
				continue
			}

			vmsg.MessageList2Send[i].Last = time.Now()

			// Jāatjauno servera laiks ziņojumā

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

	select {
	case err = <-chErr:
		return
	}

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

func waitMsg(conn *net.UDPConn, done chan int, chErr chan error) {

	for {
		// waiting, waiting, ... UDP
		time.Sleep(vomni.DelayWaitMessage)

		buff := make([]byte, 4096)
		nn, msgAddr, err := conn.ReadFromUDP(buff)
		if err != nil {
			continue
		}

		if len(buff) == 0 {
			continue
		}

		msg := string(buff[:nn])

		locErr := make(chan error)
		go vpointrun.MessageReceived(msg, locErr)
		err = <-locErr

		if err != nil {
			vutils.LogErr(fmt.Errorf("The received message %q (address %s:%d) %error %q", msg,
				msgAddr.IP.String(), msgAddr.Port, err.Error()))
		}

		buff = []byte{}
	}
}
