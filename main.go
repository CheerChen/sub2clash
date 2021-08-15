package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/NYTimes/gziphandler"
	"github.com/imroc/req"
	"github.com/julienschmidt/httprouter"
	"github.com/robfig/cron/v3"
	"github.com/urfave/negroni"

	"sub2clash/clash"
	"sub2clash/log"
)

var (
	spec, subs, api string
)

func main() {
	spec = os.Getenv("CRON")
	subs = os.Getenv("SUB_URLS")
	api = os.Getenv("CLASH_CONTROLLER")

	c := cron.New()
	_, err := c.AddFunc(spec, Update)
	if err != nil {
		log.Fatalf("cron init failed, %s", err)
	}
	c.Start()
	Update()

	router := httprouter.New()
	router.GET("/sub", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.ServeFile(w, r, "/configs/config.yaml")
	})

	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	n.UseHandler(gziphandler.GzipHandler(router))
	n.Run(":80")
}

// Update update config file from subscribe url and put it into controller api
func Update() {
	urls := strings.Split(subs, ",")
	if len(urls) == 0 {
		log.Errorf("env var SUB_URLS empty")
		return
	}

	b, err := clash.Sub2byte(urls)
	if err != nil {
		log.Errorf("Sub2byte failed, %s", err)
		return
	}

	filename := "/configs/config.yaml"
	err = ioutil.WriteFile(filename, b, 0644)
	if err != nil {
		log.Errorf("writing config file failed, %s", err)
		return
	}
	log.Infof("writes data to a file %s", filename)
	// trigger update
	u := fmt.Sprintf("http://%s/configs", api)
	foo := map[string]interface{}{
		"path":    "",
		"payload": string(b),
	}
	req.Debug = true
	_, err = req.Put(u, req.BodyJSON(&foo))
	if err != nil {
		log.Errorf("put config file failed, %s", err)
		return
	}
	log.Infof("req put data %s", u)
}
