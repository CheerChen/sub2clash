package main

import (
	"flag"
	"io/ioutil"
	slog "log"
	"net/http"
	"net/url"

	"github.com/NYTimes/gziphandler"
	"github.com/imroc/req"
	"github.com/julienschmidt/httprouter"
	cron "github.com/robfig/cron/v3"
	"github.com/urfave/negroni"

	"sub2clash/clash"
	"sub2clash/conf"
)

var (
	workDir string
)

func init() {
	flag.StringVar(&workDir, "d", ".", "specify directory")
	flag.Parse()
}

func main() {
	var err error
	err = conf.Load(workDir)
	if err != nil {
		slog.Fatalf("Error reading config file, %s", err)
	}

	c := cron.New()
	c.AddFunc(conf.Cfg.Spec, func() {
		b, err := clash.Sub2byte(conf.Cfg.Subs, workDir)
		if err != nil {
			slog.Printf("Sub2byte failed, %s", err)
			return
		}
		err = ioutil.WriteFile(workDir+"config.yaml", b, 0644)
		if err != nil {
			slog.Printf("Error writing config file, %s", err)
			return
		}
		// trigger update
		foo := map[string]interface{}{
			"path":    "",
			"payload": string(b),
		}
		req.Debug = true
		_, err = req.Put(conf.Cfg.ControllerApi+"/configs", req.BodyJSON(&foo))
		if err != nil {
			slog.Printf("Error put config file, %s", err)
			return
		}
	})
	c.Start()

	router := httprouter.New()

	router.GET("/sub", Sub)
	router.GET("/convert", Convert)

	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	n.UseHandler(gziphandler.GzipHandler(router))
	n.Run(conf.Cfg.Port)
}

// Sub
// curl "http://localhost:8080/sub" -o config.yaml
func Sub(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	b, err := ioutil.ReadFile(workDir + "config.yaml")
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
