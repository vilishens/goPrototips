package pointconfig

import (
	"fmt"
	"strings"
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
	return
}
