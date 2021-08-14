package clash

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

var regionList map[string][]string

//var avgDelay int64

type Delay struct {
	History []struct {
		Delay int64  `json:"delay"`
		Time  string `json:"time"`
	} `json:"history"`
	Name string `json:"name"`
	Type string `json:"type"`
	Udp  bool   `json:"udp"`
}

type ProxiesResp struct {
	Proxies map[string]Delay `json:"proxies"`
}

var regionCode = map[string]string{
	"jp": "日本",
	"hk": "香港",
	"us": "美国",
}

// GetProxies
// request proxies and make delay map
func GetProxies() error {
	var proxyDelayList ProxyDelayList

	api := os.Getenv("CLASH_CONTROLLER")
	u := fmt.Sprintf("http://%s/proxies", api)
	bodyString, err := HttpGet(u)
	if err != nil {
		return err
	}
	var resp ProxiesResp
	err = json.Unmarshal([]byte(bodyString), &resp)
	if err != nil {
		return err
	}

	for name, delay := range resp.Proxies {
		proxyDelay := ProxyDelay{
			Name:  name,
			Delay: math.MaxInt64,
		}
		for _, history := range delay.History {
			if history.Delay > 0 {
				proxyDelay.Delay = history.Delay
			}
		}
		proxyDelayList = append(proxyDelayList, proxyDelay)
	}

	sort.Sort(proxyDelayList)

	regionList = make(map[string][]string)
	for _, proxyDelay := range proxyDelayList {
		for code, region := range regionCode {
			if len(regionList[code]) < 10 && InRegion(proxyDelay.Name, code, region) {
				regionList[code] = append(regionList[code], proxyDelay.Name)
			}
		}
	}
	return nil
}

func InRegion(name, code, region string) bool {
	if strings.Contains(name, region) {
		return true
	}
	if strings.Contains(name, code) {
		return true
	}
	if strings.Contains(name, strings.ToUpper(code)) {
		return true
	}
	return false
}

type ProxyDelay struct {
	Name  string
	Delay int64
}

type ProxyDelayList []ProxyDelay

func (p ProxyDelayList) Len() int           { return len(p) }
func (p ProxyDelayList) Less(i, j int) bool { return p[i].Delay < p[j].Delay }
func (p ProxyDelayList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
