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

var RunningPoints map[string]*RunInterface
var RunningData map[string]*RunData

func init() {
	RunningPoints = make(map[string]*RunInterface)
	RunningData = make(map[string]*RunData)
}

func (d RunInterface) GetCfgs() (cfgDefault interface{}, cfgRun interface{}, cfgSaved interface{},
	cfgIndex interface{}, cfgState interface{}) {

	//	back, err := json.Marshal(d.CfgDefault)
	//	if nil != err {
	//		panic(err)
	//		return
	//	}

	//	json.Unmarshal(back, &cfgDefault)

	return d.CfgDefault, d.CfgRun, d.CfgSaved, RunningData[d.Point].Index, d.State
}

func (d RunInterface) LogStr(infoCd int, str string) {

	for _, v := range d.Logs {
		for k1, v1 := range v.Loggers {
			if k1 == infoCd {
				vutils.LogStr(v1.Logger, str)
			}
		}
	}
}

func (d RunInterface) LetsGo(chGoOn chan bool, chDone chan int, chErr chan error) {

	//d.UDPAddr = addr vk-xxx

	fmt.Println("$$$$$$$$$$$$$$$$ FINAL $$$$$$$$$$$$$$$$$", d.UDPAddr)

	fmt.Printf("============ UDPAddr %+v\n", d.UDPAddr)

	//d.Index = AllIndex{Start: vomni.PointNonActiveIndex, Base: vomni.PointNonActiveIndex, Finish: vomni.PointNonActiveIndex}

	RunningData[d.Point].Index = AllIndex{Start: vomni.PointNonActiveIndex, Base: vomni.PointNonActiveIndex, Finish: vomni.PointNonActiveIndex}

	locGoOn := make(chan bool)
	locDone := make(chan int)
	locErr := make(chan error)
	go d.run(locGoOn, locDone, locErr)

	<-locGoOn

	d.SetState(vomni.PointStateActive|vomni.PointStateSigned, true)

	chGoOn <- true
}

func (d RunInterface) GetDone(done int) {
	d.ChDone <- done
}

func (d RunInterface) Ready() (ready bool) {

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

func (d RunInterface) run(chGoOn chan bool, chDone chan int, chErr chan error) {

	fmt.Printf("Point %q Addr %+v Index %+v\n", d.Point, d.UDPAddr, RunningData[d.Point].Index)

	chGoOn <- true

	locDone := make(chan int)
	type stage struct {
		once  bool
		index *int
		cfg   vcfg.RelIntervalArray
	}

	allStages := []stage{
		stage{once: true, index: &RunningData[d.Point].Index.Start, cfg: d.CfgRun.Start},   // start sequence
		stage{once: false, index: &RunningData[d.Point].Index.Base, cfg: d.CfgRun.Base},    // base sequence
		stage{once: true, index: &RunningData[d.Point].Index.Finish, cfg: d.CfgRun.Finish}} // finishe sequence

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

func (d RunInterface) runArray(arr vcfg.RelIntervalArray, index *int, once bool, chDone chan int) {

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

func (d RunInterface) StartRotate() (err error) {

	if err = d.prepareRotateLoggers(); nil != err {
		return vutils.ErrFuncLine(fmt.Errorf("Couldn't prepare the point %q rotate configuration - %v", d.Point, err))
	}

	return vrotate.StartPointLoggers(d.Point, d.Logs)
}

func (d RunInterface) prepareRotateLoggers() (err error) {
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

func (d *RunInterface) SetUDPAddr(addr net.UDPAddr) {
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

func (d RunInterface) GetUDPAddr() (addr net.UDPAddr) {
	return d.UDPAddr
}

func (d *RunInterface) SetState(state int, on bool) {

	if on {
		d.State |= state
	} else {
		d.State &^= state
	}
}

func (d *RunInterface) GetState() (state int) {
	return d.State
}

func (d *RunInterface) setState(state int, on bool) {
	if on {
		d.State |= state
	} else {
		d.State &^= state
	}
}
