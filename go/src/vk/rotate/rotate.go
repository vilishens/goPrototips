package rotate

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"time"
	vomni "vk/omnibus"
	vparams "vk/params"
	vutils "vk/utils"
)

var allLoggers []activeLog

func init() {
	allLoggers = []activeLog{}
}

func addMainLogs() {

	mainList := []*log.Logger{vomni.LogData, vomni.LogErr, vomni.LogFatal, vomni.LogInfo}
	allLoggers = append(allLoggers, activeLog{path: vparams.Params.LogMainPath, file: vomni.LogMainFile, loggers: mainList})
}

func MainStart(chGoOn chan bool, chDone chan int, chErr chan error) {

	fmt.Println("ALEX PASHUTIN")

	addMainLogs()

	// start the brand new rotation with the Main Log file
	err := SetRotateCfg(vparams.Params.LogMainPath, vparams.Params.RotateMainCfg, vparams.Params.RotateRunCfg, true)
	if nil != err {
		err = vutils.ErrFuncLine(err)
		vutils.LogStr(vomni.LogErr, err.Error())
		chErr <- err
		return
	}

	chGoOn <- true

	locDone := make(chan int)
	locErr := make(chan error)
	go runRotate(locDone, locErr)

	select {
	case err = <-locErr:
		err = vutils.ErrFuncLine(err)
		vutils.LogStr(vomni.LogErr, err.Error())
		chErr <- err
		return
	case <-locDone:
		chDone <- vomni.DoneStop
		return
	}
}

func SetRotateCfg(file2Rotate string, cfgTmpl2Use string, cfg2RunFile string, newRotation bool) (err error) {

	usr := new(user.User)

	/*
	 * find the current user data (user name and user group)
	 * to create this particular rotation configuration
	 */
	usr, err = user.Current()
	if nil != err {
		return vutils.ErrFuncLine(err)
	}

	name := usr.Username

	group := ""
	if usrGrp, err := user.LookupGroupId(usr.Gid); nil != err {
		return vutils.ErrFuncLine(err)
	} else {
		group = usrGrp.Name
	}

	/*
	 * Prepare rotation configuration
	 * from the configuration template
	 */
	var format []byte
	// read necessary data file rotation configuration template
	if format, err = ioutil.ReadFile(cfgTmpl2Use); nil != err {
		return vutils.ErrFuncLine(err)
	}

	// put required data (user name and user group) into configuration template
	str := fmt.Sprintf(string(format), file2Rotate, name, group)

	/*
	 * delete the existing running configuration
	 * if it is required to start with the brand new configuration file
	 */
	if newRotation {
		has, err := vutils.PathExists(cfg2RunFile)
		if nil != err {
			return vutils.ErrFuncLine(err)
		}
		if has {
			if err = vutils.FileDelete(cfg2RunFile); nil != err {
				return vutils.ErrFuncLine(err)
			}
		}
	}

	/*
	 * use configuration prepared from the template
	 * to put it in the running configuration
	 * or to add (in case the rotation was started - newRotation is false)
	 */
	if err = vutils.FileAppend(cfg2RunFile, str); nil != err {
		return vutils.ErrFuncLine(err)
	}

	return
}

func runRotate(chDone chan int, chErr chan error) {

	if 0 >= vparams.Params.RotateRunSecs {
		chErr <- vutils.ErrFuncLine(fmt.Errorf("\nWrong point rotation interval %d", vparams.Params.RotateRunSecs))
		return
	}

	for {
		timer := time.NewTimer(time.Duration(vparams.Params.RotateRunSecs) * time.Second)

		// rotate
		if err := setRotateFiles(); nil != err {
			chErr <- vutils.ErrFuncLine(fmt.Errorf("\nRotation command failure -- %v", err))
			return
		}

		vutils.LogStr(vomni.LogInfo, "Rotate check")

		timeStr := time.Now().Format("2006-01-02 15:04:05 -07:00 MST")
		str := fmt.Sprintf("==>>>>>\n==>>>>>\n==>>>>>\n==>>>>>\n %s <<<<<< ROTATE\n==>>>>>\n==>>>>>\n==>>>>>", timeStr)
		fmt.Println(str)

		select {
		case <-timer.C:

			fmt.Println("\n\n\n$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
			fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
			fmt.Printf("$$$$$$$$$$$$$$$$$$$$$$ %q $$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$\n", timeStr)
			fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$")
			fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$\n\n\n")
		}
	}
}

func setRotateFiles() (err error) {
	// rotate files if necessary
	if err = runRotateCmd(); nil != err {
		err = vutils.ErrFuncLine(fmt.Errorf("\nRotation command failure -- %v", err))
		vutils.LogStr(vomni.LogErr, err.Error())
		return
	}

	jaturpina ar Reaassignu


	/*

		// reassign the main logger files
		if err = reassingMainFile(); nil != err {
			vomni.RootErr <- vutils.ErrFuncLine(fmt.Errorf("\nRotation main file reassign failure -- %v", err))
			return
		}

		// reassign point logger files
		if err = reassignPointFiles(); nil != err {
			vomni.RootErr <- vutils.ErrFuncLine(fmt.Errorf("\nRotation point files reassign failure -- %v", err))
			return
		}

		fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
		fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
		fmt.Println("@@@@@@@@@@@@@@@@@@@@@@ ZAPAH-ZAPAH-ZAPAH --> RESTART @@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
		fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
		fmt.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	*/
	return
}

func reassingMainFile() (err error) {

	/*
		if vomni.LogMainFile, err = vutils.LogReAssign(vomni.LogMainFile, vomni.LogMainPath); nil != err {
			vomni.RootErr <- vutils.ErrFuncLine(fmt.Errorf("\nRotation file reaasign failure -- %v", err))
			return
		}

		vomni.LogData.SetOutput(vomni.LogMainFile)
		vomni.LogErr.SetOutput(vomni.LogMainFile)
		vomni.LogFatal.SetOutput(vomni.LogMainFile)
		vomni.LogInfo.SetOutput(vomni.LogMainFile)
	*/
	return
}

/*
func reassignPointFiles() (err error) {

	if err = xrun.RotateReAssign(); nil != err {
		return err
	}

	return
}
*/

func LogReAssign(f *os.File, path string) (fNew *os.File, err error) {
	var perms os.FileMode
	flags := vomni.LogFileFlags

	if nil != f {
		var stat os.FileInfo

		stat, err = f.Stat()
		if nil != err {
			return nil, vutils.ErrFuncLine(fmt.Errorf("Could not get stat of the file %s", path))
		}

		perms = stat.Mode()

		if err = f.Close(); nil != err {
			return nil, vutils.ErrFuncLine(fmt.Errorf("Could not close the file %s", path))
		}
	} else {
		perms = vomni.LogUserPerms
	}

	return vutils.OpenFile(path, flags, perms)
}

func runRotateCmd() (err error) {
	//	find the local status file
	dirpath := filepath.Dir(vparams.Params.RotateRunCfg)
	statusF := filepath.Join(dirpath, vparams.Params.RotateStatusFileName)

	if has, _ := vutils.PathExists(statusF); has {
		err = vutils.FileDelete(statusF)
		if nil != err {
			return vutils.ErrFuncLine(err)
		}
	}

	// logrotate <conf.file> -s <localstatus.file>
	cmd := exec.Command("logrotate", vparams.Params.RotateRunCfg, "-s", statusF)

	if err = cmd.Run(); nil != err {
		return vutils.ErrFuncLine(err)
	}

	b, _ := cmd.Output()
	fmt.Println("Execution of my mind >>>", string(b))

	return
}
