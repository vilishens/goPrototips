package pointconfig

type AllCfgData struct {
	Default     CfgFileData
	DefaultJSON CfgFileJSON
	Running     CfgFileData
	RunningJSON CfgFileJSON
}

type CfgFileData map[string]PointCfgData

type CfgFileJSON map[string]JSONPointData

type JSONPointData struct {
	RelIntervalJSON    JSONRelIntervalStruct `json:"RelayOnOffIntervals"`
	TempRelayJSON      JSONTempRelay         `json:"TempRelay"`
	TempRelayArrayJSON []JSONTempRelay       `json:"TempRelayArray"`
}

type JSONData map[string]JSONPointData
