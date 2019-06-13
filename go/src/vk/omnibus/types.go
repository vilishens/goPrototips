package omnibus

import (
	"log"
)

type PointLog struct {
	LogFile string
	LogTmpl string
	Loggers map[int]PointLogger
}

type PointLogger struct {
	LogPrefix string
	Logger    *log.Logger
}

type LogCfg struct {
	File string
	List []string
}
