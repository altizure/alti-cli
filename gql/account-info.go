package gql

import (
	"errors"
	"sync"
	"time"
)

// GetAccountInfoTimeout gets the account info with timeout.
func GetAccountInfoTimeout(endpoint, key, token string, timeout time.Duration) (AccountInfo, error) {
	c := make(chan AccountInfo)
	go func() {
		defer close(c)
		c <- GetAccountInfo(endpoint, key, token)
	}()

	var ret AccountInfo
	select {
	case ret = <-c:
		return ret, nil
	case <-time.After(timeout):
		return ret, errors.New("timeout")
	}
}

// GetAccountInfo gets the account info.
func GetAccountInfo(endpoint, key, token string) AccountInfo {
	var ret AccountInfo

	var wg sync.WaitGroup
	wg.Add(6)

	go func() {
		ret.Sales = IsSales(endpoint, key, token)
		wg.Done()
	}()
	go func() {
		ret.Super = IsSuper(endpoint, key, token)
		wg.Done()
	}()
	go func() {
		ret.ImageCloud = SupportedCloud(endpoint, key, "image")
		wg.Done()
	}()
	go func() {
		ret.ModelCloud = SupportedCloud(endpoint, key, "model")
		wg.Done()
	}()
	go func() {
		ret.MetaCloud = SupportedCloud(endpoint, key, "meta")
		wg.Done()
	}()
	go func() {
		ret.Version, ret.ResponseTime = Version(endpoint, key)
		wg.Done()
	}()

	wg.Wait()
	return ret
}

// AccountInfo represents the account info.
type AccountInfo struct {
	Super        bool
	Sales        bool
	ImageCloud   []string
	ModelCloud   []string
	MetaCloud    []string
	Version      string
	ResponseTime time.Duration
}
