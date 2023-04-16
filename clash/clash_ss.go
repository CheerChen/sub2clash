package clash

import (
	"bytes"
	"net/url"
	"regexp"
	"strings"
	"sub2clash/log"
)

var ssReg = regexp.MustCompile(`(?m)ss://(\w+)@([^:]+):(\d+)#(.+)`)
var ssReg2 = regexp.MustCompile(`(?m)ss://(\w+)#(.+)`)

type ClashSS struct {
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Server     string      `json:"server"`
	Port       interface{} `json:"port"`
	Password   string      `json:"password"`
	Cipher     string      `json:"cipher"`
	Plugin     string      `json:"plugin"`
	PluginOpts PluginOpts  `json:"plugin-opts"`
	UDP        bool        `json:"udp"`
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
	s = strings.Replace(s, "=", "", -1)

	findStr := ssReg.FindStringSubmatch(s)
	if len(findStr) < 5 {
		findStr = ssReg2.FindStringSubmatch(s)
		if len(findStr) < 3 {
			log.Errorf("Decode ss config raw %s", s)
			return ClashSS{}
		}
	}
	rawSSRConfig, err := Base64DecodeStripped(findStr[1])
	if err != nil {
		log.Errorf("Decode ss config %s", err)
		return ClashSS{}
	}
	if bytes.Contains(rawSSRConfig, []byte("@")) {
		rawSSRConfig = bytes.Replace(rawSSRConfig, []byte("@"), []byte(":"), 1)
	}
	params := strings.Split(string(rawSSRConfig), `:`)

	ss := ClashSS{}
	ss.Type = "ss"
	ss.UDP = true
	ss.Cipher = params[0]
	ss.Password = params[1]
	if len(findStr) > 3 {
		ss.Server = findStr[2]
		ss.Port = findStr[3]
		ss.Name = findStr[4]
	} else {
		ss.Server = params[2]
		ss.Port = params[3]
		ss.Name = findStr[2]
	}

	return ss
}
