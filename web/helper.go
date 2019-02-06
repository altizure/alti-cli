package web

import (
	"net"
	"strings"
)

// GetOutboundIP gets the preferred outbound ip of this machine.
func GetOutboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}

// GetAllIP gets all the local ip.
func GetAllIP() ([]string, error) {
	var ret []string
	ifaces, err := net.Interfaces()
	if err != nil {
		return ret, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return ret, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			ipStr := ip.String()
			if !strings.Contains(ipStr, ":") {
				ret = append(ret, ip.String())
			}
		}
	}
	return ret, nil
}
