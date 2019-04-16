package pointconfig

import (
	"fmt"
	vutils "vk/utils"
)

var PointsAllJSON CfgJSONData
var PointsAllData AllPointCfgData

func init() {

	PointsAllData = make(map[string]PointCfgData)
	PointsAllJSON = CfgJSONData{}
}

func GetPointCfg(chGoOn chan bool, chDone chan int, chErr chan error) {
	locDone := make(chan bool)
	locErr := make(chan error)

	go preparePointCfg(locDone, locErr)

	select {
	case err := <-locErr:
		vutils.LogErr(err)
		chErr <- vutils.ErrFuncLine(err)
	case <-locDone:
		chGoOn <- true
	}

	fmt.Println(PointsAllData.String())
}

func preparePointCfg(doneCh chan bool, errCh chan error) {

	var err error

	if PointsAllJSON, err = loadPointCfg(); nil == err {
		err = PointsAllJSON.putCfg4Run()
	}

	if nil != err {
		errCh <- err
		return
	}

	doneCh <- true
}

func (d CfgJSONData) putCfg4Run() (err error) {
	// prepare Relay On/Off Interval configuration
	if err = d.RelIntervalJSON.putCfg4Run(); nil != err {
		err = vutils.ErrFuncLine(fmt.Errorf("Relay Interval configuration Error - %s", err.Error()))
		return
	}

	return
}
