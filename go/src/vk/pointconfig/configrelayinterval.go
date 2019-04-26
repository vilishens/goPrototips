package pointconfig

import (
	"strconv"
	vomni "vk/omnibus"
	vparams "vk/params"
	vutils "vk/utils"
)

func (d CfgRelIntervalPoints) putCfg4Run() (err error) {

	for k, v := range d {

		var newD RelIntervalStruct

		if newD, err = v.putCfg4Run(); nil != err {
			err = vutils.ErrFuncLine(err)
			return
		}

		if _, has := PointsAllData[k]; !has {
			PointsAllData[k] = PointCfgData{}
		}

		tmpD := PointsAllData[k]
		tmpD.List |= vomni.CfgTypeRelayInterval
		tmpD.Cfg.RelInterv = newD
		tmpD.CfgSaved.RelInterv = newD
		PointsAllData[k] = tmpD
	}

	return
}

func (d CfgRelIntervalStruct) putCfg4Run() (newD RelIntervalStruct, err error) {

	newD = RelIntervalStruct{}
	if newD.Start, err = d.Start.putCfg4Run(); nil != err {
		err = vutils.ErrFuncLine(err)
		return
	}

	return
}

func (d CfgRelIntervalArray) putCfg4Run() (newD []RelInterval, err error) {

	newD = []RelInterval{}

	for _, v := range d {
		tmpD := RelInterval{Gpio: -1, State: -1, Seconds: 0}

		if "" != v.Gpio {
			if tmpD.Gpio, err = strconv.Atoi(v.Gpio); nil != err {
				err = vutils.ErrFuncLine(err)
				return
			}
		}

		if "" != v.State {
			if tmpD.State, err = strconv.Atoi(v.State); nil != err {
				err = vutils.ErrFuncLine(err)
				return
			}
		}

		if "" != v.Interval {
			if tmpD.Seconds, err = vutils.ConfInterval2Seconds(v.Interval); nil != err {
				err = vutils.ErrFuncLine(err)
				return
			}
		}

		newD = append(newD, tmpD)
	}

	return
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
