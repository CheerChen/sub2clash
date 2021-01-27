package main

import (
	"flag"
	slog "log"
	"net/http"
	"net/url"

	"github.com/NYTimes/gziphandler"
	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"

	"sub2clash/clash"
	"sub2clash/conf"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "c", "config", "specify configuration file")
	flag.Parse()
}

func main() {
	var err error
	err = conf.Load(configFile)
	if err != nil {
		slog.Fatalf("Error reading config file, %s", err)
	}

	router := httprouter.New()

	router.GET("/sub", Sub)
	router.GET("/convert", Convert)

	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	n.UseHandler(gziphandler.GzipHandler(router))
	n.Run(":" + conf.Cfg.Port)
}

// Sub 密钥获取订阅
func Sub(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	key := r.URL.Query().Get("key")
	if key != conf.Cfg.Key {
		http.Error(w, "wrong key", http.StatusBadRequest)
		return
	}

	b, err := clash.Sub2byte(conf.Cfg.Subs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	_, _ = w.Write(b)
}

// Convert 自由转换订阅
// wget -O config.yaml -c "http://localhost:8081/convert?url={{urlencode}}
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

	b, err := clash.Sub2byte([]string{decodedValue})
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	_, _ = w.Write(b)
}
