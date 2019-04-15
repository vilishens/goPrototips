package pointconfig

import (
	"time"
)

type RelOnOffInterval struct {
	Gpio    int
	State   int
	Seconds time.Duration
}

type RelOnOffIntervalArray []RelOnOffInterval

type RelOnOffIntervalStruct struct {
	Start  RelOnOffIntervalArray
	Base   RelOnOffIntervalArray
	Finish RelOnOffIntervalArray
}

type PointCfg struct {
	RelOnOffInterv RelOnOffIntervalStruct
}

type PointCfgData struct {
	CfgRun   PointCfg
	CfgSaved PointCfg
}

type AllPointCfgData map[string]PointCfgData

// JSON data

type CfgJSONData struct {
	RelOnOffIntervalJSON CfgRelOnOffIntervalPoints `json:"RelayOnOffIntervals"`
}

type CfgRelOnOffInterval struct {
	Gpio     string `json:"Gpio"`
	State    string `json:"State"`
	Interval string `json:"Interval"`
}

type CfgRelOnOffIntervalArr []CfgRelOnOffInterval

type CfgRelOnOffIntervalStruct struct {
	Start  CfgRelOnOffIntervalArr `json:"Start"`  // array of the point relay default settings (used at the start and exit)
	Base   CfgRelOnOffIntervalArr `json:"Base"`   // array of the point relay setting sequences (used between the start and exit)
	Finish CfgRelOnOffIntervalArr `json:"Finish"` // array of the point relay setting sequences (used between the start and exit)
}

type CfgRelOnOffIntervalPoints map[string]CfgRelOnOffIntervalStruct
