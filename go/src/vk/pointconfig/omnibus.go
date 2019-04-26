package pointconfig

import (
	"time"
)

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
}

type PointCfgData struct {
	List     int      // a field contains bits of available configurations of the point
	Cfg      PointCfg // configuration to use
	CfgSaved PointCfg // saved configuration
}

type AllPointCfgData map[string]PointCfgData

//********* Relay On/Off Interval JSON Configuration ****************************
type CfgJSONData struct {
	RelIntervalJSON CfgRelIntervalPoints `json:"RelayOnOffIntervals"`
}

type CfgRelInterval struct {
	Gpio     string `json:"Gpio"`
	State    string `json:"State"`
	Interval string `json:"Interval"`
}

type CfgRelIntervalArray []CfgRelInterval

type CfgRelIntervalStruct struct {
	Start  CfgRelIntervalArray `json:"Start"`  // array of the point relay default settings (used at the start and exit)
	Base   CfgRelIntervalArray `json:"Base"`   // array of the point relay setting sequences (used between the start and exit)
	Finish CfgRelIntervalArray `json:"Finish"` // array of the point relay setting sequences (used between the start and exit)
}

type CfgRelIntervalPoints map[string]CfgRelIntervalStruct
