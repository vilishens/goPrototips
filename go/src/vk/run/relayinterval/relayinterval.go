package runrelayinterval

import (
	"fmt"
	"net"
	"strconv"
	vmsg "vk/messages"
	vomni "vk/omnibus"
	vutils "vk/utils"
)

var RunningPoints map[string]RunData

func init() {
	RunningPoints = make(map[string]RunData)
}

func (d RunData) Starter(flds []string, chGoOn chan bool, chErr chan error) {

	fmt.Println("$$$$$$$$$$$$$$$$ FINAL $$$$$$$$$$$$$$$$$")

	intNbr, err := strconv.Atoi(flds[vomni.MsgIndexPrefixNbr])
	if nil != err {
		err = vutils.ErrFuncLine(fmt.Errorf("Point %q received a message (%v) with the wrong Number string %q - %s",
			d.Point,
			flds,
			flds[vomni.MsgIndexPrefixNbr],
			err.Error()))
		vutils.LogErr(err)
	}

	if nil == err {
		vmsg.MessageMinusByNbr(intNbr)

		if d.PointUDPAddr(flds) && (0 == (d.State & vomni.PointStateActive)) {
			fmt.Printf("============ UDPAddr %+v NBR %d\n", d.UDPAddr, intNbr)

			d.State |= vomni.PointStateActive
		}
	}

	chGoOn <- true
}

func (d *RunData) PointUDPAddr(flds []string) (ok bool) {

	intPort, err := strconv.Atoi(flds[vomni.MsgPrefixLen+vomni.MsgIndexHelloFromPointPort])
	if nil != err {
		err = vutils.ErrFuncLine(fmt.Errorf("Point %q received a message (%v) with the wrong Port format %q - %s",
			d.Point,
			flds,
			flds[vomni.MsgPrefixLen+vomni.MsgIndexHelloFromPointPort],
			err.Error()))
		vutils.LogErr(err)
	}

	netIP := net.ParseIP(flds[vomni.MsgPrefixLen+vomni.MsgIndexHelloFromPointIP])
	if nil == netIP {
		err = vutils.ErrFuncLine(fmt.Errorf("Point %q received a message (%v) with the invalid IP %q",
			d.Point,
			flds,
			flds[vomni.MsgPrefixLen+vomni.MsgIndexHelloFromPointIP]))
		vutils.LogErr(err)
	}

	if nil != err {
		return false
	}

	d.UDPAddr = net.UDPAddr{IP: netIP, Port: intPort}

	return true
}

func (d RunData) PointUDPAddrY(flds []string) (addr net.UDPAddr) {

	intPort, err := strconv.Atoi(flds[vomni.MsgPrefixLen+vomni.MsgIndexHelloFromPointPort])
	if nil != err {
		err = vutils.ErrFuncLine(fmt.Errorf("Point %q received a message (%v) with the wrong Port format %q - %s",
			d.Point,
			flds,
			flds[vomni.MsgPrefixLen+vomni.MsgIndexHelloFromPointPort],
			err.Error()))
		vutils.LogErr(err)
	}

	netIP := net.ParseIP(flds[vomni.MsgPrefixLen+vomni.MsgIndexHelloFromPointIP])
	if nil == netIP {
		err = vutils.ErrFuncLine(fmt.Errorf("Point %q received a message (%v) with the invalid IP %q",
			d.Point,
			flds,
			flds[vomni.MsgPrefixLen+vomni.MsgIndexHelloFromPointIP]))
		vutils.LogErr(err)
	}

	if nil != err {
		return //false
	}

	return net.UDPAddr{IP: netIP, Port: intPort}

	//	return //true
}
