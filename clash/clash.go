package clash

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"

	"sub2clash/log"
)

type Clash struct {
	Port               int                      `yaml:"port"`
	SocksPort          int                      `yaml:"socks-port"`
	AllowLan           bool                     `yaml:"allow-lan"`
	Mode               string                   `yaml:"mode"`
	LogLevel           string                   `yaml:"log-level"`
	ExternalController string                   `yaml:"external-controller"`
	ExternalUI         string                   `yaml:"external-ui"`
	Profile            map[string]interface{}   `yaml:"profile"`
	Dns                map[string]interface{}   `yaml:"dns"`
	Proxy              []map[string]interface{} `yaml:"proxies"`
	ProxyGroup         []map[string]interface{} `yaml:"proxy-groups"`
	Rule               []string                 `yaml:"rules"`

	// 兼容
	// ProxyOld      []map[string]interface{} `yaml:"Proxy"`
	// ProxyGroupOld []map[string]interface{} `yaml:"Proxy Group"`
	// RuleOld       []string                 `yaml:"Rule"`
}

func (c *Clash) LoadTemplate(proxies []interface{}) ([]byte, error) {
	tplFile := "/configs/base.yaml"
	_, err := os.Stat(tplFile)
	if err != nil && os.IsNotExist(err) {
		return nil, fmt.Errorf("[%s] template doesn't exist", tplFile)
	}
	buf, err := ioutil.ReadFile(tplFile)
	if err != nil {
		return nil, fmt.Errorf("[%s] template open the failure", tplFile)
	}
	err = yaml.Unmarshal(buf, &c)
	if err != nil {
		return nil, fmt.Errorf("[%s] Template format error", tplFile)
	}

	c.Proxy = nil

	var proxy []map[string]interface{}
	var proxiesStr []string
	names := map[string]int{}

	for _, proto := range proxies {
		o := reflect.ValueOf(proto)
		nameField := o.FieldByName("Name")
		proxyItem := make(map[string]interface{})
		j, _ := json.Marshal(proto)
		_ = json.Unmarshal(j, &proxyItem)

		name := nameField.String()
		if index, ok := names[name]; ok {
			names[name] = index + 1
			name = fmt.Sprintf("%s%d", name, index+1)
		} else {
			names[name] = 0
		}

		proxyItem["name"] = name
		proxy = append(proxy, proxyItem)
		c.Proxy = append(c.Proxy, proxyItem)
		proxiesStr = append(proxiesStr, name)
	}

	c.Proxy = proxy

	for _, group := range c.ProxyGroup {
		groupProxies := group["proxies"].([]interface{})
		for i, proxie := range groupProxies {
			groupProxies = groupProxies[:i]
			var tmpGroupProxies []string
			for _, s := range groupProxies {
				tmpGroupProxies = append(tmpGroupProxies, s.(string))
			}
			switch proxie {
			case "all":
				tmpGroupProxies = append(tmpGroupProxies, proxiesStr...)
			case "jp":
				fallthrough
			case "hk":
				fallthrough
			case "us":
				if len(regionList[proxie.(string)]) == 0 {
					for _, ps := range proxiesStr {
						if strings.Contains(ps, regionCode[proxie.(string)]) {
							tmpGroupProxies = append(tmpGroupProxies, ps)
						}
					}
				} else {
					tmpGroupProxies = regionList[proxie.(string)]
				}
			}
			group["proxies"] = tmpGroupProxies
		}
	}
	// c.ProxyOld = c.Proxy
	// c.ProxyGroupOld = c.ProxyGroup
	// c.RuleOld = c.Rule

	return yaml.Marshal(c)
}

func Base64DecodeStripped(s string) ([]byte, error) {
	if i := len(s) % 4; i != 0 {
		s += strings.Repeat("=", 4-i)
	}
	s = strings.ReplaceAll(s, " ", "+")
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		decoded, err = base64.URLEncoding.DecodeString(s)
	}
	return decoded, err
}

func IsValidNode(nodeName string) bool {
	blacklist := strings.Split(os.Getenv("SUB_BLACKLIST"), ",")
	for _, keyword := range blacklist {
		if strings.Contains(nodeName, keyword) {
			return false
		}
	}
	whitelist := strings.Split(os.Getenv("SUB_WHITELIST"), ",")
	for _, keyword := range whitelist {
		if strings.Contains(nodeName, keyword) {
			return true
		}
	}

	return false
}

func ParseContent(content string) []interface{} {
	var proxies []interface{}
	b, err := Base64DecodeStripped(content)
	if err != nil {
		log.Errorf("Decode fail content %s", err)
		return proxies
	}

	scanner := bufio.NewScanner(bytes.NewReader(b))
	for scanner.Scan() {
		switch {
		case strings.HasPrefix(scanner.Text(), "ss://"):
			s := scanner.Text()
			s = strings.TrimSpace(s)
			ss := buildSS(s)
			if ss.Name != "" && IsValidNode(ss.Name) {
				proxies = append(proxies, ss)
			}
		case strings.HasPrefix(scanner.Text(), "ssr://"):
			s := scanner.Text()[6:]
			s = strings.TrimSpace(s)
			ssr := buildSSR(s)
			if ssr.Name != "" && IsValidNode(ssr.Name) {
				proxies = append(proxies, ssr)
			}
		case strings.HasPrefix(scanner.Text(), "vmess://"):
			s := scanner.Text()[8:]
			s = strings.TrimSpace(s)
			vmess := buildVMess(s)
			if vmess.Name != "" && IsValidNode(vmess.Name) {
				proxies = append(proxies, vmess)
			}
		}

	}

	return proxies
}
