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
	Point string
	State int
	Type  int
	// all point logger  files, key shows bitwise what type of loggers included ("info", "data", ...)
	// The file can have more than one logger (for instance, "info" and "err" info into one file by 2 loggers)
	Logs       []vomni.PointLog
	Index      AllIndex
	UDPAddr    net.UDPAddr
	ChErr      chan error
	ChDone     chan int
	ChMsg      chan string
	CfgDefault vcfg.RelIntervalStruct
	CfgRun     vcfg.RelIntervalStruct
	CfgSaved   vcfg.RelIntervalStruct
}
