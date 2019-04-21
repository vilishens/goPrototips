package cfg

type CfgData struct {
	StationName string `json:"StationName"`

	RotateMainCfg        string `json:"RotateMainCfg"`
	RotatePointCfg       string `json:"RotatePointCfg"`
	RotateRunCfg         string `json:"RotateRunCfg"`
	RotateRunSecs        string `json:"RotateRunIntervalSecs"`
	RotateStatusFileName string `json:"RotateStatusFileName"`

	PortSSHInternal string `json:"PortSSHInternal"`
	PortUDPInternal string `json:"PortUDPInternal"`
	PortWEBInternal string `json:"PortWEBInternal"`

	WebStaticPrefix string `json:"WEBStaticPrefix"`
	WebStaticDir    string `json:"WEBStaticDir"`
	WebTemplateDir  string `json:"WEBTemplateDir"`

	IPExternalAddressCmds  []string `json:"IPExternalAddressCmds"`
	NetExternalRequirement string   `json:"NetExternalRequired"`

	PointConfigOriginalFile string `json:"PointConfigOriginalFile"`
	PointConfigFile         string `json:"PointConfigFile`
}

type CfgFinalData struct {
	StationName string
	LogMainPath string

	RotateMainCfg        string
	RotatePointCfg       string
	RotateRunCfg         string
	RotateRunSecs        int
	RotateStatusFileName string

	PortSSHInternal int
	PortUDPInternal int
	PortWEBInternal int

	WebStaticPrefix string
	WebStaticDir    string
	WebTemplateDir  string

	IPExternalAddressCmds  []string
	NetExternalRequirement int

	PointConfigOriginalFile string
	PointConfigFile         string
}
