package pointrun

import "net"

type Runner interface {
	//	GetUDPAddr() (addr net.UDPAddr)
	//	IsActive() (active bool)
	//	LetsGo(chGoOn chan bool, chErr chan error)
	//	LogPointStr(cd int, logStr string)
	//	RotateReAssign() (err error)
	//	Response(msg []string, chDelete chan bool, chErr chan error)
	//	SetUDPAddr(ip string, port int)
	//	WebPointData() (v omnibus.WPointData)
	//	ReceivedWebMsg(msg string, data interface{})
	//	Finish(done chan bool)
	LetsGo(addr net.UDPAddr, chGoOn chan bool, chDone chan int, chErr chan error)
	LogStr(info int, str string)
	StartRotate() (err error)
	GetDone(done int)
	SetUDPAddr(addr net.UDPAddr)
	GetUDPAddr() (addr net.UDPAddr)
}

type PointMsg struct {
	MsgCd  int
	MsgStr string
}

type PointRunners map[int]Runner

type PointRun struct {
	Point PointData
	Run   PointRunners
}

type PointData struct {
	Point   string
	UDPAddr net.UDPAddr
	Type    int
	State   int
}
