package clash

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

const TplFile = "/configs/base.yaml"

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
}

func (c *Clash) LoadTemplate(proxies []interface{}) ([]byte, error) {
	_, err := os.Stat(TplFile)
	if err != nil && os.IsNotExist(err) {
		return nil, err
	}
	buf, err := os.ReadFile(TplFile)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(buf, &c)
	if err != nil {
		return nil, err
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
			case "sg":
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
	// Get blacklist and whitelist from env
	blacklistStr := os.Getenv("SUB_BLACKLIST")
	whitelistStr := os.Getenv("SUB_WHITELIST")

	// Check blacklist first
	if blacklistStr != "" {
		blacklist := strings.Split(blacklistStr, ",")
		for _, keyword := range blacklist {
			if keyword != "" && strings.Contains(nodeName, keyword) {
				log.Printf("Node %s blocked by blacklist rule: %s", nodeName, keyword)
				return false
			}
		}
	}

	// If whitelist is empty, allow all non-blacklisted nodes
	if whitelistStr == "" {
		return true
	}

	// Check whitelist
	whitelist := strings.Split(whitelistStr, ",")
	for _, keyword := range whitelist {
		if keyword != "" && strings.Contains(nodeName, keyword) {
			log.Printf("Node %s allowed by whitelist rule: %s", nodeName, keyword)
			return true
		}
	}

	// If whitelist is not empty and no match found, reject the node
	log.Printf("Node %s rejected (not in whitelist)", nodeName)
	return false
}

func ParseContent(content string) []interface{} {
	var proxies []interface{}
	b, err := Base64DecodeStripped(content)
	if err != nil {
		log.Printf("Decode fail content %s", err)
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
		default:
			log.Println(scanner.Text())
		}

	}

	return proxies
}
