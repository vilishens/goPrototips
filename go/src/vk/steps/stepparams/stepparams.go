package stepparams

import (
	"time"
	vomni "vk/omnibus"
	vparam "vk/params"
	"vk/steps/step"
)

type thisStep step.StepVars

var ThisStep thisStep

func init() {
	ThisStep.Name = vomni.StepNameParams
	ThisStep.Err = make(chan error)
	ThisStep.GoOn = make(chan bool)
	ThisStep.Done = make(chan int)
}

func (s *thisStep) stepDo() {
	// put code here to do what is necessary

	chErr := make(chan error)
	chDone := make(chan bool)

	go vparam.Put(chDone, chErr)

	for {
		select {
		case err := <-chErr:
			s.Err <- err
			return
		case done := <-chDone:
			s.GoOn <- done
		}
	}
}

func (s *thisStep) StepName() string {
	return ThisStep.Name
}

func (s *thisStep) StepPre(chDone chan int, chGoOn chan bool, chErr chan error) {
	// do if something is required before the step execution

	chGoOn <- true
}

func (s *thisStep) StepExec(chDone chan int, chGoOn chan bool, chErr chan error) {

	// do what you would like
	go s.stepDo()

	stop := false
	for !stop {
		select {
		case locErr := <-s.Err:
			vomni.StepErr <- locErr
			stop = true
		case locDone := <-s.Done:
			chDone <- locDone
			stop = true
		case locGoOn := <-s.GoOn:
			chGoOn <- locGoOn
		}
		time.Sleep(vomni.StepExecDelay)
	}
}

func (s *thisStep) StepPost() {
	// may be something needs to be done before leave the step
	s.Done <- vomni.DoneStop
}
