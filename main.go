package main

import (
	"flag"

	"github.com/antoineaugusti/crawler/crawlers"
	"github.com/antoineaugusti/crawler/fetchers"
	"github.com/antoineaugusti/crawler/processors"
	"github.com/asaskevich/govalidator"
)

func main() {
	// Parse CLI options
	startAddress := flag.String("start", "http://www.lemonde.fr/", "address to start from")
	recursionLevel := flag.Int("levels", 50, "depth of the web crawl")
	nbConcurrentFetchersArg := flag.Int("concurrentFetchers", 30, "number of fetchers to run concurrently")
	keepFragment := flag.Bool("fragment", false, "keep the fragment part of an URL. Example #top")
	keepQuery := flag.Bool("query", false, "keep the query part of an URL. Example ?foo=bar")
	stayOnDomain := flag.Bool("stayOnDomain", true, "do not crawl resources that are stored on another domain.")
	savePath := flag.String("savePath", "/tmp/crawl/", "where to saved crawled resources")
	flag.Parse()

	// Basic validation CLI arguments
	if !govalidator.IsURL(*startAddress) {
		panic("Expected a valid URL to start from")
	}
	depth := max(2, *recursionLevel)
	nbConcurrentFetchers := max(1, *nbConcurrentFetchersArg)

	// Create and launch the crawler
	fetcher := fetchers.NewWeb(*keepFragment, *keepQuery, *stayOnDomain)
	processor := processors.NewSaver(*savePath)
	crawler := crawlers.NewCrawler(fetcher, processor, nbConcurrentFetchers)
	crawler.Crawl(*startAddress, depth)
}

// Find the maximum between 2 integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
