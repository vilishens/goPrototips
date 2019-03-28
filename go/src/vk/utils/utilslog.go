package utils

import (
	"log"
	"os"
	vomni "vk/omnibus"
)

// String into a logger -- <PREFIX> <SEPARATOR> <DATE+TIME> <SEPARATOR> <STR>
func LogStr(d *log.Logger, str string) {
	strNew := vomni.UDPMessageSeparator + " " + str
	d.Println(strNew)
}

// Logger with a prefix -- <PREFIX> <SEPARATOR> <DATE+TIME>
func LogNew(d *os.File, prefix string) (newLog *log.Logger) {
	return log.New(d, prefix+vomni.UDPMessageSeparator+" ", vomni.LogLoggerFlags)
}

// Point logger with no prefix -- <SEPARATOR> <DATE+TIME>
func LogNewPoint(d *os.File) (newLog *log.Logger) {
	return log.New(d, vomni.UDPMessageSeparator+" ", vomni.LogLoggerFlags)
}
