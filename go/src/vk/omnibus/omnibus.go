package omnibus

import (
	"log"
	"os"
	"time"
)

var stepList map[string]bool

var RootErr = make(chan error)
var RootDone = make(chan int)
var StepErr = make(chan error)

var (
	RootPath    string
	LogMainFile *os.File
	LogData     *log.Logger
	LogErr      *log.Logger
	LogFatal    *log.Logger
	LogInfo     *log.Logger
)

// constants for log
const (
	LogFileFlags   = os.O_RDWR | os.O_CREATE | os.O_APPEND
	LogUserPerms   = os.FileMode(0666)
	LogMainPath    = "../log/main/logMain.log"
	LogLoggerFlags = log.LstdFlags | log.LUTC
	LogPrefixData  = "==== DATA === "
	LogPrefixErr   = "!!! ERROR !!! "
	LogPrefixInfo  = "**** INFO *** "
	LogPrefixFatal = "xxx FATAL xxx "
)

const (
	DoneError    = 0x0000010
	DoneReboot   = 0x0000020
	DoneRestart  = 0x0000040
	DoneStop     = 0x0000080
	DonePostStop = 0x0000100
)

const (
	NoNetError           = 0x0000
	NoNetInternal        = 0x0010
	NoNetExternal        = 0x0020
	NetExternalNone      = 0x0000
	NetExternalNice2Have = 0x0001
	NetExternalRequired  = 0x0002
	NetExternalBits      = 0x0003
)

const (
	DelayStepExec             = 10 * time.Millisecond
	DelaySendMessage          = time.Millisecond // time delay between two message send
	DelaySendMessageListEmpty = 3 * time.Millisecond
	DelaySendMessageRepeat    = 500 * time.Millisecond // interval between repeated messages

	DelayWaitMessage = time.Millisecond // time delay between two message waiting

	DelayBetweenIPHello = 1000 * time.Millisecond

	MessageSendRepeatLimit = 3
)

const (
	StepNameConfig      = "step-config"
	StepNameMessages    = "step-messages"
	StepNameNetInfo     = "step-net-info"
	StepNameNetScan     = "step-net-scan"
	StepNameParams      = "step-params"
	StepNamePointConfig = "step-point-cfg"
	StepNamePointReady  = "step-point-ready"
	StepNamePointRun    = "step-point-run"
	StepNameRotateMain  = "step-rotate-main"
	StepNameStart       = "step-start"
	StepNameUDP         = "step-udp"
	StepNameWeb         = "step-web"
)

const (
	DirPermissions         = 0744
	FileNonExecPermissions = 0666
)

const (
	TimeFormat1 = "2006-01-02 15:04:05 -07:00 MST"
)

const (
	CfgDefaultPath      = "../cfg/app/default.cfg"
	CliCfgPathFld       = "path"
	LogRotateStatusFile = "logStatus.status"
)

//################################################# Point run state codes ####################
const (
	PointStateUnknown = 0x000000
)

//################################################# Configuration parameters ####################
const (
	CfgTypeRelayInterval = 0x000001
)

//################################################# Net Info parameters (net/netinfo.go) ####################
const (
	NetInfoRepeats  = 3
	NetInfoInterval = 15 * time.Second
)

//################################################# Message ####################

var MessageNumber int // unique message number (starting from the application launch)

var AllMessages map[int]MessageData

type MessageData struct {
	FieldCount int
}

const (
	MsgCdOutputHelloFromStation = 0x00000001 // Output <station name><msgCd><msgNbr><station UTC seconds><station time offset><stationIP><stationPort>
	MsgCdInputHelloFromPoint    = 0x00000002 // Input  <point name><msgCd><msgNbr><pointIP><pointPort>
	MsgCdOutputSetRelayGpio     = 0x00000004 // Output <point name><msgCd><msgNbr><Gpio><set value>
)

const (
	UDPMessageSeparator = ":::"

	MsgIndexPrefixSender = 0
	MsgIndexPrefixCd     = 1
	MsgIndexPrefixNbr    = 2

	MsgPrefixLen = 3
)

// Hello From Station
const (
	MsgIndexHelloFromStationTime   = 0
	MsgIndexHelloFromStationOffset = 1
	MsgIndexHelloFromStationIP     = 2
	MsgIndexHelloFromStationPort   = 3

	MsgHelloFromStationLen = 4
)

// Hello From Point
const (
	MsgIndexHelloFromPointIP   = 0
	MsgIndexHelloFromPointPort = 1

	MsgHelloFromPointLen = 2
)
