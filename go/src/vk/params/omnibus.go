package params

type ParamData struct {
	StationName string

	LogMainPath  string
	LogPointPath string

	PortUDPInternal int
	PortSSHInternal int
	PortWEBInternal int
	PortSSHExternal int
	PortWEBExternal int

	RotateMainCfg        string
	RotatePointCfg       string
	RotateRunCfg         string
	RotateRunSecs        int
	RotateStatusFileName string

	WebStaticPrefix string
	WebStaticDir    string
	WebTemplateDir  string

	IPAddressInternal string
	IPAddressExternal string

	IPExternalAddressCmds  []string
	NetExternalRequirement int

	PointConfigOriginalFile string
	PointConfigFile         string

	SendGridKey         string
	MessageEmailAddress string
}
