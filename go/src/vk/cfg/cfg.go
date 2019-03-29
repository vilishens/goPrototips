package cfg

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	vomni "vk/omnibus"
	vutils "vk/utils"
)

var Final CfgFinalData

func init() {
	Final.Name = ""
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

/*

func Cfg(chDone chan bool, chErr chan error, wg *sync.WaitGroup) {

	locDone := make(chan bool)
	locErr := make(chan error)

	//	defer wg.Done()

	go load(locDone, locErr)

	select {
	case err := <-locErr:
		chErr <- err
	case <-locDone:
		fmt.Printf(">>>>>>>>>>>>>>>>>>>>> SAHAR \"%s\"\n", Final.LogPath)
		wg.Done()
		chDone <- true
	}

	fmt.Printf(">>>>>>>>>>>>>>>>>>>>> turban \"%s\"\n", Final.LogPath)
	//	wg.Done()
	//chDone <- true

	/ *
		// factory
		if err = loadCfg(true); nil != err {
			return
		}

		// customer
		if err = loadCfg(false); nil != err {
			return
		}

		fmt.Printf("#### FINAL ####\n%s\n", Final.String())
	* /
	return
}
*/

func load(chDone chan bool, chErr chan error) {

	err := error(nil)

	// factory
	if err = loadCfg(true); nil != err {
		chErr <- err
		return
	}

	// customer
	if err = loadCfg(false); nil != err {
		chErr <- err
		return
	}

	// prepare log file
	if err = prepareLogFile(); nil != err {
		chErr <- err
		return
	}

	fmt.Printf("#### FINAL ####\n%s\n", Final.String())

	chDone <- true
}

func prepareLogFile() (err error) {
	// Set Rotate config for the main log
	if err = vutils.SetRotateCfg(vomni.LogMainPath, Final.RotateMainCfg, Final.RotateRunCfg, true); nil != err {
		return vutils.ErrFuncLine(fmt.Errorf("Main rotate file error - %v", err))
	}

	return
}

func (c CfgFinalData) String() (str string) {

	str = ""

	flds := make(map[string]string)

	flds["NAME"] = c.Name
	flds["Log File"] = c.LogMainFile
	flds["Point Default Cfg File"] = c.PointDefaultCfgFile
	flds["Point Cfg File"] = c.PointCfgFile
	flds["Point Log Path"] = c.PointLogPath

	flds["Rotate Main Cfg"] = c.RotateMainCfg
	flds["Rotate Point Cfg"] = c.RotatePointCfg
	flds["Rotate Run Cfg"] = c.RotateRunCfg
	flds["Rotate Run Seconds"] = strconv.Itoa(c.RotateRunSecs)

	flds["Internal IP"] = c.InternalIP
	flds["Internal Port"] = strconv.Itoa(c.InternalPort)
	flds["External Port"] = strconv.Itoa(c.ExternalPort)
	flds["WEB Email"] = c.WebEmail
	flds["WEB Alive Interval"] = strconv.Itoa(c.WebAliveInterval)
	flds["WEB Email Mutt"] = c.WebEmailMutt
	flds["Log Path"] = c.LogPath
	flds["Event Path"] = c.EventPath
	flds["Template Path"] = c.TemplatePath
	flds["Template Ext"] = c.TemplateExt
	flds["Error Path"] = c.ErrorPath
	flds["UDPPort"] = strconv.Itoa(c.UDPPort)

	seq := []string{"NAME", "Log File", "Point Defaul Cfg File", "Point Cfg File",
		"Point Log Path", "Rotate Main Cfg", "Rotate Point Cfg", "Rotate Run Cfg", "Rotate Run Seconds",
		"Internal IP", "Internal Port", "External Port",
		"WEB Email", "WEB Alive Interval", "WEB Email Mutt", "Log Path",
		"Event Path", "Error Path", "Template Path", "Template Ext"}

	for _, k := range seq {
		str += fmt.Sprintf("%-18s : %s\n", k, flds[k])
	}

	return str
}

func readCfg(full string) (data cfgData, err error) {

	data = cfgData{}
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

	fmt.Printf(" CFG >>> \n%s\n", data.String())

	return data, err
}

func loadCfg(factory bool) (err error) {

	full := ""
	err = error(nil)
	if full, err = cfgPath(factory); nil != err || "" == full {
		return
	}

	data := cfgData{}
	if data, err = readCfg(full); nil != err {
		return
	}

	//	readCfg(full)
	if factory {
		fmt.Printf("Factory *****\n%s\n", data.String())
	} else {
		fmt.Printf("Customer *****\n%s\n", data.String())
	}

	if err = data.Put(); nil != err {
		return
	}

	return
}

func cfgPath(factory bool) (path string, err error) {
	defPath := flag.Lookup(vomni.CfgFldPath).DefValue

	cPath := flag.Lookup(vomni.CfgFldPath).Value.String()

	usePath := ""
	useFile := ""

	if factory {
		usePath = defPath
	} else {
		if !(defPath == cPath) {
			usePath = cPath
		}
	}

	if "" == usePath {
		return
	}

	path = vutils.FileAbsPath(usePath, useFile)

	ok := false
	if ok, err = vutils.PathExists(path); !ok {
		err = fmt.Errorf("File \"%s\" doesn't exist", path)
	}
	err = vutils.ErrFuncLine(err)

	return
}

func (c *cfgData) Put() (err error) {

	if (nil == err) && ("" != c.Name) {
		Final.Name = c.Name
	}
	//	if (nil == err) && ("" != c.LogMainFile) {
	//		Final.LogMainFile = vutils.FileAbsPath(c.LogMainFile, "")
	//	}

	if (nil == err) && ("" != c.RotateMainCfg) {
		Final.RotateMainCfg = vutils.FileAbsPath(c.RotateMainCfg, "")
	}
	if (nil == err) && ("" != c.RotatePointCfg) {
		Final.RotatePointCfg = vutils.FileAbsPath(c.RotatePointCfg, "")
	}
	if (nil == err) && ("" != c.RotateRunCfg) {
		Final.RotateRunCfg = vutils.FileAbsPath(c.RotateRunCfg, "")
	}

	if (nil == err) && ("" != c.PointDefaultCfgFile) {
		Final.PointDefaultCfgFile = vutils.FileAbsPath(c.PointDefaultCfgFile, "")
	}
	if (nil == err) && ("" != c.PointCfgFile) {
		Final.PointCfgFile = vutils.FileAbsPath(c.PointCfgFile, "")
	}
	if (nil == err) && ("" != c.PointLogPath) {
		Final.PointLogPath = c.PointLogPath
	}

	//###############################################
	//###############################################

	if (nil == err) && ("" != c.InternalPort) {
		Final.InternalPort, err = strconv.Atoi(c.InternalPort)
		err = vutils.ErrFuncLine(err)
	}
	if (nil == err) && ("" != c.InternalIP) {
		Final.InternalIP = c.InternalIP
	}
	if (nil == err) && ("" != c.ExternalPort) {
		Final.ExternalPort, err = strconv.Atoi(c.ExternalPort)
		err = vutils.ErrFuncLine(err)
	}
	if (nil == err) && ("" != c.WebEmail) {
		Final.WebEmail = c.WebEmail
	}
	if (nil == err) && ("" != c.WebAliveInterval) {
		Final.WebAliveInterval, err = strconv.Atoi(c.WebAliveInterval)
		err = vutils.ErrFuncLine(err)
	}
	if (nil == err) && ("" != c.WebEmailMutt) {
		Final.WebEmailMutt = c.WebEmailMutt
	}
	if (nil == err) && ("" != c.ScriptPath) {
		Final.ScriptPath = c.ScriptPath
	}

	if (nil == err) && ("" != c.UDPPort) {
		Final.UDPPort, err = strconv.Atoi(c.UDPPort)
	}

	if (nil == err) && ("" != c.RotateRunSecs) {
		Final.RotateRunSecs, err = strconv.Atoi(c.RotateRunSecs)
		err = vutils.ErrFuncLine(err)
	}

	fmt.Println("----------------------- PereTerpi   ", c.LogPath)

	fmt.Println("----------------------- CfgOfPoints ", Final.PointCfgFile)
	fmt.Println("-----------------------  CfgLogPath ", Final.PointLogPath)

	if (nil == err) && ("" != c.LogPath) {
		Final.LogPath = c.LogPath
	}

	if (nil == err) && ("" != c.WebPort) {
		Final.WebPort, err = strconv.Atoi(c.WebPort)
	}
	if (nil == err) && ("" != c.EventPath) {
		Final.EventPath = c.EventPath
	}
	if (nil == err) && ("" != c.TemplatePath) {
		Final.TemplatePath = c.TemplatePath
	}
	if (nil == err) && ("" != c.TemplateExt) {
		Final.TemplateExt = c.TemplateExt
	}
	if (nil == err) && ("" != c.ErrorPath) {
		Final.ErrorPath = c.ErrorPath
	}

	return
}

func (c *cfgData) String() (str string) {

	str = ""

	flds := make(map[string]string)

	flds["NAME"] = c.Name
	flds["Log Main File"] = vomni.LogMainPath
	flds["Point Log Path"] = c.PointLogPath

	flds["Rotate Main Cfg"] = c.RotateMainCfg
	flds["Rotate Point Cfg"] = c.RotatePointCfg
	flds["Rotate Run Cfg"] = c.RotateRunCfg

	flds["Internal IP"] = c.InternalIP
	flds["Internal Port"] = c.InternalPort
	flds["External Port"] = c.ExternalPort
	flds["WEB Email"] = c.WebEmail
	flds["WEB Alive"] = c.WebAliveInterval
	flds["Script Path"] = c.ScriptPath
	flds["Event Path"] = c.EventPath
	flds["Template Path"] = c.TemplatePath
	flds["Template Ext"] = c.TemplateExt
	flds["UDP Port"] = c.UDPPort

	seq := []string{"NAME", "Log Main File", "Point Log Path",
		"Rotate Main Cfg", "Rotate Point Cfg", "Rotate Run Cfg",
		"Internal IP", "Internal Port", "External Port",
		"WEB Email", "WEB Alive", "Script Path", "Event Path", "Template Path", "Template Ext",
		"UDP Port"}

	for _, k := range seq {
		str += fmt.Sprintf("%-18s : %s\n", k, flds[k])
	}

	return str
}
