package runrelayinterval

import (
	"net"
	vomni "vk/omnibus"
	vcfg "vk/pointconfig"
)

type AllIndex struct {
	Start  int
	Base   int
	Finish int
}

type RunData struct {
	Point    string
	State    int
	Type     int
	Logs     map[int]vomni.PointLog
	Index    AllIndex
	UDPAddr  net.UDPAddr
	ChError  chan error
	ChDone   chan int
	Cfg      vcfg.RelIntervalStruct
	CfgSaved vcfg.RelIntervalStruct
}
