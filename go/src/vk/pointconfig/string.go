package pointconfig

import (
	"fmt"
	"strings"
	"time"
	vutils "vk/utils"
)

var offset string

func init() {
	offset = "\t"
}

func (d AllPointCfgData) String() (str string) {

	str = ""
	str += "### Running Configuration ###\n"
	for k, v := range d {
		str += offset + fmt.Sprintf("STATION %q", k) + "\n"
		str += strings.Repeat(offset, 2) + fmt.Sprintf("Config Run") + "\n"
		str += v.CfgRun.String(3)
		str += strings.Repeat(offset, 2) + fmt.Sprintf("Config Saved") + "\n"
		str += v.CfgSaved.String(3)
	}

	return
}

func (d PointCfg) String(offN int) (str string) {

	str += strings.Repeat(offset, offN) + "Relay On/Off Interval:\n"
	str += d.RelInterv.String(offN + 1)

	return
}

func (d RelIntervalStruct) String(offN int) (str string) {

	str += strings.Repeat(offset, offN) + "START\n"
	str += d.Start.String(offN + 1)
	str += strings.Repeat(offset, offN) + "BASE\n"
	str += d.Start.String(offN + 1)
	str += strings.Repeat(offset, offN) + "FINISH\n"
	str += d.Start.String(offN + 1)

	return
}

func (d RelIntervalArray) String(offN int) (str string) {

	for _, v := range d {
		str += v.String(offN)
	}

	return
}

func (d RelInterval) String(offN int) (str string) {

	str += strings.Repeat(offset, offN) + fmt.Sprintf("GPIO: %2d STATE: %d SECONDS %d (%s)\n", d.Gpio, d.State,
		d.Seconds/time.Second, vutils.Duration2ConfInterval(d.Seconds, true))

	return
}