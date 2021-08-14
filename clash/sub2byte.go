package clash

import (
	"errors"
	"net/url"
	"path/filepath"
	"strings"
	"sub2clash/log"

	"github.com/imroc/req"
)

func Sub2byte(subs []string, workDir string) (b []byte, err error) {
	clash := &Clash{}
	tplFile := filepath.Join(workDir, "base.yaml")

	var proxies []interface{}
	for _, u := range subs {
		var bodyString string
		if _, err = url.Parse(strings.TrimSpace(u)); err != nil {
			log.Errorf("parse err in url %s, %s", u, err)
			continue
		}
		bodyString, err = HttpGet(u)
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

	err = GetProxies()
	if err != nil {
		regionList = make(map[string][]string)
		log.Warnf("get proxies err %s", err)
	}

	if len(proxies) == 0 {
		return nil, errors.New("proxies is empty")
	}

	return clash.LoadTemplate(tplFile, proxies)
}

func HttpGet(u string) (string, error) {
	log.Infof("req get %s", u)

	r, err := req.Get(u)

	if err != nil {
		return "", err
	}

	return r.ToString()
}
