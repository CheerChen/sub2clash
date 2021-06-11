package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sub2clash/log"

	"github.com/NYTimes/gziphandler"
	"github.com/imroc/req"
	"github.com/julienschmidt/httprouter"
	cron "github.com/robfig/cron/v3"
	"github.com/urfave/negroni"

	"sub2clash/clash"
)

var (
	workDir string
)

func init() {
	flag.StringVar(&workDir, "d", ".", "specify directory")
	flag.Parse()
}

func main() {
	spec := os.Getenv("CRON")

	c := cron.New()
	_, err := c.AddFunc(spec, Update)
	if err != nil {
		log.Fatalf("cron init failed, %s", err)
	}
	c.Start()

	router := httprouter.New()

	router.GET("/sub", Sub)
	router.GET("/convert", Convert)

	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	n.UseHandler(gziphandler.GzipHandler(router))
	n.Run(":80")
}

// Update update config file from subscribe url and put it into controller api
func Update() {
	subs := strings.Split(os.Getenv("SUB_URLS"), ",")
	if len(subs) == 0 {
		log.Errorf("env var SUB_URLS empty")
		return
	}

	b, err := clash.Sub2byte(subs, workDir)
	if err != nil {
		log.Errorf("Sub2byte failed, %s", err)
		return
	}

	err = ioutil.WriteFile(filepath.Join(workDir, "config.yaml"), b, 0644)
	if err != nil {
		log.Errorf("writing config file failed, %s", err)
		return
	}
	// trigger update
	ctrlUrl := fmt.Sprintf("http://%s/configs", os.Getenv("CLASH_CONTROLLER"))
	foo := map[string]interface{}{
		"path":    "",
		"payload": string(b),
	}
	// req.Debug = true
	_, err = req.Put(ctrlUrl, req.BodyJSON(&foo))
	if err != nil {
		log.Errorf("put config file failed, %s", err)
		return
	}
}

// Sub
// curl "http://localhost:8080/sub" -o config.yaml
func Sub(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b, err := ioutil.ReadFile(filepath.Join(workDir, "config.yaml"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	_, _ = w.Write(b)
}

// Convert
// curl "http://localhost:8080/convert?url={{urlencode}}" -o config.yaml
func Convert(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decodedValue, err := url.QueryUnescape(r.URL.Query().Get("url"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}
	_, err = url.Parse(decodedValue)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	b, err := clash.Sub2byte([]string{decodedValue}, workDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	_, _ = w.Write(b)
}
