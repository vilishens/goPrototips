package omnibus

const (
	MsgCdOutputHelloFromStation = 0x00000001 // Output <station name><msgCd><msgNbr><station UTC seconds><station time offset><stationIP><stationPort>
	MsgCdInputHelloFromPoint    = 0x00000002 // Input  <point name><msgCd><msgNbr><pointIP><pointPort>
	MsgCdOutputSetGpio          = 0x00000004 // Output <station name><msgCd><msgNbr><Gpio><set value>
)

const (
	UDPMessageSeparator = ":::"

	MsgIndexPrefixSender = 0
	MsgIndexPrefixCd     = 1
	MsgIndexPrefixNbr    = 2

	MsgPrefixLen = 3
)

// Hello From Station
const (
	MsgIndexHelloFromStationTime   = 0
	MsgIndexHelloFromStationOffset = 1
	MsgIndexHelloFromStationIP     = 2
	MsgIndexHelloFromStationPort   = 3

	MsgHelloFromStationLen = 4
)

// Hello From Point
const (
	MsgIndexHelloFromPointIP   = 0
	MsgIndexHelloFromPointPort = 1

	MsgHelloFromPointLen = 2
)

// Set Gpio
const (
	MsgIndexSetGpioGpio = 0
	MsgIndexSetGpioSet  = 1

	MsgSetGpioLen = 2
)
