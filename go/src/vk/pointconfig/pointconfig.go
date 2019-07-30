package pointconfig

import (
	"fmt"
	vutils "vk/utils"
)

var PointsAllJSON CfgJSONData
var PointsAllDefaultJSON CfgJSONData
var PointsAllData AllPointCfgData
var PointsAllDefaultData AllPointCfgData

func init() {
	PointsAllData = make(map[string]PointCfgData)
	PointsAllDefaultData = make(map[string]PointCfgData)
	PointsAllJSON = CfgJSONData{}
	PointsAllDefaultJSON = CfgJSONData{}
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

	if PointsAllDefaultJSON, err = loadPointDefaultCfg(); nil == err {
		err = PointsAllDefaultJSON.putCfgDefault4Run()
	}

	if nil != err {
		errCh <- err
		return
	}

	doneCh <- true
}

func (d CfgJSONData) putCfg4Run() (err error) {

	for k, v := range d {
		if _, has := PointsAllData[k]; !has {
			PointsAllData[k] = PointCfgData{}
		}

		var newStruct PointCfgData
		if v.RelIntervalJSON.hasCfgRelInterval() {
			if newStruct, err = v.RelIntervalJSON.putCfg4Run(PointsAllData[k]); nil != err {
				err = vutils.ErrFuncLine(fmt.Errorf("Relay Interval configuration Error - %s", err.Error()))
				return
			} else {
				PointsAllData[k] = newStruct
			}
		}

		//		if v.RelIntervalJSON.hasCfgRelInterval() {
		//			if err = v.RelIntervalJSON.putCfg4Run(k); nil != err {
		//				err = vutils.ErrFuncLine(fmt.Errorf("Relay Interval configuration Error - %s", err.Error()))
		//				return
		//			}
		//		}
	}

	return
}

func (d CfgJSONData) putCfgDefault4Run() (err error) {

	for k, v := range d {
		if _, has := PointsAllDefaultData[k]; !has {
			PointsAllDefaultData[k] = PointCfgData{}
		}

		var newStruct PointCfgData
		if v.RelIntervalJSON.hasCfgRelInterval() {
			if newStruct, err = v.RelIntervalJSON.putCfgDefault4Run(PointsAllDefaultData[k]); nil != err {
				err = vutils.ErrFuncLine(fmt.Errorf("Relay Interval Default configuration Error - %s", err.Error()))
				return
			} else {
				PointsAllDefaultData[k] = newStruct
			}
		}
	}

	return
}

func (d CfgRelIntervalStruct) hasCfgRelInterval() (has bool) {
	if 0 < len(d.Start) {
		return true
	}

	if 0 < len(d.Base) {
		return true
	}

	if 0 < len(d.Finish) {
		return true
	}

	return false
}
