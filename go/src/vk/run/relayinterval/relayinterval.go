package runrelayinterval

import (
	"fmt"
	"log"
	"net"
	"strings"
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

	dx := RunningData[d.Point]
	return dx.CfgDefault, dx.CfgRun, dx.CfgSaved, dx.Index, dx.State
}

func (d RunInterface) ReceiveWeb(cmd int, data interface{}) {

	newCmd := 0

	switch cmd {
	case vomni.PointCmdLoadCfgIntoPoint:
		RunningData[d.Point].CfgRun = webInterface2Struct(data)
		newCmd = cmdRestart
	case vomni.PointCmdSaveCfg:

		if err := webSavePointCfg(d.Point, data); nil != err {
			vomni.LogErr.Println(vutils.ErrFuncLine(err))
		} else {
			dNew := webInterface2Struct(data)
			RunningData[d.Point].CfgRun = dNew
			RunningData[d.Point].CfgSaved = dNew
		}
		newCmd = cmdRestart
	default:
		log.Fatal("RelayInterval received ", cmd, ". What to do?")
	}

	RunningData[d.Point].ChCmd <- newCmd

	//data.(vcfg.RelIntervalStruct)

	/*
		RunningData[d.Point]

		cfg := 147

		switch cfg {
		default:
			str := fmt.Sprintf("\n\nDon't know how te receive configuration %08X for %q\n\n", cfg, d.Point)
			panic(str)
		}
	*/
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

	stop := false
	for !stop {

		fmt.Println("wwwww\nwwwww\nZIRGS-ZIRGS-ZIRGS\nwwwww\nwwwww")

		go d.run(locGoOn, locDone, locErr)

		waitNext := true
		for waitNext {
			select {
			case <-locGoOn:
				d.SetState(vomni.PointStateActive|vomni.PointStateSigned, true)
				chGoOn <- true
			case rc := <-locDone:

				fmt.Println("=====\n=====\nGruzovik", rc, " && ", vomni.PointCmdLoadCfgIntoPoint, "\n=====\n=====")

				if vomni.PointCmdLoadCfgIntoPoint == rc {

					d.SetState(vomni.PointStateSigned, false)
					waitNext = false
				}
			}
		}

		fmt.Println("xxxxx\nxxxxx\nNehuj!!!\n=====\n=====")

	}

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

	chGoOn <- true

	locDone := make(chan int)

	stop := false
	for !stop {

		fmt.Println("$$$\n$$$\nVICINS-GEIRGS\n$$$\n$$$")

		RunningData[d.Point].Index = AllIndex{Start: vomni.PointNonActiveIndex,
			Base: vomni.PointNonActiveIndex, Finish: vomni.PointNonActiveIndex}

		allStages := []stage{
			stage{once: true, runEmptyArr: false, index: &RunningData[d.Point].Index.Start, cfg: RunningData[d.Point].CfgRun.Start},   // start sequence
			stage{once: false, runEmptyArr: true, index: &RunningData[d.Point].Index.Base, cfg: RunningData[d.Point].CfgRun.Base},     // base sequence
			stage{once: true, runEmptyArr: false, index: &RunningData[d.Point].Index.Finish, cfg: RunningData[d.Point].CfgRun.Finish}} // finishe sequence

		for _, v := range allStages {
			go d.runArray(v, locDone)
			rc := <-locDone

			if vomni.DoneDisconnected == rc {
				d.SetState(vomni.DoneDisconnected, true)
				str := fmt.Sprintf("Point %q lost connection", d.Point)
				d.LogStr(vomni.LogFileCdErr, str)

				fmt.Printf("***\n***\n*** Nutivara %q \n***\n***\n", d.Point)

				stop = true
				break
			}

			/*
				if vomni.PointCmdLoadCfgIntoPoint == rc {
					RunningData[d.Point].Index = AllIndex{Start: vomni.PointNonActiveIndex,
						Base: vomni.PointNonActiveIndex, Finish: vomni.PointNonActiveIndex}

					fmt.Println("$$$$$\n$$$$$\nIgor Botvin\n@@@@\n@@@@@@")

					//chDone <- rc

					//				return
					//stop = true
					break
				}
			*/
			if cmdRestart == rc {
				break
			}

		}
	}
}

func (d RunInterface) runArray(st stage, chDone chan int) {

	if !st.runEmptyArr && 0 == len(st.cfg) {
		chDone <- vomni.DoneStop
		return
	}

	*st.index = nextIndex(*st.index, len(st.cfg))

	for {

		var tick *time.Ticker

		if !(st.runEmptyArr && 0 == len(st.cfg)) {

			// set the interval for this new state
			tick = time.NewTicker(st.cfg[*st.index].Seconds)
			// put the message in the send queue
			msg := vmsg.QeueuGpioSet(d.Point, d.UDPAddr, st.cfg[*st.index].Gpio, st.cfg[*st.index].State)

			fmt.Printf("vk-xxx SHADOW *** -------> POINT %15s ADDR %20s MSG %s\n", d.Point, d.UDPAddr.IP.String(), msg)

			d.LogStr(vomni.LogFileCdInfo, fmt.Sprintf("Send message: %q", msg))
		}

		done := 0

		select {

		case cmd := <-RunningData[d.Point].ChCmd:

			// Seit jāieliek msg apstrāde

			chDone <- cmd
			return

		case done = <-d.ChDone:

		case <-tick.C:
			*st.index = nextIndex(*st.index, len(st.cfg))

			if st.once && 0 == *st.index {
				done = vomni.DoneStop
			}
		}

		if 0 < done {
			*st.index = vomni.PointNonActiveIndex

			chDone <- done
			return
		}

	}
	//	chDone <- vomni.DoneStop

}

/*
func (d RunInterface) runArray(arr vcfg.RelIntervalArray, index *int, once bool, runEmpty bool, chDone chan int) {

	if !runEmpty && 0 == len(arr) {
		chDone <- vomni.DoneStop
		return
	}

	*index = nextIndex(*index, len(arr))

	for {

		var tick *time.Ticker

		if !(runEmpty && 0 == len(arr)) {

			// set the interval for this new state
			tick = time.NewTicker(arr[*index].Seconds)
			// put the message in the send queue
			msg := vmsg.QeueuGpioSet(d.Point, d.UDPAddr, arr[*index].Gpio, arr[*index].State)

			fmt.Printf("vk-xxx SHADOW *** -------> POINT %15s ADDR %20s MSG %s\n", d.Point, d.UDPAddr.IP.String(), msg)

			d.LogStr(vomni.LogFileCdInfo, fmt.Sprintf("Send message: %q", msg))
		}

		done := 0

		select {

		case cmd := <-RunningData[d.Point].ChCmd:

			// Seit jāieliek msg apstrāde

			chDone <- cmd
			return

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
*/
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

func webSavePointCfg(point string, data interface{}) (err error) {

	var newData vcfg.CfgRelIntervalStruct

	if newData, err = webInterface2SaveCfg(data); nil != err {
		err = vutils.ErrFuncLine(err)
		return
	}

	fmt.Println("NEW DATA\n", newData)

	whole := vcfg.AllPointData.RunningJSON //CfgJSONData

	fmt.Println("AKTIVS\n", whole, "\nDEFUALTE\n", vcfg.PointsAllDefaultData)

	pData := whole[point]

	pData.RelIntervalJSON = newData

	whole[point] = pData

	fmt.Println("PEC\n", whole)

	return whole.Save()
}

func webInterface2SaveCfg(inter interface{}) (web vcfg.CfgRelIntervalStruct, err error) {
	// WEB struct
	web = vcfg.CfgRelIntervalStruct{}

	for part, v := range inter.(map[string]interface{}) { // list add configuration parts
		d := vcfg.CfgRelIntervalArray{}        // array for the configuration part records
		for _, v1 := range v.([]interface{}) { // fill part record array
			rec := vcfg.CfgRelInterval{} // storage for a record data

			for k2, v2 := range v1.(map[string]interface{}) {
				switch strings.ToUpper(k2) {
				case "GPIO":
					rec.Gpio = v2.(string)
				case "STATE":
					rec.State = v2.(string)
				case "SECONDS":
					if rec.Interval, err = vutils.DurationStrToIntervalStr(v2.(string)); nil != err {
						err = vutils.ErrFuncLine(err)
						return
					}
				default:
					log.Fatal(fmt.Sprintf("Unknow WEB interface record field \"%s\"", k2))
				}
			}

			d = append(d, rec)
		}

		switch strings.ToUpper(part) {
		case "START":
			web.Start = d
		case "BASE":
			web.Base = d
		case "FINISH":
			web.Finish = d
		default:
			log.Fatal(fmt.Sprintf("Unknow WEB interface part \"%s\"", part))
		}
	}

	return
}
