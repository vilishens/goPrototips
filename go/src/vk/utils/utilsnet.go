package utils

import (
	"net"
	vomni "vk/omnibus"
)

// LocalIP returns the non loopback local IP of the host
func InternalIPv4() (ip string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		err = ErrFuncLine(err)
		vomni.LogErr.Println(err)
		return
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ip = ipnet.IP.String()
				return
			}
		}
	}
	return
}
