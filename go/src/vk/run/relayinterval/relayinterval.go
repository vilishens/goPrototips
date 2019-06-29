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

var RunningPoints map[string]*RunData

func init() {
	RunningPoints = make(map[string]*RunData)
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

func (d RunData) LetsGo(chGoOn chan bool, chDone chan int, chErr chan error) {

	//d.UDPAddr = addr vk-xxx

	fmt.Println("$$$$$$$$$$$$$$$$ FINAL $$$$$$$$$$$$$$$$$", d.UDPAddr)

	fmt.Printf("============ UDPAddr %+v\n", d.UDPAddr)

	d.Index = AllIndex{Start: vomni.PointNonActiveIndex, Base: vomni.PointNonActiveIndex, Finish: vomni.PointNonActiveIndex}

	locGoOn := make(chan bool)
	locDone := make(chan int)
	locErr := make(chan error)
	go d.run(locGoOn, locDone, locErr)

	<-locGoOn

	d.SetState(vomni.PointStateActive|vomni.PointStateSigned, true)

	chGoOn <- true
}

func (d RunData) GetDone(done int) {
	d.ChDone <- done
}

func (d RunData) Ready() (ready bool) {

	ready = true

	/*
		if !ready {
				d.Point,
				vomni.PointCfgData[d.Type].CfgStr)

			d.LogStr(vomni.LogFileCdErr, str)
		} else {
			d.SetState(vomni.PointCfgStateReady, true)

			str := fmt.Sprintf("Point %q - %q configuration ready",
				d.Point,
				vomni.PointCfgData[d.Type].CfgStr)

			d.LogStr(vomni.LogFileCdInfo, str)
		}
	*/
	return
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

	allStages := []stage{
		stage{once: true, index: &d.Index.Start, cfg: d.Cfg.Start},   // start sequence
		stage{once: false, index: &d.Index.Base, cfg: d.Cfg.Base},    // base sequence
		stage{once: true, index: &d.Index.Finish, cfg: d.Cfg.Finish}} // finishe sequence

	for _, v := range allStages {
		go d.runArray(v.cfg, v.index, v.once, locDone)
		rc := <-locDone
		if vomni.DoneDisconnected == rc {
			d.SetState(vomni.DoneDisconnected, true)
			str := fmt.Sprintf("Point %q lost connection", d.Point)
			d.LogStr(vomni.LogFileCdErr, str)

			fmt.Printf("***\n***\n*** Nutivara %q \n***\n***\n", d.Point)

			break
		}
	}
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

		fmt.Printf("vk-xxx SHADOW *** -------> POINT %15s ADDR %20s MSG %s\n", d.Point, d.UDPAddr.IP.String(), msg)

		d.LogStr(vomni.LogFileCdInfo, fmt.Sprintf("Send message: %q", msg))

		done := 0

		select {

		case msg := <-d.ChMsg:
			fmt.Printf("vk-xxx ###\n###\n###\n Point %q received a message %q *** HEAVY METAL\n###\n###\n###\n", d.Point, msg)

		case done = <-d.ChDone:

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
	//	chDone <- vomni.DoneStop

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

func (d *RunData) SetUDPAddr(addr net.UDPAddr) {
	/*
		fAddr := reflect.ValueOf(&d.UDPAddr)

		elemAddr := fAddr.Elem()
		if elemAddr.Kind() == reflect.Struct {
			//		fmt.Println("ADDRESE ir struktūra")
			fIP := elemAddr.FieldByName("IP")
			if fIP.IsValid() && fIP.CanSet() && fIP.Kind() == reflect.Slice {
				fIP.SetBytes(addr.IP)
			}

			fPort := elemAddr.FieldByName("Port")
			if fPort.IsValid() && fPort.CanSet() && fPort.Kind() == reflect.Int {
				fPort.SetInt(int64(addr.Port))
			}
		}
	*/
	d.UDPAddr = addr
}

func (d RunData) GetUDPAddr() (addr net.UDPAddr) {
	return d.UDPAddr
}

func (d *RunData) SetState(state int, on bool) {

	if on {
		d.State |= state
	} else {
		d.State &^= state
	}
}

func (d *RunData) GetState() (state int) {
	return d.State
}

func (d *RunData) setState(state int, on bool) {
	if on {
		d.State |= state
	} else {
		d.State &^= state
	}
}
