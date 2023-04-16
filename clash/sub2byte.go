package clash

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sub2clash/log"

	req "github.com/imroc/req/v3"
)

func Sub2byte(subs []string) (b []byte, err error) {
	var proxies []interface{}
	for _, u := range subs {
		var bodyString string
		if _, err = url.Parse(strings.TrimSpace(u)); err != nil {
			log.Errorf("parse err in url %s, %s", u, err)
			continue
		}
		bodyString, err = HttpGet(u, false)
		if err != nil {
			log.Errorf("get sub url err, %s", err)
			continue
		}
		if len(bodyString) == 0 {
			log.Errorf("get sub url err, %s", errors.New("the request body content is empty"))
			continue
		}

		p := ParseContent(bodyString)
		log.Infof("parse content found %d proxies", len(p))

		proxies = append(proxies, p...)
	}
	if len(proxies) == 0 {
		return nil, errors.New("proxies is empty")
	}

	err = GetProxiesWithDelay(proxies)
	if err != nil {
		regionList = make(map[string][]string)
		log.Warnf("get proxies err %s", err)
	}

	clash := &Clash{}
	return clash.LoadTemplate(proxies)
}

func HttpGet(u string, useProxy bool) (string, error) {

	log.Infof("req get %s", u)
	client := req.C().DevMode()

	if useProxy {
		api := os.Getenv("CLASH_CONTROLLER")
		api = strings.ReplaceAll(api, "9090", "7891")
		proxyUrl := fmt.Sprintf("socks5://%s", api)
		client.SetProxyURL(proxyUrl)
		log.Infof("SetProxyUrl %s", proxyUrl)
	}
	r, err := client.R().Get(u)

	if err != nil {
		return "", err
	}

	return r.ToString()
}
