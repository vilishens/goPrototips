package params

import (
	vconf "vk/cfg"
)

var Params ParamData

func init() {

	Params = ParamData{}

	Params.StationName = ""

	Params.LogMainPath = ""
	Params.LogPointPath = ""

	Params.PortUDPInternal = -1
	Params.PortSSHInternal = -1
	Params.PortWEBInternal = -1
	Params.PortSSHExternal = -1
	Params.PortWEBExternal = -1

	Params.RotateMainTmpl = ""
	Params.RotatePointDataTmpl = ""
	Params.RotatePointInfoTmpl = ""
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

	Params.SendGridKey = ""
	Params.MessageEmailAddress = ""
}

func Put(chDone chan bool, chErr chan error) {

	err := error(nil)

	data := vconf.Final

	if "" != data.StationName {
		Params.StationName = data.StationName
	}

	if "" != data.LogMainPath {
		Params.LogMainPath = data.LogMainPath
	}
	if "" != data.LogPointPath {
		Params.LogPointPath = data.LogPointPath
	}

	if 0 <= data.PortUDPInternal {
		Params.PortUDPInternal = data.PortUDPInternal
	}
	if 0 <= data.PortSSHInternal {
		Params.PortSSHInternal = data.PortSSHInternal
	}
	if 0 <= data.PortWEBInternal {
		Params.PortWEBInternal = data.PortWEBInternal
	}
	if 0 <= data.PortSSHExternal {
		Params.PortSSHExternal = data.PortSSHExternal
	}
	if 0 <= data.PortWEBExternal {
		Params.PortWEBExternal = data.PortWEBExternal
	}

	if "" != data.RotateMainTmpl {
		Params.RotateMainTmpl = data.RotateMainTmpl
	}
	if "" != data.RotatePointDataTmpl {
		Params.RotatePointDataTmpl = data.RotatePointDataTmpl
	}
	if "" != data.RotatePointInfoTmpl {
		Params.RotatePointInfoTmpl = data.RotatePointInfoTmpl
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

	// SendGrid key and email address
	if (nil == err) && ("" != data.SendGridKey) {
		Params.SendGridKey = data.SendGridKey
	}
	if (nil == err) && ("" != data.MessageEmailAddress) {
		Params.MessageEmailAddress = data.MessageEmailAddress
	}

	if nil != err {
		chErr <- err
	} else {
		chDone <- true
	}
}
