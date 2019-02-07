package web

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/jackytck/alti-cli/gql"
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

// CheckVisibility checks if api server could reach this client machine
// for each newtwork interface.
func CheckVisibility() (map[string]bool, error) {
	ret := make(map[string]bool)

	// create local web server
	s := Server{Directory: "/tmp"}
	server, port, err := s.ServeStatic(false)
	if err != nil {
		return nil, err
	}

	// check each ip
	ips, err := GetAllIP()
	if err != nil {
		return nil, err
	}
	for _, ip := range ips {
		url := fmt.Sprintf("http://%v:%v", ip, port)
		res := gql.CheckDirectNetwork(url)
		ret[url] = res
	}

	// close down temp server
	if err := server.Shutdown(context.TODO()); err != nil {
		return nil, err
	}

	return ret, nil
}
