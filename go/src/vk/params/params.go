package params

import (
	vconf "vk/cfg"
)

var Params ParamData

func init() {

	Params = ParamData{}

	Params.Name = ""

	Params.LogMainPath = ""

	Params.PortSSHInternal = -1
	Params.PortUDPInternal = -1
	Params.PortWEBInternal = -1

	Params.RotateMainCfg = ""
	Params.RotatePointCfg = ""
	Params.RotateRunCfg = ""
	Params.RotateRunSecs = -1
	Params.RotateStatusFileName = ""

	Params.WebStaticPrefix = ""
	Params.WebStaticDir = ""
	Params.WebTemplateDir = ""

	Params.IPAddressInternal = ""
	Params.IPAddressExternal = ""

	Params.IPExternalAddressCmds = []string{} // commands to find the station external IP address
	Params.NetExternalRequirement = 0         // no the external net required at this moment

	Params.PointConfigOriginalFile = ""
	Params.PointConfigFile = ""
}

func Put(chDone chan bool, chErr chan error) {

	err := error(nil)

	data := vconf.Final

	if "" != data.Name {
		Params.Name = data.Name
	}

	if "" != data.LogMainPath {
		Params.LogMainPath = data.LogMainPath
	}

	if 0 <= data.PortSSHInternal {
		Params.PortSSHInternal = data.PortSSHInternal
	}
	if 0 <= data.PortSSHInternal {
		Params.PortUDPInternal = data.PortUDPInternal
	}
	if 0 <= data.PortSSHInternal {
		Params.PortWEBInternal = data.PortWEBInternal
	}

	if "" != data.RotateMainCfg {
		Params.RotateMainCfg = data.RotateMainCfg
	}
	if "" != data.RotatePointCfg {
		Params.RotatePointCfg = data.RotatePointCfg
	}
	if "" != data.RotateRunCfg {
		Params.RotateRunCfg = data.RotateRunCfg
	}
	if 0 <= data.RotateRunSecs {
		Params.RotateRunSecs = data.RotateRunSecs
	}
	if "" != data.RotateStatusFileName {
		Params.RotateStatusFileName = data.RotateStatusFileName
	}

	if "" != data.WebStaticPrefix {
		Params.WebStaticPrefix = data.WebStaticPrefix
	}
	if "" != data.WebStaticDir {
		Params.WebStaticDir = data.WebStaticDir
	}
	if "" != data.WebTemplateDir {
		Params.WebTemplateDir = data.WebTemplateDir
	}

	if 0 <= data.PortSSHInternal {
		Params.PortWEBInternal = data.PortWEBInternal
	}

	if (nil == err) && (0 < len(data.IPExternalAddressCmds)) {
		Params.IPExternalAddressCmds = make([]string, len(data.IPExternalAddressCmds))
		copy(Params.IPExternalAddressCmds, data.IPExternalAddressCmds)
	}
	if (nil == err) && (0 < data.NetExternalRequirement) {
		Params.NetExternalRequirement = data.NetExternalRequirement
	}

	// point config file
	if (nil == err) && ("" != data.PointConfigOriginalFile) {
		Params.PointConfigOriginalFile = data.PointConfigOriginalFile
	}
	if (nil == err) && ("" != data.PointConfigFile) {
		Params.PointConfigFile = data.PointConfigFile
	}

	if nil != err {
		chErr <- err
	} else {
		chDone <- true
	}
}
