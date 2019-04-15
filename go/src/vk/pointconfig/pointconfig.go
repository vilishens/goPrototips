package pointconfig

import (
	"fmt"
	vparams "vk/params"
	vutils "vk/utils"
)

var PointsCfgJSON CfgJSONData
var PointsCfgData AllPointCfgData

func init() {

	PointsCfgData = make(map[string]PointCfgData)
	PointsCfgJSON = CfgJSONData{}
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
}

func preparePointCfg(doneCh chan bool, errCh chan error) {

	var err error

	if PointsCfgJSON, err = loadPointCfg(); nil == err {
		err = PointsCfgJSON.putCfg4Run()
	}

	if nil != err {
		errCh <- err
		return
	}

	doneCh <- true
}

func (d CfgJSONData) putCfg4Run() (err error) {
	if err = d.RelOnOffIntervalJSON.putCfg4Run(); nil != err {
		return
	}

	return
}

func (d CfgRelOnOffIntervalPoints) putCfg4Run() (err error) {

	for k, _ := range d {

		fmt.Println("RelOnOffInterval", k)
	}

	return fmt.Errorf("Bambarbiya!")
}

func loadPointCfg() (data CfgJSONData, err error) {

	if has, _ := vutils.PathExists(vparams.Params.PointConfigFile); !has {
		if err := vutils.FileCopy(vparams.Params.PointConfigOriginalFile, vparams.Params.PointConfigFile); nil != err {
			return CfgJSONData{}, vutils.ErrFuncLine(err)
		}
	}

	if err = vutils.ReadJson(vparams.Params.PointConfigFile, &data); nil != err {
		return CfgJSONData{}, vutils.ErrFuncLine(err)
	}

	return
}
