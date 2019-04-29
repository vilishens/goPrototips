package pointconfig

import (
	"strconv"
	vomni "vk/omnibus"
	vparams "vk/params"
	vutils "vk/utils"
)

func (d CfgRelIntervalStruct) putCfg4Run(point string) (err error) {

	newD := RelIntervalStruct{}
	if newD.Start, err = d.Start.putCfg4Run(); nil != err {
		err = vutils.ErrFuncLine(err)
		return
	}

	if newD.Base, err = d.Base.putCfg4Run(); nil != err {
		err = vutils.ErrFuncLine(err)
		return
	}

	if newD.Finish, err = d.Finish.putCfg4Run(); nil != err {
		err = vutils.ErrFuncLine(err)
		return
	}

	tmpD := PointsAllData[point]
	tmpD.List |= vomni.CfgTypeRelayInterval
	tmpD.Cfg.RelInterv = newD
	tmpD.CfgSaved.RelInterv = newD
	PointsAllData[point] = tmpD

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
