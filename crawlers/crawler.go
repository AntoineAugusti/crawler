package crawlers

import (
	"errors"
	"fmt"
	"sync"

	"github.com/antoineaugusti/crawler/contracts"
)

// Store already visited URLs. Since map cannot be accessed concurrently,
// we need to embed a mutex.
type visitedUrls struct {
	m map[string]error
	sync.Mutex
}

var loading = errors.New("URL load in progress")

type Crawler struct {
	// A fetcher is used to fetch a resource and find its children resources
	fetcher contracts.Fetcher
	// A processor is in charge of processing a single fetched resource
	processor contracts.Processor
	// Store already visited URLs, that we don't want to fetch or process twice
	urls visitedUrls
	// Dummy channel to coordinate the number of concurrent fetchers
	concurrentFetchers chan struct{}
}

// Crawl recursively starting from a given URL with a maximum depth
func (c Crawler) Crawl(url string, depth int) {
	// Let's go, crawl from the given URL
	c.crawlRecursive(url, depth)
}

func (c *Crawler) crawlRecursive(url string, depth int) {
	if depth <= 0 {
		return
	}

	c.urls.Lock()
	// If we have already visited this link, stop here
	if _, ok := c.urls.m[url]; ok {
		c.urls.Unlock()
		return
	}
	c.urls.m[url] = loading
	c.urls.Unlock()

	// Fetch children URLs. We don't want to launch too much concurrent
	// fetchers, so we wait for an available spot
	<-c.concurrentFetchers
	content, title, urls, err := c.fetcher.Fetch(url)
	// Done with the work, release one spot for another fetcher
	c.concurrentFetchers <- struct{}{}

	c.urls.Lock()
	c.urls.m[url] = err
	c.urls.Unlock()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Process the fetched resource
	c.processor.Process(url, content, title, urls, err)

	// Crawl children URLs
	done := make(chan bool)
	for _, childrenUrl := range urls {
		go func(url string) {
			c.crawlRecursive(url, depth-1)
			done <- true
		}(childrenUrl)
	}

	// Wait for children goroutines to finish
	for _ = range urls {
		<-done
	}
}

// Create a new crawler.
func NewCrawler(fetcher contracts.Fetcher, processor contracts.Processor, maxNbConcurrentFetchers int) Crawler {
	crawler := Crawler{
		fetcher:            fetcher,
		processor:          processor,
		urls:               visitedUrls{m: make(map[string]error)},
		concurrentFetchers: make(chan struct{}, maxNbConcurrentFetchers),
	}

	// Fill up the channel with dummy data to represent the number
	// of fetchers that are authorized to run concurrently.
	for i := 0; i < maxNbConcurrentFetchers; i++ {
		crawler.concurrentFetchers <- struct{}{}
	}

	return crawler
}
