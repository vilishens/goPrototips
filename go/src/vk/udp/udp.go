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

		if len(vmsg.MessageList2Send) == 0 {
			// no messages to send
			time.Sleep(vomni.DelaySendMessageListEmpty)
			continue
		}

		for i := 0; i < len(vmsg.MessageList2Send); i++ {
			time.Sleep(vomni.DelaySendMessage)

			radit := "192.168.0.182" == vmsg.MessageList2Send[i].UDPAddr.IP.String()

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

			if radit {
				fmt.Println("###### PIRMS ", vmsg.MessageList2Send[i].Repeat)
			}

			was := vmsg.MessageList2Send[i].Repeat

			vmsg.MessageList2Send[i].Repeat++

			if radit {
				fmt.Println("###### Process  ", vmsg.MessageList2Send[i].Repeat, " WAS ", was, " NOW ", vmsg.MessageList2Send[i].Repeat)
			}

			if vmsg.MessageList2Send[i].Repeat >= vomni.MessageSendRepeatLimit {

				if radit {
					fmt.Println("###*************************************************### Paddington")
				}

				go vmsg.MessageList2Send.MinusIndex(i, chDone)
				vutils.LogInfo(fmt.Sprintf("Deleted message #%d due to the exceeded repeat limit"))
				<-chDone
				continue
			}

			vmsg.MessageList2Send[i].Last = time.Now()

			if radit {
				fmt.Println("***** mizandarI  ", vmsg.MessageList2Send[i].Repeat)
			}

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
