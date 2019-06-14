package pointready

import (
	"fmt"
	"path/filepath"
	vomni "vk/omnibus"
	vparams "vk/params"
	vpointconfig "vk/pointconfig"
	vrunrelayinterval "vk/run/relayinterval"
	vutils "vk/utils"
)

func relayInterval() {

	for k, v := range vpointconfig.PointsAllData {
		if 0 < (v.List & vomni.CfgTypeRelayInterval) {
			d := vrunrelayinterval.RunData{}
			d.Point = k
			d.State = vomni.PointStateUnknown
			d.Type = vomni.CfgTypeRelayInterval
			d.Cfg = v.Cfg.RelInterv
			d.CfgSaved = v.Cfg.RelInterv

			//logs, err := pointLoggers(d.Point, d.Type)
			// handle all loggers of the point
			logs, err := relayIntervalPointLoggers(d.Point, d.Type)

			if nil != err {
				vomni.RootErr <- err
				return
			}

			_ = logs
			d.Logs = make(map[int]vomni.PointLog)
			//d.Logs[d.Type] = logs

			vrunrelayinterval.RunningPoints[k] = d
		}
	}
}

func relayIntervalPointLoggers(point string, cd int) (logs map[int]vomni.PointLog, err error) {

	logKey := vomni.LogFileErr | vomni.LogFileInfo

	i := 0
	j := 1
	ending := ""
	for i = 0; i < logKey; i++ {

		j <<= 1

		fmt.Printf("KEY %2d %2d I %2d J%2d\n", logKey, logKey&j, i, j)

		if 0 == logKey&j {
			continue
		}

		if "" != ending {
			ending += "-"
		}
		ending += vomni.PointLogData[j].FileEnd
	}

	fmt.Printf("ENDING  %q\n", ending)

	// rotate configuration template
	tmplF := vutils.FileAbsPath(vparams.Params.RotatePointInfoTmpl, "")
	// rotate log file
	logF := vutils.FileAbsPath(filepath.Join(vparams.Params.LogPointPath, point), vomni.PointCfgData[cd].CfgStr+vomni.LogFileEndInfo)

	fmt.Printf("%q ***** Folderītis  %q\n", point, tmplF)
	fmt.Printf("%q ***** File        %q\n", point, logF)

	return
}
