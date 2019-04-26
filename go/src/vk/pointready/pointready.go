package pointready

func Prepare(chGoOn chan bool, chDone chan int, chErr chan error) {

	relayInterval()

	chGoOn <- true
	//	for {
	//		time.Sleep(vomni.DelayStepExec)
	//	}
}
