package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
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
	workDir, spec, subs, api string
)

func init() {
	flag.StringVar(&workDir, "d", ".", "specify directory")
	flag.Parse()
}

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
		b, err := ioutil.ReadFile(filepath.Join(workDir, "config.yaml"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}
		_, _ = w.Write(b)
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

	b, err := clash.Sub2byte(urls, workDir)
	if err != nil {
		log.Errorf("Sub2byte failed, %s", err)
		return
	}

	filename := filepath.Join(workDir, "config.yaml")
	err = ioutil.WriteFile(filepath.Join(workDir, "config.yaml"), b, 0644)
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
	// req.Debug = true
	_, err = req.Put(u, req.BodyJSON(&foo))
	if err != nil {
		log.Errorf("put config file failed, %s", err)
		return
	}
	log.Infof("req put data %s", u)
}
