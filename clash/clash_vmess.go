package clash

import (
	"encoding/json"
	"log"
	"strings"
)

type Vmess struct {
	Add  string      `json:"add"`
	Aid  interface{} `json:"aid"`
	Host string      `json:"host"`
	ID   string      `json:"id"`
	Net  string      `json:"net"`
	Path string      `json:"path"`
	Port interface{} `json:"port"`
	PS   string      `json:"ps"`
	TLS  string      `json:"tls"`
	Type string      `json:"type"`
	V    interface{} `json:"v"`
}

type ClashVmess struct {
	Name           string            `json:"name,omitempty"`
	Type           string            `json:"type,omitempty"`
	Server         string            `json:"server,omitempty"`
	Port           interface{}       `json:"port,omitempty"`
	UUID           string            `json:"uuid,omitempty"`
	AlterID        interface{}       `json:"alterId,omitempty"`
	Cipher         string            `json:"cipher,omitempty"`
	TLS            bool              `json:"tls,omitempty"`
	Network        string            `json:"network,omitempty"`
	WSPATH         string            `json:"ws-path,omitempty"`
	WSHeaders      map[string]string `json:"ws-headers,omitempty"`
	SkipCertVerify bool              `json:"skip-cert-verify,omitempty"`
}

func buildVMess(s string) ClashVmess {
	vmconfig, err := Base64DecodeStripped(s)
	if err != nil {
		return ClashVmess{}
	}
	vmess := Vmess{}
	err = json.Unmarshal(vmconfig, &vmess)
	if err != nil {
		log.Println(err)
		return ClashVmess{}
	}
	clashVmess := ClashVmess{}
	clashVmess.Name = vmess.PS

	clashVmess.Type = "vmess"
	clashVmess.Server = vmess.Add
	switch vmess.Port.(type) {
	case string:
		clashVmess.Port, _ = vmess.Port.(string)
	case int:
		clashVmess.Port, _ = vmess.Port.(int)
	case float64:
		clashVmess.Port, _ = vmess.Port.(float64)
	default:

	}
	clashVmess.UUID = vmess.ID
	clashVmess.AlterID = vmess.Aid
	clashVmess.Cipher = "auto"
	if strings.EqualFold(vmess.TLS, "tls") {
		clashVmess.TLS = true
	} else {
		clashVmess.TLS = false
	}
	clashVmess.Network = vmess.Net
	if vmess.Net == "ws" {
		clashVmess.WSPATH = vmess.Path
	}
	wsHeaders := make(map[string]string)
	if vmess.Host != "" {
		wsHeaders["Host"] = vmess.Host
	} else {
		wsHeaders["Host"] = vmess.Add
	}
	clashVmess.WSHeaders = wsHeaders

	return clashVmess
}
