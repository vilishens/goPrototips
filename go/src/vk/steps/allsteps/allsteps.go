package allsteps

import (
	"fmt"

	//	vcfg "vk/cfg"
	vomni "vk/omnibus"
	vutils "vk/utils"

	vstep "vk/steps/step"

	//	schecknet "vk/steps/stepchecknet"
	scfg "vk/steps/stepconfig"
	sparams "vk/steps/stepparams"

	//	spointcfg "vk/steps/steppointconfig"
	//	spointscan "vk/steps/steppointscan"
	//	srunpoints "vk/steps/steprunpoints"
	sstart "vk/steps/stepstart"
	//	sudp "vk/steps/stepudp"
	sweb "vk/steps/stepweb"
)

var steps = make(map[string]vstep.Step)
var stepSequence []string

func init() {
	initSteps()
}

func initSteps() {
	addStep(&(sstart.ThisStep))
	addStep(&(scfg.ThisStep))
	addStep(&(sparams.ThisStep))
	addStep(&(sweb.ThisStep))
	// seit jaieliek rotateMain solis
	//	addStep(&(schecknet.ThisStep))
	//	addStep(&(sweb.ThisStep)) // WEB step must be before point steps
	//	addStep(&(spointcfg.ThisStep))
	//	addStep(&(srunpoints.ThisStep))
	//	addStep(&(sudp.ThisStep))
	//	addStep(&(spointscan.ThisStep))
}

func addStep(s vstep.Step) {
	sName := s.StepName()

	if _, exists := steps[sName]; exists {
		panic(fmt.Sprintf("ALERT! Step '%s' used more than once (ind %d)", sName, len(stepSequence)))
	}

	stepSequence = append(stepSequence, sName)
	steps[sName] = s
}

func DoSteps(chDone chan int) {

	locDone := make(chan int)

	go doAllSteps(locDone)

	done := <-locDone

	chDone <- done
}

func doAllSteps(chanDone chan int) {

	chErr := make(chan error)
	chDone := make(chan int)
	chGoOn := make(chan bool)

	count := -1
	stop := false
	err := error(nil)
	done := 0

	ind := 0

	for _, s := range stepSequence {
		thisS := steps[s]

		ind++

		str := fmt.Sprintf("===== Step %q -> started", thisS.StepName())
		fmt.Println(str)
		vutils.LogStr(vomni.LogInfo, fmt.Sprintf(str))
		go vstep.Execute(thisS, chDone, chGoOn, chErr)

		select {
		case <-chGoOn:
			str = fmt.Sprintf("===== Step %q -> sent GoOn", thisS.StepName())
			vutils.LogStr(vomni.LogInfo, str)
		case err = <-chErr:
			stop = true
		case done = <-chDone:
			stop = true
		}

		if stop {
			break
		}

		count++
	}

	fmt.Println("ooooooooooooooooooooooooooooooooooooooooooooo END-APP")
	fmt.Println("ooooooooooooooooooooooooooooooooooooooooooooo END-APP")
	fmt.Println("ooooooooooooooooooooooooooooooooooooooooooooo END-APP")

	if !stop {

		str := fmt.Sprintf("===== All steps are running")
		vutils.LogStr(vomni.LogInfo, str)
		fmt.Println(str)

		select {
		case err = <-chErr:
			str := fmt.Sprintf("Steps need to be closed due to an err - %q", err)
			vutils.LogStr(vomni.LogErr, str)
			stop = true
		case done = <-vomni.RootDone:
			stop = true
		}
	}

	fmt.Println("Tagad jābeidz...")

	for ; count >= 0; count-- {
		// let's do Post of each step starting from the last one
		thisS := steps[stepSequence[count]]
		thisS.StepPost()

		str := fmt.Sprintf("===== Step %q -> closed", thisS.StepName())
		vutils.LogStr(vomni.LogInfo, str)
		fmt.Println(str)
	}

	if stop {
		if nil != err {
			vomni.RootErr <- err
		}
		if 0 != done {
			chanDone <- done
		}
		return
	}
}
