package omnibus

import (
	"fmt"
)

func init() {
	stepList = make(map[string]bool)
}

func StepInList(step string) {
	stepList[step] = true
}

func StepRemoveFromList(step string) {
	delete(stepList, step)
}

func AreStepsInList(steps []string) (err error) {
	for _, v := range steps {
		if _, ok := stepList[v]; !ok {
			err = fmt.Errorf("%q is not in the running step list", v)
			break
		}
	}

	return
}

func StepCount() (count int) {
	return len(stepList)
}
