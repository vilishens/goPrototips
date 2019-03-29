package omnibus

import (
	"log"
	"os"
	"time"
)

var RootErr = make(chan error)
var RootDone = make(chan int)

var (
	RootPath    string
	LogMainFile *os.File
	LogErr      *log.Logger
	LogInfo     *log.Logger
	LogData     *log.Logger
	LogFatal    *log.Logger
)

// constants for loag
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
	DoneError   = 0x0000010
	DoneReboot  = 0x0000020
	DoneRestart = 0x0000040
	DoneStop    = 0x0000080
)

const (
	StepExecDelay = 10 * time.Millisecond
)

const (
	StepNameConfig = "step-config"
	StepNameStart  = "step-start"
	StepNameWeb    = "step-web"
)

const (
	DirPermissions = 0744
)

const (
	TimeFormat1 = "2006-01-02 15:04:05 -07:00 MST"
)

const (
	WebPrefix     = "/xK-@eRty$Wj7e$yZ/"
	WebStaticPath = "webstatic/"
)

const (
	CfgDefaultPath = "../cfg/default/default.cfg"
)
