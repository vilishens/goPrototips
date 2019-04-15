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

	IPAddressInternal string
	IPAddressExternal string

	IPExternalAddressCmds  []string
	NetExternalRequirement int

	PointConfigOriginalFile string
	PointConfigFile         string
}
