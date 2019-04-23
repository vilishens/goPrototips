package netinfo

import (
	"fmt"
	"net"
	"strings"
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
			vutils.LogErr(err)
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

	first := true
	repeatMax := vomni.NetInfoRepeats
	netDuration := vomni.NetInfoInterval

	repeat := 0
	for {
		intIP, extIP, errCd, err := getIPv4Addrs()

		if 0 != (vomni.NoNetInternal & errCd) {
			// no sense to continue if there is no the internal net
			err = vutils.ErrFuncLine(fmt.Errorf("Couldn't find the internal IP - %s", err.Error()))

			chErr <- err
			return
		} else if 0 != (vomni.NoNetExternal & errCd) {
			repeat++

			err = vutils.ErrFuncLine(fmt.Errorf("Couldn't find the external IP (attempt #%d) - %s", repeat, err.Error()))

			if repeat <= repeatMax {
				vutils.LogErr(err)
				continue
			}

			if 0 != (vomni.NetExternalRequired & vparams.Params.NetExternalRequirement) {
				// in case of absence of the required external net stop with error
				chErr <- err
				return
			}
		}

		repeat = 0
		if setCurrentNet(intIP, extIP) && (0 < (vparams.Params.NetExternalRequirement & vomni.NetExternalBits)) {
			// send IP email only if the external net required or nice to have

			// Ko darīt ja nevajag ārējo IP adresi? vai jāsūta emails?
			// Liekas, ka nevajag sūtīt email, ja nevajag ārejo tīklu
			// Jāņem vērā arī NetExternalTreatment
			/*
				if err = sendNetInfov4(); nil != err {
					err = vutils.ErrFuncLine(fmt.Errorf("Couldn't send new IPv4 - %v", err))

					vomni.LogFatal.Println(err)

					vomni.RootErr <- err
					done <- vomni.DoneOK
				}
			*/
		}

		str := `=== Iternal IP "` + vparams.Params.IPAddressInternal + `" External IP "` + vparams.Params.IPAddressExternal + `"`
		fmt.Println(str)

		if first {
			chGoOn <- true
			first = false
		}

		tick := time.NewTicker(netDuration)
		<-tick.C
	}
}

func setCurrentNet(intIP string, extIP string) (newIP bool) {

	wasInternal := vparams.Params.IPAddressInternal
	wasExternal := vparams.Params.IPAddressExternal

	if (nil != net.ParseIP(intIP)) && (vparams.Params.IPAddressInternal != intIP) {
		vparams.Params.IPAddressInternal = intIP
		newIP = true
	}

	if (nil != net.ParseIP(extIP)) && (vparams.Params.IPAddressExternal != extIP) {
		vparams.Params.IPAddressExternal = extIP
		newIP = true
	}

	if newIP {
		str := fmt.Sprintf("New IP: Internal %q (was %q)", vparams.Params.IPAddressInternal, wasInternal)
		if 0 != (vomni.NetExternalBits & vparams.Params.NetExternalRequirement) {
			str += fmt.Sprintf(", External %q (was %q)", vparams.Params.IPAddressExternal, wasExternal)
		}

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

func getIPv4Addrs() (intIP string, extIP string, errCd int, err error) {
	if intIP, err = vutils.InternalIPv4(); nil != err {
		err = vutils.ErrFuncLine(fmt.Errorf("Couldn'get Internal IPv4 - %v", err))
		errCd = vomni.NoNetInternal
		return
	}

	if 0 != (vomni.NetExternalBits & vparams.Params.NetExternalRequirement) {
		if extIP, err = ExternalIPv4(); nil != err {
			err = vutils.ErrFuncLine(fmt.Errorf("Couldn't get External IPv4 - %v", err))
			errCd = vomni.NoNetExternal
			return
		}
	}

	return
}

func ExternalIPv4() (ip string, err error) {

	cmds := vparams.Params.IPExternalAddressCmds

	tmpIP := ""

	ind := 0
	for "" == tmpIP && nil == err && ind < len(cmds) {
		tick := time.NewTicker(2 * time.Second)
		cmd := cmds[ind]
		chStr := make(chan string)
		chErr := make(chan error)

		go vutils.DoCmd(cmd, chStr, chErr)
		select {
		case <-tick.C:
			ind++
		case tmpIP = <-chStr:
			return strings.Trim(tmpIP, "\n"), nil
		case err = <-chErr:
			err1 := vutils.ErrFuncLine(fmt.Errorf("Failed CMD \"%s\" --- %v (index %d)", cmd, err, ind))

			vomni.LogErr.Println(err1)
			ind++
		}
	}

	if "" == tmpIP {
		return "", vutils.ErrFuncLine(fmt.Errorf("Couldn't get the external IP"))
	}

	return
}
