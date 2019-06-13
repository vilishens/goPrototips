package pointready

import (
	vomni "vk/omnibus"
	vpointconfig "vk/pointconfig"
	vrunrelayinterval "vk/run/relayinterval"
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
			logs, err := pointLoggers(d.Point, d.Type)

			if nil != err {
				vomni.RootErr <- err
			}

			d.Logs = make(map[int]vomni.PointLog)
			d.Logs[d.Type] = logs

			vrunrelayinterval.RunningPoints[k] = d
		}
	}
}
