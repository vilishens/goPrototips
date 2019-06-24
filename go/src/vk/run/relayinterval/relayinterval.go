package runrelayinterval

import (
	"fmt"
	"net"
	"reflect"
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

func (d RunData) LetsGo(addr net.UDPAddr, chGoOn chan bool, chDone chan int, chErr chan error) {

	//d.UDPAddr = addr vk-xxx

	fmt.Println("$$$$$$$$$$$$$$$$ FINAL $$$$$$$$$$$$$$$$$", d.UDPAddr)

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

	ip := d.UDPAddr.IP.String()

	if d.UDPAddr.IP == nil {
		fmt.Println("vk-xxx AKEX->SITKOVE IP nil ################################################")
	}

	fmt.Printf("vk-xxx =======> POINT %s IP %s 0x%04x\n", d.Point, ip, done)
	fmt.Printf("vk-xxx =======> POINT %s IP %s 0x%04x\n", d.Point, "kira", done)
	fmt.Printf("vk-xxx =======> POINT %s IP %+v 0x%04x\n", d.Point, d.UDPAddr, done)

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

	allStages := []stage{
		stage{once: true, index: &d.Index.Start, cfg: d.Cfg.Start},   // start sequence
		stage{once: false, index: &d.Index.Base, cfg: d.Cfg.Base},    // base sequence
		stage{once: true, index: &d.Index.Finish, cfg: d.Cfg.Finish}} // finishe sequence

	for _, v := range allStages {
		go d.runArray(v.cfg, v.index, v.once, locDone)
		rc := <-locDone
		if vomni.DoneDisconnected == rc {

			fmt.Printf("vk-xxx ,,,,,,,,,,,,,,,,,,,,,Susuman %+v IZEJU\n", d.UDPAddr)
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

		fmt.Printf("vk-xxx -------> POINT %15s ADDR %20s MSG %s\n", d.Point, d.UDPAddr.IP.String(), msg)

		d.LogStr(vomni.LogFileCdInfo, fmt.Sprintf("Send message: %q", msg))

		done := 0

		select {
		case done = <-d.ChDone:

			fmt.Println("vk-xxx ###\n###\n###\n", d.Point, "*** ŅUTA FEDERMESSER saņēmu\n###\n###\n###")

		case <-tick.C:
			*index = nextIndex(*index, len(arr))

			if once && 0 == *index {
				done = vomni.DoneStop
			}
		}

		if 0 < done {
			fmt.Printf("vk-xxx ***\n***\n***\n %s *** tamara soboļ beidzu %x \n***\n***\n***\n", d.Point, done)

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

func (d *RunData) SetUDPAddr(addr net.UDPAddr) {

	fmt.Println("vk-xxx ADDR ", addr)

	adrese := reflect.ValueOf(&d.UDPAddr)

	s := reflect.ValueOf(&d.UDPAddr).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		fmt.Printf("gorgonzola %d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}

	fmt.Println("ORDENA BY DAL Type:", reflect.TypeOf(d.UDPAddr))
	fmt.Printf("Ferratum Type: %T\n", adrese)

	sAddr := adrese.Elem()
	if sAddr.Kind() == reflect.Struct {
		fmt.Println("ADDRESE ir struktūra")
		fIP := sAddr.FieldByName("IP")
		fPort := sAddr.FieldByName("Port")
		if fIP.IsValid() && fIP.CanSet() && fIP.Kind() == reflect.Slice {

			fIP.SetBytes(addr.IP)

			fmt.Println("ADDRESE ir valida IP")

		}

		if fPort.IsValid() && fPort.CanSet() && fPort.Kind() == reflect.Int {

			fPort.SetInt(int64(addr.Port))

			fmt.Println("ADDRESE ir valida Port")

		}
	}

	ss := reflect.ValueOf(&d.UDPAddr).Elem()
	typeOfT = ss.Type()
	for i := 0; i < ss.NumField(); i++ {
		f := ss.Field(i)
		fmt.Printf("orizava %d: %s %s = %v\n", i,
			typeOfT.Field(i).Name, f.Type(), f.Interface())
	}

	//	s := a.MapKeys()
	/*
	   	s := ps.Elem()
	       if s.Kind() == reflect.Struct {
	           // exported field
	           f := s.FieldByName("N")
	           if f.IsValid() {
	               // A Value can be changed only if it is
	               // addressable and was not obtained by
	               // the use of unexported struct fields.
	               if f.CanSet() {
	                   // change value of N
	                   if f.Kind() == reflect.Int {
	                       x := int64(7)
	                       if !f.OverflowInt(x) {
	                           f.SetInt(x)
	                       }
	                   }
	               }
	           }
	       }
	*/

	//	fmt.Printf("vk-xxx Type %T %T\n %+v\nIP %v\n", ip, port, s, ps)

	//	a. UDPAddr = addr
	//	fmt.Printf("vk-xxx TIPS %T\nSTRUCTURA %+v\nADRESE %v", a, a, a.UDPAddr)

	//	a := d

	//_ = s

	//a.UDPAddr = addr

	//d = a

	//d.UDPAddr = addr

	fmt.Println("vk-xxx SLOMANNIJ KLINOK ", d.UDPAddr)
}

func (d RunData) GetUDPAddr() (addr net.UDPAddr) {

	fmt.Printf("vk-xxx KVASUBA adrese:\n\tIP   %v\n\tPORT %d\n", d.UDPAddr.IP, d.UDPAddr.Port)

	return d.UDPAddr
}
