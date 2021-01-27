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

	"sub2clash/conf"
	"sub2clash/log"
)

type Clash struct {
	Port      int `yaml:"port"`
	SocksPort int `yaml:"socks-port"`
	// RedirPort          int                      `yaml:"redir-port"`
	// Authentication     []string                 `yaml:"authentication"`
	AllowLan           bool   `yaml:"allow-lan"`
	Mode               string `yaml:"mode"`
	LogLevel           string `yaml:"log-level"`
	ExternalController string `yaml:"external-controller"`
	// ExternalUI         string                   `yaml:"external-ui"`
	// Secret             string                   `yaml:"secret"`
	// Experimental map[string]interface{}   `yaml:"experimental"`
	Dns        map[string]interface{}   `yaml:"dns"`
	Proxy      []map[string]interface{} `yaml:"proxies"`
	ProxyGroup []map[string]interface{} `yaml:"proxy-groups"`
	Rule       []string                 `yaml:"rules"`

	// 兼容
	ProxyOld      []map[string]interface{} `yaml:"Proxy"`
	ProxyGroupOld []map[string]interface{} `yaml:"Proxy Group"`
	RuleOld       []string                 `yaml:"Rule"`
}

var religionCode = map[string]string{
	"jp": "日本",
	"hk": "香港",
	"us": "美国",
}

func (c *Clash) LoadTemplate(path string, proxies []interface{}) []byte {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		log.Errorf("[%s] template doesn't exist.", path)
		return nil
	}
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("[%s] template open the failure.", path)
		return nil
	}
	err = yaml.Unmarshal(buf, &c)
	if err != nil {
		log.Errorf("[%s] Template format error.", path)
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
				for _, ps := range proxiesStr {
					if strings.Contains(ps, religionCode[proxie.(string)]) {
						tmpGroupProxies = append(tmpGroupProxies, ps)
					}
				}
			}
			group["proxies"] = tmpGroupProxies
		}
	}
	c.ProxyOld = c.Proxy
	c.ProxyGroupOld = c.ProxyGroup
	c.RuleOld = c.Rule

	d, err := yaml.Marshal(c)
	if err != nil {
		return nil
	}

	return d
}

func GenerateClashConfig(proxies []interface{}, tplFile string) ([]byte, error) {
	clash := Clash{}
	r := clash.LoadTemplate(tplFile, proxies)
	if r == nil {
		return nil, fmt.Errorf("sublink 返回数据格式不对")
	}
	return r, nil
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

func filterNode(nodeName string) bool {
	for _, keyword := range conf.Cfg.Filter {
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
			s := strings.TrimSpace(scanner.Text())
			ss := buildSS(s)
			if ss.Name != "" && !filterNode(ss.Name) {
				proxies = append(proxies, ss)
			}
		case strings.HasPrefix(scanner.Text(), "ssr://"):
			s := scanner.Text()[6:]
			s = strings.TrimSpace(s)
			ssr := buildSSR(s)
			if ssr.Name != "" && !filterNode(ssr.Name) {
				proxies = append(proxies, ssr)
			}
		case strings.HasPrefix(scanner.Text(), "vmess://"):
			s := scanner.Text()[8:]
			s = strings.TrimSpace(s)
			vmess := buildVMess(s)
			if vmess.Name != "" && !filterNode(vmess.Name) {
				proxies = append(proxies, vmess)
			}
		}

	}

	return proxies
}
