package cfg

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	vomni "vk/omnibus"
	vutils "vk/utils"
)

var Final CfgFinalData

func init() {
	Final.Name = ""

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

	/*
		// customer
		//	if err = loadCfg(false); nil != err {
		//		chErr <- err
		//		return
		//	}

		// prepare log file
		if err = prepareLogFile(); nil != err {
			chErr <- err
			return
		}

	*/
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

	return
}

//###################################################################################
//###################################################################################
