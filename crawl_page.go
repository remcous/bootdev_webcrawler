package main

import (
	"fmt"
	"net/url"
	"sync"
)

type config struct {
	pages              map[string]int
	baseURL            *url.URL
	mu                 *sync.Mutex
	concurrencyControl chan struct{}
	wg                 *sync.WaitGroup
}

func newConfig(rawBaseURL string, bufferSize int) *config {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("error parsing base url: %v", err)
		return nil
	}
	return &config{
		pages:              make(map[string]int),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, bufferSize),
		wg:                 &sync.WaitGroup{},
	}
}

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		fmt.Printf("Error - crawlPage: couldn't parse URL '%s': %v\n", rawCurrentURL, err)
		return
	}
	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - crawlPage: couldn't parse URL '%s': %v\n", rawBaseURL, err)
		return
	}

	// Skip other websites
	if currentURL.Hostname() != baseURL.Hostname() {
		return
	}

	normalizedCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - normalizedURL: %v", err)
		return
	}

	if _, visited := pages[normalizedCurrentURL]; visited {
		pages[normalizedCurrentURL]++
		return
	}

	pages[normalizedCurrentURL] = 1

	fmt.Printf("crawling %s\n", rawCurrentURL)

	body, err := getHTML(rawCurrentURL)
	if err != nil {
		fmt.Printf("Error - getHTML: %v\n", err)
		return
	}

	urls, err := getURLsFromHTML(body, rawBaseURL)
	if err != nil {
		fmt.Printf("Error - getURLsFromHTML: %v\n", err)
		return
	}

	for _, url := range urls {
		crawlPage(rawBaseURL, url, pages)
	}
}
