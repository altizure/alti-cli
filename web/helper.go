package web

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/jackytck/alti-cli/errors"
	"github.com/jackytck/alti-cli/gql"
)

// StartLocalServer starts a local server serving dir on random port.
// If ip is not provided, non-local ip will be used.
// If port is not provided, a random port will be used.
func StartLocalServer(dir, ip, port string, verbose bool) (string, func(), error) {
	var address string

	// use prefered ip if not provided
	if ip == "" {
		pu, _, err := PreferedLocalURL(verbose)
		if err != nil {
			return "", nil, err
		}
		address = pu.Hostname() + ":" + port
	} else {
		address = ip + ":" + port
	}

	s := Server{Directory: dir, Address: address}
	hs, p, err := s.ServeStatic(verbose)
	if err != nil {
		return "", nil, err
	}
	// using random port
	ps := strconv.Itoa(p)
	if ps != port {
		address += ps
	}

	baseURL := fmt.Sprintf("http://%s", address)
	log.Printf("Serving files at %s\n", baseURL)
	done := func() {
		log.Println("Shutting down local server...")
		if err = hs.Shutdown(context.TODO()); err != nil {
			panic(err)
		}
	}
	return baseURL, done, nil
}

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
func CheckVisibility(verbose bool) (map[string]bool, error) {
	ret := make(map[string]bool)

	// tmp dir for server
	tmpDir, err := ioutil.TempDir("", "alti-cli-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	// create local web server
	s := Server{Directory: tmpDir}
	server, port, err := s.ServeStatic(false)
	if err != nil {
		return nil, err
	}

	// check each ip
	ips, err := GetAllIP()
	if err != nil {
		return nil, err
	}
	ch := make(chan netCheckResult)
	var wg sync.WaitGroup
	wg.Add(len(ips))

	go func() {
		wg.Wait()
		close(ch)
	}()
	for _, ip := range ips {
		go func(ip string) {
			url := fmt.Sprintf("http://%v:%v", ip, port)
			if verbose {
				log.Printf("Checking %q...", url)
			}
			res := gql.CheckDirectNetwork(url)
			ch <- netCheckResult{url, res}
			wg.Done()
		}(ip)
	}
	for r := range ch {
		ret[r.url] = r.visibility
	}

	// close down temp server
	if err := server.Shutdown(context.TODO()); err != nil {
		return nil, err
	}

	return ret, nil
}

type netCheckResult struct {
	url        string
	visibility bool
}

// PreferedLocalURL returns visible url in following preference:
// non-localhost > localhost url.
func PreferedLocalURL(verbose bool) (*url.URL, map[string]bool, error) {
	checks, err := CheckVisibility(verbose)
	if err != nil {
		return nil, nil, err
	}

	// classify local or non-local ips
	var local, nonLocal []string
	for k, v := range checks {
		if !v {
			continue
		}
		if strings.Contains(k, "127.0.0.1") {
			local = append(local, k)
		} else {
			nonLocal = append(nonLocal, k)
		}
	}
	if len(local)+len(nonLocal) == 0 {
		return nil, checks, errors.ErrClientInvisible
	}
	sort.Strings(nonLocal)

	// prefer non-local over local ip
	var ret string
	if len(nonLocal) > 0 {
		ret = nonLocal[0]
	} else {
		ret = local[0]
	}

	u, err := url.ParseRequestURI(ret)
	if err != nil {
		return nil, checks, err
	}
	return u, checks, nil
}

// CheckVisibilityIPPort checks if starting a local server over the given
// ip and port could be visible by the api server.
func CheckVisibilityIPPort(ip, port string, verbose bool) (bool, error) {
	url := fmt.Sprintf("http://%v:%v", ip, port)
	if verbose {
		log.Printf("Checking %q...", url)
	}

	// tmp dir for server
	tmpDir, err := ioutil.TempDir("", "alti-cli-")
	if err != nil {
		return false, err
	}
	defer os.RemoveAll(tmpDir)

	// create local web server
	s := Server{
		Directory: tmpDir,
		Address:   fmt.Sprintf("%s:%s", ip, port),
	}
	server, _, err := s.ServeStatic(false)
	if err != nil {
		return false, err
	}

	// check given ip + port over api server
	res := gql.CheckDirectNetwork(url)

	// close down temp server
	if err := server.Shutdown(context.TODO()); err != nil {
		return false, err
	}

	return res, nil
}
