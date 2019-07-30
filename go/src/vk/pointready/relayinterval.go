package pointready

import (
	"fmt"
	vomni "vk/omnibus"
	vparams "vk/params"
	vpointconfig "vk/pointconfig"
	vrotate "vk/rotate"
	vrunrelayinterval "vk/run/relayinterval"
	vutils "vk/utils"
)

func relayInterval() {

	for k, v := range vpointconfig.AllPointData.Default {

		//	for k, v := range vpointconfig.AllPointData.Running {
		if 0 < (v.List & vomni.CfgTypeRelayInterval) {
			d := NewRunInterface(k, v)
			/*
				vrunrelayinterval.RunInterface{}
				d.Point = k
				d.State = vomni.PointCfgStateUnknown
				d.Type = vomni.CfgTypeRelayInterval

				fmt.Println("Nepiemirsti, ka vajag FACTORY conf!!!")

				d.CfgDefault = v.Cfg.RelInterv
				d.CfgRun = v.Cfg.RelInterv
				d.CfgSaved = v.Cfg.RelInterv

				d.Index = vrunrelayinterval.AllIndex{Start: vomni.PointNonActiveIndex, Base: vomni.PointNonActiveIndex, Finish: vomni.PointNonActiveIndex}

				d.ChDone = make(chan int)
				d.ChErr = make(chan error)
				d.ChMsg = make(chan string)
			*/
			vrunrelayinterval.RunningPoints[k] = &d

			dd := NewRunData(k, v) //vrunrelayinterval.RunInfo(d)
			vrunrelayinterval.RunningData[k] = &dd

			//logs, err := pointLoggers(d.Point, d.Type)
			// handle all loggers of the point
			logs, err := relayIntervalPointLoggers(d.Point, d.Type)
			if nil != err {
				vomni.RootErr <- err
				return
			}

			d.Logs = logs
		}
	}
}

func relayIntervalPointLoggers(point string, cd int) (logs []vomni.PointLog, err error) {

	key := vomni.LogFileCdErr | vomni.LogFileCdInfo

	// find the path of the data log file
	logF := vrotate.RotatePointFilePath(key, vparams.Params.LogPointPath, point, cd)

	// the rotate configuration template
	tmplF := vutils.FileAbsPath(vparams.Params.RotatePointInfoTmpl, "")

	// loggers into the data log file
	loggers := vrotate.RotatePointLoggers(key)

	logs = append(logs, vomni.PointLog{LogFile: logF, LogTmpl: tmplF, Loggers: loggers})

	fmt.Printf("%q ***** Folderītis  %q\n", point, tmplF)
	fmt.Printf("%q ***** File        %q\n", point, logF)
	fmt.Printf("%q ***** Loggers     %+v\n", point, loggers)

	return
}

func NewRunInterface(point string, cfg vpointconfig.PointCfgData) (d vrunrelayinterval.RunInterface) {
	//d := vrunrelayinterval.RunInterface{}
	d.Point = point
	d.State = vomni.PointCfgStateUnknown
	d.Type = vomni.CfgTypeRelayInterval

	fmt.Println("Nepiemirsti, ka vajag FACTORY conf!!!")

	d.CfgDefault = cfg.Cfg.RelInterv
	d.CfgRun = cfg.Cfg.RelInterv
	d.CfgSaved = cfg.Cfg.RelInterv

	//	d.Index = vrunrelayinterval.AllIndex{Start: vomni.PointNonActiveIndex, Base: vomni.PointNonActiveIndex, Finish: vomni.PointNonActiveIndex}

	d.ChDone = make(chan int)
	d.ChErr = make(chan error)
	d.ChMsg = make(chan string)

	return d
}

func NewRunData(point string, cfg vpointconfig.PointCfgData) (d vrunrelayinterval.RunData) {
	//d := vrunrelayinterval.RunInterface{}
	d.Point = point
	d.State = vomni.PointCfgStateUnknown
	d.Type = vomni.CfgTypeRelayInterval

	fmt.Println("Nepiemirsti, ka vajag FACTORY conf!!!")

	d.CfgDefault = cfg.Cfg.RelInterv
	d.CfgRun = cfg.Cfg.RelInterv
	d.CfgSaved = cfg.Cfg.RelInterv

	d.Index = vrunrelayinterval.AllIndex{Start: vomni.PointNonActiveIndex, Base: vomni.PointNonActiveIndex, Finish: vomni.PointNonActiveIndex}

	d.ChDone = make(chan int)
	d.ChErr = make(chan error)
	d.ChMsg = make(chan string)

	return d
}
