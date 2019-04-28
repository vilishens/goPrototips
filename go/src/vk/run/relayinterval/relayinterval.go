package runrelayinterval

var RunningPoints map[string]RunData

func init() {
	RunningPoints = make(map[string]RunData)
}

func (d RunData) Starter(chGoOn chan bool, chErr chan error) {

}
