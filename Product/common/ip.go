package common

import (
	"errors"
	"net"
)

func GetIntranceIP() (string, error) {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addr {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("Fail to get intrance IP")
}
