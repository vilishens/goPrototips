package pointrun

import (
	"fmt"
	"sort"
	vomni "vk/omnibus"
)

func AllPointData() (data vomni.WebAllPointData) {

	pts := make(map[string]vomni.WebPointData)

	list := []string{}

	fmt.Println()
	fmt.Println()

	for k, v := range Points {

		list = append(list, k)

		d := vomni.WebPointData{}

		d.Point = k
		d.State = v.Point.State
		d.Type = v.Point.Type

		d.Signed = 0 != (v.Point.State & vomni.PointStateSigned)
		d.Disconnected = 0 != (v.Point.State & vomni.PointStateDisconnected)

		d.CfgList = vomni.CfgListSequence
		d.CfgInfo = webCfgInfo(d.CfgList) // ziņas par konfigurācijām: pāšlaik konfigurāciju nosaukumi

		d.CfgDefault = make(map[int]interface{})
		d.CfgRun = make(map[int]interface{})
		d.CfgSaved = make(map[int]interface{})
		d.CfgState = make(map[int]interface{})
		d.CfgIndex = make(map[int]interface{})

		// get all configurations of this poist
		for _, cc := range d.CfgList {
			cfgDef, cfgRun, cfgSaved, cfgIndex, cfgState := Points[k].Run[cc].GetCfgs()

			d.CfgDefault[cc], d.CfgRun[cc], d.CfgSaved[cc], d.CfgIndex[cc], d.CfgState[cc] =
				cfgDef, cfgRun, cfgSaved, cfgIndex, cfgState //Points[k].Run[cc].GetCfgs()

			_, _, _ = cfgDef, cfgRun, cfgSaved
		}

		/*
			cfgDef, cfgRun, cfgSaved := Points[k].Run[vomni.CfgTypeRelayInterval].GetCfgs()

			back, err := json.Marshal(cfgDef)
			if nil != err {
				panic(err)
				return
			}

			json.Unmarshal(back, &cfgDef)

			d.CfgDefault = cfgDef

			d.CfgDefault, d.CfgRun, d.CfgSaved = Points[k].Run[vomni.CfgTypeRelayInterval].GetCfgs()
		*/
		fmt.Printf("Point %q Signed %t Disconn %t\n", d.Point, d.Signed, d.Disconnected)

		pts[k] = d
	}

	fmt.Println()
	fmt.Println()

	sort.Strings(list)

	data.List = list
	data.Data = pts

	return
}

func webCfgInfo(list []int) (d map[int]vomni.CfgPlusData) {

	d = make(map[int]vomni.CfgPlusData)

	for _, v := range list {
		dd := vomni.CfgPlusData{}

		dd.Name = vomni.PointCfgData[v].CfgStr

		d[v] = dd
	}

	return
}

func WebSent(todo int, point string, data interface{}) {

	cmd := todo & vomni.PointCmdBits
	switch cmd {
	case vomni.PointCmdLoadCfgIntoPoint:
		cfg := todo & vomni.PointCmdOptionBits
		Points[point].Run[cfg].ReceiveCfg(cmd, data)

	default:
		str := fmt.Sprintf("\n\nDon't know what to do with %08X for %s\n\n", todo, point)
		panic(str)
	}
}
