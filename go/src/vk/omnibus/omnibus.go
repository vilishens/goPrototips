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
	StepNameStart = "step-start"
)

//#################################
var MessageTypeLimits = map[string]int{
	MessageOmnibus: 3} // type,point, cmd

const (
	MessageTypeCmd   = "CMD"
	MessageTypeError = "ERROR"
	MessageTypeEvent = "EVENT"
	MessageTypeStart = "START"
	MessageTypeStop  = "STOP"
	MessageOmnibus   = "OMNIBUS"
)

const (
	LogStatusFile    = "logstatus.status"
	PointLogDataFile = "data.log"
	PointLogInfoFile = "info.log"
)

const (
	DIR_PERMISSIONS = 0744

//FILE_PERMISSIONS       = 0644
//GAMMU_CFG_FILE_INDEX   = 3
//GAMMU_PIN_SUBMIT_INDEX = 5
)
