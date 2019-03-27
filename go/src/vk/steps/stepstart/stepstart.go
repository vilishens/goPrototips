package stepstart

import (
	"fmt"
	"time"
	vomni "vk/omnibus"
	"vk/steps/step"
)

type thisStep step.StepVars

var ThisStep thisStep

func init() {
	ThisStep.Name = vomni.StepNameStart
	ThisStep.Err = make(chan error)
	ThisStep.GoOn = make(chan bool)
	ThisStep.Done = make(chan int)
}

func (s *thisStep) stepDo() {
	// put code here to do what is necessary

	//	err := vcli.Init()
	err := error(nil)
	if nil != err {
		s.Err <- err
		return
	}

	fmt.Println("BACK goON")
	time.Sleep(5 * time.Second)

	s.GoOn <- true
}

func (s *thisStep) StepName() string {
	return ThisStep.Name
}

func (s *thisStep) StepPost() {
	s.Done <- vomni.DoneStop
	return
}

func (s *thisStep) StepPre(chDone chan int, chGoOn chan bool, chErr chan error) {
	chGoOn <- true
}

func (s *thisStep) StepExec(chDone chan int, chGoOn chan bool, chErr chan error) {

	go s.stepDo()

	stop := false
	for !stop {
		select {
		case locErr := <-s.Err:
			chErr <- locErr
			stop = true
		case locDone := <-s.Done:
			chDone <- locDone
			stop = true
		case locGoOn := <-s.GoOn:
			chGoOn <- locGoOn
		}
	}
}
