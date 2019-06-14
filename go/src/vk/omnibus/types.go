package omnibus

import (
	"log"
)

// the point log data file configuration
type PointLog struct {
	LogFile string              // the full path of the data file
	LogTmpl string              // the full path of the rotate configuration template file
	Loggers map[int]PointLogger // all loggers linked to the data file with the key of the logger bitwise code
}

// the logger configuration
type PointLogger struct {
	LogPrefix string      // the prefix
	Logger    *log.Logger // the logger
}

type LogPointData struct {
	LogCd     int
	FileEnd   string
	LogPrefix string
}

type CfgPointData struct {
	CfgCd  int
	CfgStr string
}

//#################################################################
//#################################################################
//#################################################################

// vai šitas vajadzīgs????
type LogCfg struct {
	File string
	List []string
}

//??????????????????????
var LogPointInfo map[int]LogCfg
