package cfg

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	vomni "vk/omnibus"
	vutils "vk/utils"
)

var Final CfgFinalData

func init() {

	Final.Name = ""

	Final.LogMainPath = ""

	Final.PortSSHInternal = -1
	Final.PortUDPInternal = -1
	Final.PortWEBInternal = -1

	Final.RotateMainCfg = ""
	Final.RotatePointCfg = ""
	Final.RotateRunCfg = ""
	Final.RotateRunSecs = -1
	Final.RotateStatusFileName = ""

	Final.WebStaticPrefix = ""
	Final.WebStaticDir = ""
	Final.WebTemplateDir = ""

	Final.IPExternalAddressCmds = []string{}
	Final.NetExternalRequirement = -1

	/*
		Final.LogMainFile = ""
		Final.PointDefaultCfgFile = ""
		Final.PointCfgFile = ""
		Final.PointLogPath = ""

		Final.RotateMainCfg = ""
		Final.RotatePointCfg = ""
		Final.RotateRunCfg = ""

		Final.InternalPort = -5
		Final.InternalIP = ""
		Final.ExternalPort = -11
		Final.WebEmail = ""
		Final.WebAliveInterval = -7
		Final.ScriptPath = ""
		Final.LogPath = ""
		Final.WebPort = -11
		Final.TemplatePath = ""
		Final.TemplateExt = ""
		Final.ErrorPath = ""
	*/
}

func Cfg(chDone chan bool, chErr chan error) {

	locDone := make(chan bool)
	locErr := make(chan error)

	go load(locDone, locErr)

	select {
	case err := <-locErr:
		chErr <- err
	case <-locDone:
		chDone <- true
	}
}

func load(chDone chan bool, chErr chan error) {

	var err error

	if err = loadCfg(); nil != err {
		vutils.LogStr(vomni.LogErr, err.Error())
		chErr <- err
		return
	}

	chDone <- true
}

func loadCfg() (err error) {

	full := ""
	err = error(nil)
	if full, err = cfgPath(); nil != err {
		return
	}

	if "" == full {
		err = fmt.Errorf("There is no Application configuration")
		return vutils.ErrFuncLine(err)
	}

	data, err := readCfg(full)
	if nil != err {
		return
	}

	if err = data.Put(); nil != err {
		return
	}

	return
}

func cfgPath() (path string, err error) {
	// configuration data path found in CLI flags
	cpath := flag.Lookup(vomni.CliCfgPathFld).Value.String()

	if "" == cpath {
		return
	}

	path = vutils.FileAbsPath(cpath, "")

	ok := false
	if ok, err = vutils.PathExists(path); !ok {
		err = fmt.Errorf("File \"%s\" doesn't exist", path)
		err = vutils.ErrFuncLine(err)
	}

	return
}

func readCfg(full string) (data CfgData, err error) {

	data = CfgData{}
	if ok, err := vutils.PathExists(full); !ok {
		return data, vutils.ErrFuncLine(err)
	}

	raw, err := ioutil.ReadFile(full)
	if err != nil {
		return data, vutils.ErrFuncLine(err)
	}

	if err = json.Unmarshal(raw, &data); nil != err {
		return data, vutils.ErrFuncLine(err)
	}

	return data, err
}

func (c *CfgData) Put() (err error) {

	if (nil == err) && ("" != c.Name) {
		Final.Name = c.Name
	}

	// hard coded Main log file path
	Final.LogMainPath = filepath.Join(vomni.RootPath, vomni.LogMainPath)

	// rotation of logs
	if (nil == err) && ("" != c.RotateMainCfg) {
		Final.RotateMainCfg = c.RotateMainCfg
	}
	if (nil == err) && ("" != c.RotatePointCfg) {
		Final.RotatePointCfg = c.RotatePointCfg
	}
	if (nil == err) && ("" != c.RotateRunCfg) {
		Final.RotateRunCfg = c.RotateRunCfg
	}
	if (nil == err) && ("" != c.RotateRunSecs) {
		Final.RotateRunSecs, err = strconv.Atoi(c.RotateRunSecs)
	}
	if (nil == err) && ("" != c.RotateStatusFileName) {
		Final.RotateStatusFileName = c.RotateStatusFileName
	}

	// internal ports
	if (nil == err) && ("" != c.PortSSHInternal) {
		Final.PortSSHInternal, err = strconv.Atoi(c.PortSSHInternal)
	}
	if (nil == err) && ("" != c.PortUDPInternal) {
		Final.PortUDPInternal, err = strconv.Atoi(c.PortUDPInternal)
	}
	if (nil == err) && ("" != c.PortWEBInternal) {
		Final.PortWEBInternal, err = strconv.Atoi(c.PortWEBInternal)
	}

	// WEB configuration
	if (nil == err) && ("" != c.WebStaticPrefix) {
		Final.WebStaticPrefix = c.WebStaticPrefix
	}
	if (nil == err) && ("" != c.WebStaticDir) {
		Final.WebStaticDir = c.WebStaticDir
	}

	if (nil == err) && ("" != c.WebTemplateDir) {
		Final.WebTemplateDir = c.WebTemplateDir
	}

	// External net settings
	if (nil == err) && (0 < len(c.IPExternalAddressCmds)) {
		Final.IPExternalAddressCmds = make([]string, len(c.IPExternalAddressCmds))
		copy(Final.IPExternalAddressCmds, c.IPExternalAddressCmds)
	}
	if (nil == err) && ("" != c.NetExternalRequirement) {
		Final.NetExternalRequirement, err = strconv.Atoi(c.NetExternalRequirement)
	}

	return
}
