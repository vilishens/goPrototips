package runrelayinterval

import "fmt"

var RunningPoints map[string]RunData

func init() {
	RunningPoints = make(map[string]RunData)
}

func (d RunData) Starter(chGoOn chan bool, chErr chan error) {
	fmt.Println("============ XXXXX ====================")
}
