package clash

import (
	"errors"
	"net/url"
	"sub2clash/log"
	"time"

	"github.com/imroc/req"
)

const (
	tplFile        = "base.yaml"
	timeoutDefault = 10 * time.Second
)

func Sub2byte(subs []string, workDir string) (b []byte, err error) {
	var proxies []interface{}
	for _, u := range subs {
		var bodyString string
		if _, err = url.Parse(u); err != nil {
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

		proxies = append(proxies, ParseContent(bodyString)...)
	}

	log.Infof("parse content found %d proxies", len(proxies))
	if len(proxies) == 0 {
		return nil, errors.New("proxies is empty")
	}

	return GenerateClashConfig(proxies, workDir+tplFile)
}

func HttpGet(u string) (string, error) {
	log.Infof("gorequest get %s", u)

	r, err := req.Get(u)

	if err != nil {
		return "", err
	}

	return r.ToString()
}
