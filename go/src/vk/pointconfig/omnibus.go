package pointconfig

import (
	"time"
)

type AllCfgData struct {
	Default     CfgFileData
	DefaultJSON CfgFileJSON
	Running     CfgFileData
	RunningJSON CfgFileJSON
}

type CfgFileData map[string]PointCfgData

type CfgFileJSON map[string]CfgJSONPointData

//##############################################################################
//######### RELAY ON/OFF INTERVAL ##############################################
//##############################################################################

//********* Relay On/Off Interval Run Configuration ****************************
type RelInterval struct {
	Gpio    int
	State   int
	Seconds time.Duration
}

type RelIntervalArray []RelInterval

type RelIntervalStruct struct {
	Start  RelIntervalArray
	Base   RelIntervalArray
	Finish RelIntervalArray
}

type PointCfg struct {
	RelInterv RelIntervalStruct
	TempRelay []RunTempRelay
}

type PointCfgData struct {
	List     int      // a field contains bits of available configurations of the point
	Cfg      PointCfg // configuration to use
	CfgSaved PointCfg // saved configuration
}

type AllPointCfgData map[string]PointCfgData

//********* Relay On/Off Interval JSON Configuration ****************************
type CfgJSONData map[string]CfgJSONPointData

type CfgJSONPointData struct {
	RelIntervalJSON    CfgRelIntervalStruct `json:"RelayOnOffIntervals"`
	TempRelayJSON      JSONTempRelay        `json:"TempRelay"`
	TempRelayArrayJSON []JSONTempRelay      `json:"TempRelayArray"`
}

type CfgRelInterval struct {
	Gpio     string `json:"Gpio"`
	State    string `json:"State"`
	Interval string `json:"Interval"`
}

type CfgRelIntervalStruct struct {
	Start  CfgRelIntervalArray `json:"Start"`  // array of the point relay default settings (used at the start and exit)
	Base   CfgRelIntervalArray `json:"Base"`   // array of the point relay setting sequences (used between the start and exit)
	Finish CfgRelIntervalArray `json:"Finish"` // array of the point relay setting sequences (used between the start and exit)
}

type CfgRelIntervalArray []CfgRelInterval
