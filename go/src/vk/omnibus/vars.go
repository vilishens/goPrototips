package omnibus

import (
	"log"
	"os"
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

var PointLogData map[int]LogPointData
var PointCfgData map[int]CfgPointData

func init() {

	PointCfgData = make(map[int]CfgPointData)
	PointCfgData[CfgTypeRelayInterval] = CfgPointData{CfgCd: CfgTypeRelayInterval, CfgStr: "relay-interval"}

	PointLogData = make(map[int]LogPointData)
	PointLogData[LogFileData] = LogPointData{LogCd: LogFileData, FileEnd: LogFileEndData, LogPrefix: LogPointPrefixData}
	PointLogData[LogFileErr] = LogPointData{LogCd: LogFileErr, FileEnd: LogFileEndErr, LogPrefix: LogPointPrefixErr}
	PointLogData[LogFileInfo] = LogPointData{LogCd: LogFileInfo, FileEnd: LogFileEndInfo, LogPrefix: LogPointPrefixInfo}
}
