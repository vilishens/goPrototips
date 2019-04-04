package params

type ParamData struct {
	Name string

	LogMainPath string

	PortSSHInternal int
	PortUDPInternal int
	PortWEBInternal int

	RotateMainCfg        string
	RotatePointCfg       string
	RotateRunCfg         string
	RotateRunSecs        int
	RotateStatusFileName string

	WebStaticPrefix string
	WebStaticDir    string
	WebTemplateDir  string

	//###################################

	PointLogPath string

	InternalPort     int
	InternalIPv4     string
	ExternalPort     int
	ExternalIPv4     string
	WebEmail         string
	WebAliveInterval int
	WebEmailMutt     string
	ScriptPath       string
	LogPath          string
	PointModeFiles   map[string]string
	TemplatePath     string
	TemplateExt      string
	DevModes         []string
	WebPort          int
	UDPPort          int
	EventPath        string
	ErrorPath        string
}
