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

var MessageNumber int // unique message number (starting from the application launch)

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
	UDPMessageSeparator = ":::"
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
)

const (
	StepNameConfig     = "step-config"
	StepNameMessages   = "step-messages"
	StepNameNetInfo    = "step-net-info"
	StepNameParams     = "step-params"
	StepNameRotateMain = "step-rotate-main"
	StepNameStart      = "step-start"
	StepNameUDP        = "step-udp"
	StepNameWeb        = "step-web"
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
