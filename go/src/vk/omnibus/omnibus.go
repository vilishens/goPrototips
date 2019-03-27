package omnibus

import (
	"log"
	"os"
)

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
	DoneError   = 0x0000001
	DoneRestart = 0x0000002
	DoneStop    = 0x0000004
)
