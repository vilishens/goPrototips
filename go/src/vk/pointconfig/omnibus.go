package pointconfig

import (
	"time"
)

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
	CfgRun   PointCfg
	CfgSaved PointCfg
}

type AllPointCfgData map[string]PointCfgData

// JSON data

type CfgJSONData struct {
	RelIntervalJSON CfgRelIntervalPoints `json:"RelayOnOffIntervals"`
}

type CfgRelInterval struct {
	Gpio     string `json:"Gpio"`
	State    string `json:"State"`
	Interval string `json:"Interval"`
}

type CfgRelIntervalArr []CfgRelInterval

type CfgRelIntervalStruct struct {
	Start  CfgRelIntervalArr `json:"Start"`  // array of the point relay default settings (used at the start and exit)
	Base   CfgRelIntervalArr `json:"Base"`   // array of the point relay setting sequences (used between the start and exit)
	Finish CfgRelIntervalArr `json:"Finish"` // array of the point relay setting sequences (used between the start and exit)
}

type CfgRelIntervalPoints map[string]CfgRelIntervalStruct
