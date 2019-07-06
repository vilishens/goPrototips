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
