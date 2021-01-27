package clash

import (
	"net/url"
	"regexp"
	"strings"
	"sub2clashr/log"
)

var ssReg = regexp.MustCompile(`(?m)ss://(\w+)@([^:]+):(\d+)#(.+)`)

type ClashSS struct {
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Server     string      `json:"server"`
	Port       interface{} `json:"port"`
	Password   string      `json:"password"`
	Cipher     string      `json:"cipher"`
	Plugin     string      `json:"plugin"`
	PluginOpts PluginOpts  `json:"plugin-opts"`
}

type PluginOpts struct {
	Mode string `json:"mode"`
	Host string `json:"host"`
}

func buildSS(s string) ClashSS {
	s, err := url.PathUnescape(s)
	if err != nil {
		log.Errorf("Decode ss config err %s", err)
		return ClashSS{}
	}

	findStr := ssReg.FindStringSubmatch(s)
	if len(findStr) < 5 {
		log.Infof("Decode ss config err %s", "findStr<5")
		return ClashSS{}
	}
	rawSSRConfig, err := Base64DecodeStripped(findStr[1])
	if err != nil {
		log.Errorf("Decode ss config %s", err)
		return ClashSS{}
	}
	params := strings.Split(string(rawSSRConfig), `:`)
	if 2 != len(params) {
		log.Errorf("Decode ss config %s", "params<2")
		return ClashSS{}
	}

	ss := ClashSS{}
	ss.Type = "ss"
	ss.Cipher = params[0]
	ss.Password = params[1]
	ss.Server = findStr[2]
	ss.Port = findStr[3]

	// ss.Plugin = findStr[4]
	// switch {
	// case strings.Contains(ss.Plugin, "obfs"):
	// 	ss.Plugin = "obfs"
	// }

	// p := PluginOpts{
	// 	Mode: findStr[5],
	// }
	// p.Host = findStr[6]

	ss.Name = findStr[4]
	// ss.PluginOpts = p

	return ss
}
