package runrelayinterval

import (
	"net"
	vcfg "vk/pointconfig"
)

type RunData struct {
	Point    string
	State    int
	Type     int
	UDPAddr  net.UDPAddr
	Cfg      vcfg.RelIntervalStruct
	CfgSaved vcfg.RelIntervalStruct
}
