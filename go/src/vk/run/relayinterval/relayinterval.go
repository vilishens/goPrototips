package runrelayinterval

import (
	"fmt"
	"net"
	"time"
	vmsg "vk/messages"
	vomni "vk/omnibus"
	vcfg "vk/pointconfig"
	vrotate "vk/rotate"
	vutils "vk/utils"
)

var RunningPoints map[string]RunData

func init() {
	RunningPoints = make(map[string]RunData)
}

func (d RunData) LogStr(infoCd int, str string) {

	for _, v := range d.Logs {
		for k1, v1 := range v.Loggers {
			if k1 == infoCd {
				vutils.LogStr(v1.Logger, str)
			}
		}
	}
}

func (d RunData) LetsGo(addr net.UDPAddr, chGoOn chan bool, chDone chan int, chErr chan error) {

	fmt.Println("$$$$$$$$$$$$$$$$ FINAL $$$$$$$$$$$$$$$$$", d.UDPAddr)

	d.UDPAddr = addr

	fmt.Printf("============ UDPAddr %+v\n", d.UDPAddr)

	d.Index = AllIndex{Start: vomni.PointNonActiveIndex, Base: vomni.PointNonActiveIndex, Finish: vomni.PointNonActiveIndex}

	locGoOn := make(chan bool)
	locDone := make(chan int)
	locErr := make(chan error)
	go d.run(locGoOn, locDone, locErr)

	<-locGoOn

	d.State |= vomni.PointStateActive
	d.State |= vomni.PointStateSigned

	chGoOn <- true
}

func (d RunData) GetDone(done int) {
	d.ChDone <- done
}

func (d RunData) run(chGoOn chan bool, chDone chan int, chErr chan error) {

	fmt.Printf("Point %q Addr %+v Index %+v\n", d.Point, d.UDPAddr, d.Index)

	chGoOn <- true

	locDone := make(chan int)
	type stage struct {
		once  bool
		index *int
		cfg   vcfg.RelIntervalArray
	}

	allStages := []stage{stage{once: true, index: &d.Index.Start, cfg: d.Cfg.Start},
		stage{once: false, index: &d.Index.Base, cfg: d.Cfg.Base},
		stage{once: true, index: &d.Index.Finish, cfg: d.Cfg.Finish}}

	for _, v := range allStages {
		go d.runArray(v.cfg, v.index, v.once, locDone)
		rc := <-locDone
		if vomni.DoneDisconnected == rc {
			return
		}
	}

	/*
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
	*/
}

func (d RunData) runArray(arr vcfg.RelIntervalArray, index *int, once bool, chDone chan int) {

	if 0 == len(arr) {
		chDone <- vomni.DoneStop
		return
	}

	*index = nextIndex(*index, len(arr))

	for {
		// set the interval for this new state
		tick := time.NewTicker(arr[*index].Seconds)
		// put the message in the send queue
		msg := vmsg.QeueuGpioSet(d.Point, d.UDPAddr, arr[*index].Gpio, arr[*index].State)

		fmt.Println("vk-xxx -------> POINT", d.Point, "Karolina", msg, "ADDR", d.UDPAddr)

		d.LogStr(vomni.LogFileCdInfo, fmt.Sprintf("Send message: %q", msg))

		done := 0

		select {
		case done = <-d.ChDone:

			fmt.Println("###\n###\n###\n", d.Point, "*** THREE MAIN SECTIONS\n###\n###\n###")

		case <-tick.C:
			*index = nextIndex(*index, len(arr))

			if once && 0 == *index {
				done = vomni.DoneStop
			}
		}

		if 0 < done {
			*index = vomni.PointNonActiveIndex
			chDone <- done
			return
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

func (d RunData) StartRotate() (err error) {

	if err = d.prepareRotateLoggers(); nil != err {
		return vutils.ErrFuncLine(fmt.Errorf("Couldn't prepare the point %q rotate configuration - %v", d.Point, err))
	}

	return vrotate.StartPointLoggers(d.Point, d.Logs)
}

func (d RunData) prepareRotateLoggers() (err error) {
	for k, v := range d.Logs {
		// Let's open the log data fiel
		d.Logs[k].LogFilePtr, err = vutils.OpenFile(v.LogFile, vomni.LogFileFlags, vomni.LogUserPerms)
		if nil != err {
			return vutils.ErrFuncLine(fmt.Errorf("Could not open the point %q data log file --- %v", d.Point, err))
		}
		// prepare Logger fields
		for k1, v1 := range v.Loggers {
			log := vomni.PointLogger{LogPrefix: v1.LogPrefix, Logger: vutils.LogNew(d.Logs[k].LogFilePtr, v1.LogPrefix)}
			d.Logs[k].Loggers[k1] = log
		}
	}

	return
}
