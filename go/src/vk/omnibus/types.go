package omnibus

import (
	"log"
)

type PointLogger struct {
	LogFile string
	LogTmpl string
	LogPrefix string
	Logger  *log.Logger
}

type PointLogs map[string]PointLogger

type 