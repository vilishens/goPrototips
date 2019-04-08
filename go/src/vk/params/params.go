package params

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"sync"
	vconf "vk/cfg"
	vutils "vk/utils"
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

	Params.NetRequirement = -1
	Params.IPExternalAddressCmds = []string{} // commands to find the station external IP address

	/*


		Params.UDPPort = -14
		Params.PointLogPath = "../data"

		Params.RotateMainCfg = "" // jāliek vērtības iz cfg faila
		Params.RotatePointCfg = ""
		Params.RotateRunCfg = ""
		Params.RotateRunSecs = 0

		Params.Name = "Don't know yet"
		Params.InternalPort = 7520
		Params.InternalIPv4 = "127.0.0.1"
		Params.ExternalPort = 8025
		Params.ExternalIPv4 = "Don't know yet"
		Params.WebEmail = "arduins@gmx.com"
		Params.WebAliveInterval = 7200
		Params.WebEmailMutt = "../cfg/factory/muttrc.set"
		Params.ScriptPath = "../cfg/factory/scripts"
		Params.LogPath = "../data"
		Params.TemplatePath = "tmpl"
		Params.TemplateExt = ".tmpl"
		Params.WebPort = 7015
		Params.EventPath = "../data/event"
		Params.ErrorPath = "../data/error"

		// JĀIZDOMĀ, KO DaRīT ar modēm
		Params.PointModeFiles = make(map[string]string)
		Params.PointModeFiles[vomni.PointModeIntervalOnOff] = "../cfg/factory/modes/intervalOnOff.js"

		Params.DevModes = []string{vomni.PointModeIntervalOnOff}
	*/
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

	if 0 <= data.NetRequirement {
		Params.NetRequirement = data.NetRequirement
	}
	if (nil == err) && (0 < len(data.IPExternalAddressCmds)) {
		Params.IPExternalAddressCmds = make([]string, len(data.IPExternalAddressCmds))
		copy(Params.IPExternalAddressCmds, data.IPExternalAddressCmds)
	}

	/*
		if "" != data.PointLogPath {
			Params.PointLogPath = data.PointLogPath
		}



		Params.PointLogPath = vutils.FileAbsPath(data.PointLogPath, "")

		// Jāliek pārbaude vai konfiga dati nav tukši vai negatīvi
		//
		if "" != data.Name {
			Params.Name = data.Name
		}
		if 0 <= data.InternalPort {
			Params.InternalPort = data.InternalPort
		}
		if "" != data.InternalIP {
			Params.InternalIPv4 = data.InternalIP
		}
		if 0 <= data.ExternalPort {
			Params.ExternalPort = data.ExternalPort
		}
		if "" != data.WebEmail {
			Params.WebEmail = data.WebEmail
		}
		if 0 <= data.WebAliveInterval {
			Params.WebAliveInterval = data.WebAliveInterval
		}
		if "" != data.WebEmailMutt {
			Params.WebEmailMutt = data.WebEmailMutt
		}
		if "" != data.ScriptPath {
			Params.ScriptPath = data.ScriptPath
		}
		Params.ScriptPath = vutils.FileAbsPath(Params.ScriptPath, "")
		if "" != data.TemplatePath {
			Params.TemplatePath = data.TemplatePath
		}
		Params.TemplatePath = vutils.FileAbsPath(data.TemplatePath, "")
		if "" != data.TemplateExt {
			Params.TemplateExt = data.TemplateExt
		}
		if 0 <= data.WebPort {
			Params.WebPort = data.WebPort
		}
		if 0 <= data.UDPPort {
			Params.UDPPort = data.UDPPort
		}
		if "" != data.EventPath {
			Params.EventPath = data.EventPath
		}
		Params.EventPath = vutils.FileAbsPath(data.EventPath, "")

		if "" != data.ErrorPath {
			Params.ErrorPath = data.ErrorPath
		}
		Params.ErrorPath = vutils.FileAbsPath(data.ErrorPath, "")

		if "" != data.RotateMainCfg {
			Params.RotateMainCfg = data.RotateMainCfg
		}
		Params.RotateMainCfg = vutils.FileAbsPath(data.RotateMainCfg, "")
		if "" != data.RotatePointCfg {
			Params.RotatePointCfg = data.RotatePointCfg
		}
		Params.RotatePointCfg = vutils.FileAbsPath(data.RotatePointCfg, "")
		if "" != data.RotateRunCfg {
			Params.RotateRunCfg = data.RotateRunCfg
		}
		Params.RotateRunCfg = vutils.FileAbsPath(data.RotateRunCfg, "")

		if 0 < data.RotateRunSecs {
			Params.RotateRunSecs = data.RotateRunSecs
		}

		fmt.Printf("<<<<<<<<<<<<<<<<<<<<<<<<<<<< ERRORe === \"%s\"\n\tSUDAK \"%s\"\n", Params.ErrorPath, data.ErrorPath)

		//???
		//???
		//???
		//???

		path := vutils.FileAbsPath(data.WebEmailMutt, "")

		fmt.Printf("SI === %s\n", path)

		ok := false
		ok, err := vutils.PathExists(path)

		if !ok {
			chErr <- vutils.ErrFuncLine(fmt.Errorf("Mutt file \"%s\" doesn't exist", path))
		} else if nil == err {
			Params.WebEmailMutt = path
		}

		path = vutils.FileAbsPath(data.ScriptPath, "")

		if _, errScript := os.Stat(path); errScript == nil {
			// path/to/whatever exists
			Params.ScriptPath = path
		}

		if "" != Params.ScriptPath {
			err = handleScripts()
		}

		if "" != data.LogPath {
			Params.LogPath = data.LogPath
		}

		fmt.Printf("\n\n\n\n\nBORTICH === %s (%s)\n\n\n\n", data.LogPath, Params.LogPath)

		fmt.Printf("MUTT FILE %s --- OK %v -- Err %v -- PATH %s\n", Params.WebEmailMutt, ok, err, path)

		if nil == err {

			fmt.Println("Zemieki")

			err = setNetAddrs()

			fmt.Println("Fetisoff")
		}

	*/

	if nil != err {
		chErr <- err
	} else {
		chDone <- true
	}
}

//######################################################################################

func Scripts(data vconf.CfgFinalData, wg *sync.WaitGroup, err *error) {

	defer wg.Done()
}

func PutXXX(data vconf.CfgFinalData, chDone chan bool, chErr chan error) {

	//
	/*
		--	Name             string
		--	InternalPort     int
		--	InternalIP       string
		--	ExternalPort     int
		--	ExternalIP       string
		--	WebEmail         string
		--	WebAliveInterval int
		--	WebEmailMutt     string
		--	ScriptPath       string
			LogPath          string
		???	PointModeFiles   map[string]string
		--	TemplatePath     string
		--	TemplateExt      string
		???	DevModes         []string
		--	WebPort          int
		--	EventPath        string
	*/

	if "" != data.PointLogPath {
		Params.PointLogPath = data.PointLogPath
	}
	Params.PointLogPath = vutils.FileAbsPath(data.PointLogPath, "")

	// Jāliek pārbaude vai konfiga dati nav tukši vai negatīvi
	//
	if "" != data.Name {
		Params.Name = data.Name
	}
	if 0 <= data.InternalPort {
		Params.InternalPort = data.InternalPort
	}
	if "" != data.InternalIP {
		Params.InternalIPv4 = data.InternalIP
	}
	if 0 <= data.ExternalPort {
		Params.ExternalPort = data.ExternalPort
	}
	if "" != data.WebEmail {
		Params.WebEmail = data.WebEmail
	}
	if 0 <= data.WebAliveInterval {
		Params.WebAliveInterval = data.WebAliveInterval
	}
	if "" != data.WebEmailMutt {
		Params.WebEmailMutt = data.WebEmailMutt
	}
	if "" != data.ScriptPath {
		Params.ScriptPath = data.ScriptPath
	}
	Params.ScriptPath = vutils.FileAbsPath(Params.ScriptPath, "")
	if "" != data.TemplatePath {
		Params.TemplatePath = data.TemplatePath
	}
	Params.TemplatePath = vutils.FileAbsPath(data.TemplatePath, "")
	if "" != data.TemplateExt {
		Params.TemplateExt = data.TemplateExt
	}
	if 0 <= data.WebPort {
		Params.WebPort = data.WebPort
	}
	if 0 <= data.UDPPort {
		Params.UDPPort = data.UDPPort
	}
	if "" != data.EventPath {
		Params.EventPath = data.EventPath
	}
	Params.EventPath = vutils.FileAbsPath(data.EventPath, "")

	if "" != data.ErrorPath {
		Params.ErrorPath = data.ErrorPath
	}
	Params.ErrorPath = vutils.FileAbsPath(data.ErrorPath, "")

	if "" != data.RotateMainCfg {
		Params.RotateMainCfg = data.RotateMainCfg
	}
	Params.RotateMainCfg = vutils.FileAbsPath(data.RotateMainCfg, "")
	if "" != data.RotatePointCfg {
		Params.RotatePointCfg = data.RotatePointCfg
	}
	Params.RotatePointCfg = vutils.FileAbsPath(data.RotatePointCfg, "")
	if "" != data.RotateRunCfg {
		Params.RotateRunCfg = data.RotateRunCfg
	}
	Params.RotateRunCfg = vutils.FileAbsPath(data.RotateRunCfg, "")

	if 0 < data.RotateRunSecs {
		Params.RotateRunSecs = data.RotateRunSecs
	}

	fmt.Printf("<<<<<<<<<<<<<<<<<<<<<<<<<<<< ERRORe === \"%s\"\n\tSUDAK \"%s\"\n", Params.ErrorPath, data.ErrorPath)

	//???
	//???
	//???
	//???

	path := vutils.FileAbsPath(data.WebEmailMutt, "")

	fmt.Printf("SI === %s\n", path)

	ok := false
	ok, err := vutils.PathExists(path)

	if !ok {
		chErr <- vutils.ErrFuncLine(fmt.Errorf("Mutt file \"%s\" doesn't exist", path))
	} else if nil == err {
		Params.WebEmailMutt = path
	}

	path = vutils.FileAbsPath(data.ScriptPath, "")

	if _, errScript := os.Stat(path); errScript == nil {
		// path/to/whatever exists
		Params.ScriptPath = path
	}

	if "" != Params.ScriptPath {
		err = handleScripts()
	}

	if "" != data.LogPath {
		Params.LogPath = data.LogPath
	}

	fmt.Printf("\n\n\n\n\nBORTICH === %s (%s)\n\n\n\n", data.LogPath, Params.LogPath)

	fmt.Printf("MUTT FILE %s --- OK %v -- Err %v -- PATH %s\n", Params.WebEmailMutt, ok, err, path)

	if nil == err {

		fmt.Println("Zemieki")

		err = setNetAddrs()

		fmt.Println("Fetisoff")
	}

	if nil != err {
		chErr <- err
	} else {
		chDone <- true
	}

}

func handleScripts() (err error) {
	fmt.Println("### IGOR BOTVIN ###", Params.ScriptPath)

	if err = filepath.Walk(Params.ScriptPath, putScript); nil != err {
		err = vutils.ErrFuncLine(err)
	}

	return
}

func putScript(path string, info os.FileInfo, errX error) (err error) {
	if info.IsDir() {
		fmt.Printf("skipping a dir without errors: %+v \n", info.Name())
		return
	}

	err = copyScript(info.Name())

	return
}

func copyScript(file string) (err error) {

	usr := new(user.User)

	usr, err = user.Current()
	if nil != err {
		return
	}

	dst := vutils.FileAbsPath(usr.HomeDir, "bin")

	if _, err = os.Stat(dst); os.IsNotExist(err) {
		err = os.MkdirAll(dst, 0755)

		fmt.Println("Marina Fomina ", dst, " file ", file)

		if err != nil {
			return
		}
	} else {
		fmt.Println("Marina Lebedeva ", dst, " file ", file)
	}

	dst = vutils.FileAbsPath(dst, file)
	src := vutils.FileAbsPath(Params.ScriptPath, file)

	var dstf *os.File
	var srcf *os.File

	if srcf, err = os.Open(src); err != nil {
		return err
	}
	defer srcf.Close()

	if dstf, err = os.Create(dst); err != nil {
		return err
	}
	defer dstf.Close()

	if _, err = io.Copy(dstf, srcf); err != nil {
		return err
	}

	err = os.Chmod(dst, 0700)

	//	if srcinfo, err = os.Stat(src); err != nil {
	return err

	//dst := ""
	return
}

func setNetAddrs() (err error) {

	fmt.Println("2018.līgo")

	// if Params.InternalIPv4, err = vutils.InternalIPv4(); nil == err {
	//	//Params.ExternalIPv4, err = vutils.ExternalIPv4()
	// }

	fmt.Println("Maradona")

	if nil != err {
		err = vutils.ErrFuncLine(err)
	}

	return
}
