package main

import (
	"log"
	"net/http"
	"os"

	"github.com/NYTimes/gziphandler"
	"github.com/julienschmidt/httprouter"
	"github.com/robfig/cron/v3"
	"github.com/urfave/negroni"

	"sub2clash/clash"
)

func main() {
	// Initialize environment variables
	spec := os.Getenv("CRON")
	subs := os.Getenv("SUB_URLS")

	// Start cron job for subscription updates
	c := cron.New()
	_, err := c.AddFunc(spec, func() { clash.Update(subs) })
	if err != nil {
		log.Printf("cron init failed: %v\n", err)
		return
	}
	c.Start()
	clash.Update(subs) // Initial update

	// Start template file watcher
	err = clash.StartTemplateWatcher(subs)
	if err != nil {
		log.Printf("failed to start template watcher: %v\n", err)
		return
	}

	// Setup HTTP server
	router := httprouter.New()
	router.GET("/sub", serveConfig)

	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	n.UseHandler(gziphandler.GzipHandler(router))
	n.Run(":80")
}

func serveConfig(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.ServeFile(w, r, "/configs/config.yaml")
}
