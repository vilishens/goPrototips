package cfg

type CfgData struct {
	Name string `json:"Name"`

	//==========================================================================================

	//	LogMainFile         string `json:"mainLogFile"`
	PointDefaultCfgFile string `json:"pointDefaultCfgFile"`
	PointCfgFile        string `json:"pointCfgFile"`
	PointLogPath        string `json:"pointLogPath"`
	RotateMainCfg       string `json:"rotateMainCfg"`
	RotatePointCfg      string `json:"rotatePointCfg"`
	RotateRunCfg        string `json:"rotateRunCfg"`
	RotateRunSecs       string `json:"rotateRunIntervalSecs"`

	//#########################################

	InternalPort     string `json:"internalPort"`
	InternalIP       string `json:"internalIP"`
	ExternalPort     string `json:"externalPort"`
	WebEmail         string `json:"webEmail"`
	WebAliveInterval string `json:"webAliveInternal"`
	WebEmailMutt     string `json:"webEmailMutt"`
	ScriptPath       string `json:"scriptPath"`
	LogPath          string `json:"logPath"`
	WebPort          string `json:"webPort"`
	UDPPort          string `json:"udpPort"`
	EventPath        string `json:"eventPath"`
	TemplatePath     string `json:"templatePath"`
	TemplateExt      string `json:"templateExt"`
	ErrorPath        string `json:"errorPath"`
}

type CfgFinalData struct {
	Name                string
	LogMainFile         string
	PointDefaultCfgFile string
	PointCfgFile        string
	PointLogPath        string
	RotateMainCfg       string
	RotatePointCfg      string
	RotateRunCfg        string
	RotateRunSecs       int

	//#####################################################

	UDPPort int

	InternalPort     int
	InternalIP       string
	ExternalPort     int
	WebEmail         string
	WebAliveInterval int
	WebEmailMutt     string
	ScriptPath       string
	LogPath          string
	WebPort          int
	EventPath        string
	TemplatePath     string
	TemplateExt      string
	ErrorPath        string
}
