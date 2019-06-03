package runrelayinterval

import (
	"fmt"
	"net"
	"strconv"
	"time"
	vmsg "vk/messages"
	vomni "vk/omnibus"
	vcfg "vk/pointconfig"
	vutils "vk/utils"
)

var RunningPoints map[string]RunData

func init() {
	RunningPoints = make(map[string]RunData)
}

func (d RunData) LetsGo(addr net.UDPAddr, flds []string, chGoOn chan bool, chDone chan int, chErr chan error) {

	fmt.Println("$$$$$$$$$$$$$$$$ FINAL $$$$$$$$$$$$$$$$$", d.UDPAddr)

	d.UDPAddr = addr

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

		if 0 == (d.State & vomni.PointStateSigned) {
			fmt.Printf("============ UDPAddr %+v NBR %d\n", d.UDPAddr, intNbr)

			d.Index = AllIndex{Start: -1, Base: -1, Finish: -1}

			locGoOn := make(chan bool)
			locDone := make(chan int)
			locErr := make(chan error)
			go d.run(locGoOn, locDone, locErr)

			<-locGoOn

			d.State |= vomni.PointStateActive
			d.State |= vomni.PointStateSigned
		}
	}

	chGoOn <- true
}

func (d RunData) run(chGoOn chan bool, chDone chan int, chErr chan error) {

	fmt.Printf("Point %q Addr %+v Index %+v\n", d.Point, d.UDPAddr, d.Index)

	chGoOn <- true

	locDone := make(chan int)

	once := true
	go d.runArray(d.Cfg.Start, locDone, &d.Index.Start, once)
	rc := <-locDone

	once = false
	go d.runArray(d.Cfg.Base, locDone, &d.Index.Base, once)
	rc = <-locDone

	once = true
	go d.runArray(d.Cfg.Finish, locDone, &d.Index.Finish, once)
	rc = <-locDone

	_ = rc
}

func (d RunData) runArray(arr vcfg.RelIntervalArray, chDone chan int, index *int, once bool) {

	if 0 == len(arr) {
		chDone <- vomni.DoneStop
		return
	}

	*index = nextIndex(*index, len(arr))

	for {
		tick := time.NewTicker(arr[*index].Seconds)

		t := time.Now()
		fmt.Println(d.Point, "@@@@@@@@@@@@@@@@", t.Format(vomni.TimeFormat1), "*************** INDEX ", *index, "JĀSŪTA CMD PIRMS INTERVALA", arr[*index].Seconds.Seconds())

		select {
		case <-tick.C:
			*index = nextIndex(*index, len(arr))

			if once && 0 == *index {
				chDone <- vomni.DoneStop
				return
			}
		}

	}
	chDone <- vomni.DoneStop

}

func nextIndex(ind int, count int) (index int) {

	index = ind + 1

	if (index < 0) || (index >= count) {
		index = 0
	}

	return
}