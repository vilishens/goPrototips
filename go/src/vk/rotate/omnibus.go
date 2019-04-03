package rotate

import (
	"log"
	"os"
)

type activeLog struct {
	path    string
	file    *os.File
	loggers []*log.Logger
}
