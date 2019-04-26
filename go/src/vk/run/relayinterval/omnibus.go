package runrelayinterval

import (
	vcfg "vk/pointconfig"
)

type RunData struct {
	State    int
	Type     int
	Cfg      vcfg.RelIntervalStruct
	CfgSaved vcfg.RelIntervalStruct
}
