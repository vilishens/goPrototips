package messages

import (
	"net"
	"time"
)

type SendMsg struct {
	UDPAddr net.UDPAddr
	MsgNbr  int // if negative, means nothing to send, it must be deleted
	Repeat  int
	Msg     string
	Last    time.Time
}

type SendMsgArray []SendMsg
