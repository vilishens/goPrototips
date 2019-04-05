package netinfo

import (
	"fmt"
	"net"
	"time"
	vomni "vk/omnibus"

	//	vsgrid "vk/sendgrid"
	vparams "vk/params"
	vutils "vk/utils"
)

func NetInfo(chGoOn chan bool, chDone chan int, chErr chan error) {

	locGoOn := make(chan bool)
	locDone := make(chan int)
	locErr := make(chan error)

	go netInfo(locGoOn, locDone, locErr)

	stop := false
	for !stop {
		select {
		case <-locGoOn:
			chGoOn <- true
		case err := <-locErr:
			err = vutils.ErrFuncLine(err)
			vutils.LogStr(vomni.LogErr, err.Error())
			chErr <- err
			stop = true
			return
		case done := <-locDone:
			chDone <- done
			stop = true
			return
		}
	}
}

func netInfo(chGoOn chan bool, chDone chan int, chErr chan error) {

	netDuration := netInfoInterval
	first := true

	repeatMax := netInfoRepeats

	repeat := 0
	for {
		if first {
			chGoOn <- true
			first = false
		} else {
			tick := time.NewTicker(netDuration)
			<-tick.C
		}

		intIP, extIP, err := getIPv4Addrs()
		if nil != err {
			repeat++

			if repeat <= repeatMax {
				err = vutils.ErrFuncLine(fmt.Errorf("Couldn't get current IP addresses (attempt #%d)", repeat))
				vutils.LogStr(vomni.LogErr, err.Error())
				continue
			}

			err := vutils.ErrFuncLine(fmt.Errorf("Couldn't find IP (int \"%s\", ext \"%s\", reapeat number %d)",
				intIP, extIP, repeat))
			vutils.LogStr(vomni.LogErr, err.Error())

			chErr <- err
			return
		}

		repeat = 0
		if setCurrentNet(intIP, extIP) {
			/*
				if err = sendNetInfov4(); nil != err {
					err = vutils.ErrFuncLine(fmt.Errorf("Couldn't send new IPv4 - %v", err))

					vomni.LogFatal.Println(err)

					vomni.RootErr <- err
					done <- vomni.DoneOK
				}
			*/
		}
	}
}

func setCurrentNet(intIP string, extIP string) (send bool) {

	wasInternal := vparams.Params.IPAddressInternal
	wasExternal := vparams.Params.IPAddressExternal

	if (nil != net.ParseIP(intIP)) && (vparams.Params.IPAddressInternal != intIP) {
		vparams.Params.IPAddressInternal = intIP
		send = true
	}

	if (nil != net.ParseIP(extIP)) && (vparams.Params.IPAddressExternal != extIP) {
		vparams.Params.IPAddressExternal = extIP
		send = true
	}

	if send {
		str := fmt.Sprintf("New IP addresses: Internal %q (was %q), External %q (was %q)",
			vparams.Params.IPAddressInternal, wasInternal, vparams.Params.IPAddressExternal, wasExternal)
		vutils.LogStr(vomni.LogInfo, str)
	}

	return
}

/*
func sendNetInfov4() (err error) {

	emails := vparam.Params.WebEmail

	subj := vparam.Params.Name + " --- " + vutils.TimeNow(vomni.TimeFormat1) + " --- NET"
	msg_txt := fmt.Sprintf("WEB: %s:%d\nSSH: %s:%d\n\n",
		vparam.Params.ExternalIPv4, 50177,
		vparam.Params.ExternalIPv4, 50354)
	msg_html := fmt.Sprintf("</h2>WEB: <strong>%s:%d</strong><br />SSH: <strong>%s:%d<strong><br /><br /></h2>",
		vparam.Params.ExternalIPv4, 50177, //vparam.Params.InternalPort,
		vparam.Params.ExternalIPv4, 50354) //vparam.Params.ExternalPort)

	return vsgrid.Send(emails, subj, msg_txt, msg_html)
}
*/
func getIPv4Addrs() (intIP string, extIP string, err error) {
	if intIP, err = vutils.InternalIPv4(); nil != err {
		err = vutils.ErrFuncLine(fmt.Errorf("Couldn'get Internal IPv4 - %v", err))
		return
	}

	/*
		if extIP, err = vutils.ExternalIPv4(); nil != err {
			err = vutils.ErrFuncLine(fmt.Errorf("Couldn't get External IPv4 - %v", err))
			return
		}
	*/

	return
}

/*
func ExternalIPv4() (ip string, err error) {

	cmds := []string{"dig +short myip.opendns.com @resolver1.opendns.com",
		"curl http://myip.dnsomatic.com"}

	tmpIP := ""

	ind := 0
	for "" == tmpIP && nil == err && ind < len(cmds) {
		tick := time.NewTicker(2 * time.Second)
		cmd := cmds[ind]
		chStr := make(chan string)
		chErr := make(chan error)

		go doCmd(cmd, chStr, chErr)
		select {
		case <-tick.C:
			ind++
		case tmpIP = <-chStr:
			return strings.Trim(tmpIP, "\n"), nil
		case err = <-chErr:
			// return "", ErrFuncLine(fmt.Errorf("Failed CMD \"%s\" --- %v", cmd, err))
			err1 := ErrFuncLine(fmt.Errorf("Failed CMD \"%s\" --- %v (index %d)", cmd, err, ind))

			vomni.LogErr.Println(err1)

			fmt.Println("########### OSET ### jāliek kļūdas ieraksts %v", err1)
			ind++
		}
	}

	if "" == tmpIP {
		return "", ErrFuncLine(fmt.Errorf("Couldn't get the external IP"))
	}

	return
}
*/
