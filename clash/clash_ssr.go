package clash

import (
	"net/url"
	"strings"
	"sub2clash/log"
)

type ClashSSR struct {
	Name             string      `json:"name"`
	Type             string      `json:"type"`
	Server           string      `json:"server"`
	Port             interface{} `json:"port"`
	Password         string      `json:"password"`
	Cipher           string      `json:"cipher"`
	Protocol         string      `json:"protocol"`
	ProtocolParam    string      `json:"protocol-param"`
	ProtocolParamOld string      `json:"protocolparam"`
	OBFS             string      `json:"obfs"`
	OBFSParam        string      `json:"obfs-param"`
	OBFSParamOld     string      `json:"obfsparam"`
}

const (
	SSRServer = iota
	SSRPort
	SSRProtocol
	SSRCipher
	SSROBFS
	SSRSuffix
)

func buildSSR(s string) ClashSSR {
	rawSSRConfig, err := Base64DecodeStripped(s)
	if err != nil {
		log.Errorf("Decode ssr config %s", err)
		return ClashSSR{}
	}
	//log.Infof(string(rawSSRConfig))
	params := strings.Split(string(rawSSRConfig), `:`)
	if len(params) != 6 {
		return ClashSSR{}
	}
	ssr := ClashSSR{}
	ssr.Type = "ssr"
	ssr.Server = params[SSRServer]
	ssr.Port = params[SSRPort]
	ssr.Protocol = params[SSRProtocol]
	ssr.Cipher = params[SSRCipher]
	ssr.OBFS = params[SSROBFS]

	// 如果兼容ss协议，就转换为clash的ss配置
	if ssr.Protocol == "origin" && ssr.OBFS == "plain" {
		switch ssr.Cipher {
		case "aes-128-gcm", "aes-192-gcm", "aes-256-gcm",
			"aes-128-cfb", "aes-192-cfb", "aes-256-cfb",
			"aes-128-ctr", "aes-192-ctr", "aes-256-ctr",
			"rc4-md5", "chacha20", "chacha20-ietf", "xchacha20",
			"chacha20-ietf-poly1305", "xchacha20-ietf-poly1305":
			ssr.Type = "ss"
		}
	}
	suffix := strings.Split(params[SSRSuffix], "/?")
	if len(suffix) != 2 {
		return ClashSSR{}
	}
	passwordBase64 := suffix[0]
	password, err := Base64DecodeStripped(passwordBase64)
	if err != nil {
		log.Errorf("Decode password %s", err)
		return ClashSSR{}
	}
	ssr.Password = string(password)

	m, err := url.ParseQuery(suffix[1])
	if err != nil {
		log.Errorf("Parse url %s", err)
		return ClashSSR{}
	}

	var de []byte
	var errs error
	for k, v := range m {
		switch k {
		case "obfsparam":
			de, errs = Base64DecodeStripped(v[0])
			ssr.OBFSParam = string(de)
			ssr.OBFSParamOld = ssr.OBFSParam
		case "protoparam":
			de, errs = Base64DecodeStripped(v[0])
			ssr.ProtocolParam = string(de)
			ssr.ProtocolParamOld = ssr.ProtocolParam
		case "remarks":
			de, errs = Base64DecodeStripped(v[0])
			ssr.Name = string(de)
		}
	}
	if errs != nil {
		log.Errorf("Decode param failed, %s", err)
		return ClashSSR{}
	}

	return ssr
}
