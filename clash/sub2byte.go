package clash

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func Sub2byte(subs []string) ([]byte, error) {
	var proxies []interface{}

	for _, subURL := range subs {
		// Validate URL
		if _, err := url.Parse(strings.TrimSpace(subURL)); err != nil {
			log.Printf("Parse error in URL %s: %v\n", subURL, err)
			continue
		}

		// Fetch subscription content
		bodyString, err := HttpGet(subURL)
		if err != nil {
			log.Printf("Failed to get subscription URL: %v\n", err)
			continue
		}

		if len(bodyString) == 0 {
			log.Printf("Error: the request body content is empty")
			continue
		}

		// Parse the subscription content
		p := ParseContent(bodyString)
		log.Printf("Parsed content found %d proxies\n", len(p))

		proxies = append(proxies, p...)
	}

	// Check if we got any valid proxies
	if len(proxies) == 0 {
		return nil, errors.New("no valid proxies found in subscriptions")
	}

	// Get delay information for proxies
	err := GetProxiesWithDelay(proxies)
	if err != nil {
		regionList = make(map[string][]string)
		log.Printf("Warning: failed to get proxy delays: %v\n", err)
	}

	// Create and load Clash configuration
	clash := &Clash{}
	config, err := clash.LoadTemplate(proxies)
	if err != nil {
		return nil, fmt.Errorf("failed to load template: %v", err)
	}

	return config, nil
}

func HttpGet(targetURL string) (string, error) {
	client := &http.Client{}
	resp, err := client.Get(targetURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
