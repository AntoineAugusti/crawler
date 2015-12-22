package contracts

type Fetcher interface {
	// Returns the content and the title of the crawled resource and
	// a slice of URLs to children resources found on that resource
	Fetch(url string) (content, title string, urls []string, err error)
	// Determine if we should crawl a link present on the baseURL
	ShouldCrawl(baseURL, href string) bool
}
