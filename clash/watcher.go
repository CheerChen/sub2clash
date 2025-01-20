package clash

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
)

// StartTemplateWatcher starts watching the template file for changes
func StartTemplateWatcher(subs string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("error creating watcher: %w", err)
	}

	err = watcher.Add(TplFile)
	if err != nil {
		watcher.Close()
		return fmt.Errorf("error adding file to watcher: %w", err)
	}

	go watchTemplate(watcher, subs)
	return nil
}

func watchTemplate(watcher *fsnotify.Watcher, subs string) {
	defer watcher.Close()

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				log.Printf("Template file updated: %s\n", event.Name)
				Update(subs)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Error watching template file: %v\n", err)
		}
	}
}

func Update(subs string) {
	urls := strings.Split(subs, ",")
	if len(urls) == 0 {
		log.Printf("Error: SUB_URLS environment variable is empty")
		return
	}

	b, err := Sub2byte(urls)
	if err != nil {
		log.Printf("Sub2byte failed: %v\n", err)
		return
	}

	filename := "/configs/config.yaml"
	err = os.WriteFile(filename, b, 0644)
	if err != nil {
		log.Printf("Writing config file failed: %v\n", err)
		return
	}

	log.Printf("Successfully wrote data to file: %s\n", filename)
}
