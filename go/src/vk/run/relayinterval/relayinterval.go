package runrelayinterval

import (
	"fmt"
	"net"
	"time"
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

	d.Index = AllIndex{Start: -1, Base: -1, Finish: -1}

	locGoOn := make(chan bool)
	locDone := make(chan int)
	locErr := make(chan error)
	go d.run(locGoOn, locDone, locErr)

	<-locGoOn

	d.State |= vomni.PointStateActive
	d.State |= vomni.PointStateSigned

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

		//						t := time.Now()

		//vk-xxx
		type dst struct {
			name string
			host net.UDPAddr
		}

		pref := dst{name: "BĻITVINGS", host: net.UDPAddr{IP: []byte{192, 168, 7, 15}, Port: 45678}}

		txt := "ZIRGS!!!"

		msg := fmt.Sprintf("DST: %+v, MSG: %q", pref, txt)

		// vk-xxx

		//		dst :=
		//		fmt.Println(d.Point, "@@@@@@@@@@@@@@@@", t.Format(vomni.TimeFormat1), "*************** INDEX ", *index, "JĀSŪTA CMD PIRMS INTERVALA", arr[*index].Seconds.Seconds())

		d.LogStr(vomni.LogFileCdInfo, msg)

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
