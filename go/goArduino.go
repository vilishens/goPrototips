package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	vomni "vk/omnibus"
	sall "vk/steps/allsteps"
	vutils "vk/utils"
)

func init() {
	root()
}

func main() {

	end := false
	endCd := -1

	for !end {
		if 0 > endCd {
			vutils.LogStr(vomni.LogInfo, "***** App - START *****")
		}

		endCd := runApp()

		switch endCd {
		case vomni.DoneRestart, vomni.DoneStop, vomni.DoneError, vomni.DoneReboot:
			end = true
		}

		switch endCd {
		case vomni.DoneRestart:
			vutils.LogStr(vomni.LogInfo, "***** App - RESTART *****")
		case vomni.DoneStop:
			vutils.LogStr(vomni.LogInfo, "***** App - STOP *****")
		case vomni.DoneError:
			str := fmt.Sprintf("***** App - ERROR *****")
			vutils.LogStr(vomni.LogInfo, str)
		case vomni.DoneReboot:
			vutils.LogStr(vomni.LogInfo, "***** App - REBOOT *****")
		default:
			str := fmt.Sprintf("***** App - unknown Exit code %d *****", endCd)
			vutils.LogStr(vomni.LogInfo, str)
		}

		if end {
			os.Exit(endCd)
		}
	}
}

func runApp() (cd int) {
	chDone := make(chan int)

	go sall.DoSteps(chDone)

	select {
	case err := <-vomni.RootErr:
		fmt.Printf("App finished due to an error ---> %v\n", err)
		cd = vomni.DoneError
		break
	case cd = <-chDone:
		str := fmt.Sprintf("***** App - received code %d *****", cd)
		vutils.LogStr(vomni.LogInfo, str)
		break
	}

	return
}

func root() {
	rootPath()
	rootLog()
}

func rootPath() {
	// It is necessary to keep the root caller path to create
	// correct file paths further
	if _, rootFile, _, ok := runtime.Caller(0); !ok {
		err := fmt.Errorf("Could not get Root Path")
		log.Fatal(err)
	} else {
		vomni.RootPath = filepath.Dir(rootFile)
	}

	return
}

func rootLog() {
	var err error

	path := filepath.Join(vomni.RootPath, vomni.LogMainPath)

	vomni.LogMainFile, err = vutils.OpenFile(path, vomni.LogFileFlags, vomni.LogUserPerms)
	if nil != err {
		log.Fatal(fmt.Errorf("Could not open the main log file --- %v", err))
	}

	vomni.LogData = vutils.LogNew(vomni.LogMainFile, vomni.LogPrefixData)
	vomni.LogErr = vutils.LogNew(vomni.LogMainFile, vomni.LogPrefixErr)
	vomni.LogFatal = vutils.LogNew(vomni.LogMainFile, vomni.LogPrefixFatal)
	vomni.LogInfo = vutils.LogNew(vomni.LogMainFile, vomni.LogPrefixInfo)
}
