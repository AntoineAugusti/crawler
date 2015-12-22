package fetchers

import (
	"bytes"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// A fetcher that fetches resources from the web
type Web struct {
	// Keep the fragment part of an URL. Example #top"
	KeepFragment bool
	// Keep the query part of an URL. Example ?foo=bar
	KeepQuery bool
	// Is it okay to crawl resources that are stored on another domain?
	StayOnDomain bool
}

// Returns the title of the crawled URL and
// a slice of URLs found on that resource
func (w Web) Fetch(url string) (content, title string, urls []string, err error) {
	var buffer bytes.Buffer
	response, err := http.Get(url)
	if response != nil {
		defer response.Body.Close()
	}
	// Problem while making the request?
	if err != nil {
		return "", "", nil, err
	}

	n, err := io.Copy(&buffer, response.Body)
	// Empty response or error?
	if err != nil || n == 0 {
		return "", "", nil, err
	}

	content = buffer.String()
	urls, title = w.findAllLinks(&buffer, url)
	return
}

// Determine if we should crawl a link present on the baseURL.
func (w Web) ShouldCrawl(baseURL, href string) bool {
	uri, err := url.Parse(href)
	// Malformed URL?
	if err != nil {
		return false
	}
	// Relative link, so we are on the same domain
	if w.StayOnDomain && isRelativeLink(*uri) {
		return true
	}

	baseUrl, err := url.Parse(baseURL)
	// Malformed URL?
	if err != nil {
		return false
	}

	// Check that we stay on the same domain
	if w.StayOnDomain {
		return uri.Host == baseUrl.Host
	}

	return true
}

// Unify a href link from a base page to have a full URL
func (w Web) unifyURL(href, base string) (fullURL string) {
	uri, err := url.Parse(href)
	if err != nil {
		panic(err)
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		panic(err)
	}
	uri = baseUrl.ResolveReference(uri)
	// Remove the URL fragment. Parts like #top, #content...
	if !w.KeepFragment {
		uri.Fragment = ""
	}
	// Remove the query part of the URL. Parts like ?foo=bar
	if !w.KeepQuery {
		uri.RawQuery = ""
	}
	return uri.String()
}

// Find all children links on a page and the title of the page from an HTTP response
func (w Web) findAllLinks(httpBody io.Reader, baseURL string) (links []string, title string) {
	page := html.NewTokenizer(httpBody)
	for {
		tokenType := page.Next()
		// End of the page, we are done
		if tokenType == html.ErrorToken {
			return
		}
		token := page.Token()

		// Extract the page title
		// React uses <title> tags also, but they have got special attributes
		if tokenType == html.StartTagToken && token.DataAtom.String() == "title" && len(token.Attr) == 0 {
			page.Next()
			title = page.Token().Data
		}

		// Parse a link
		if tokenType == html.StartTagToken && token.DataAtom.String() == "a" {
			href, hasLink := w.extractLink(token)
			if hasLink && w.ShouldCrawl(baseURL, href) {
				links = append(links, w.unifyURL(href, baseURL))
			}
		}
	}
}

// Get the href attribute from a token
func (w Web) extractLink(token html.Token) (link string, hasLink bool) {
	link = ""
	hasLink = false
	for _, attr := range token.Attr {
		if attr.Key == "href" {
			link = attr.Val
			hasLink = true
		}
	}
	return
}

// Check if a URL is relative (host and scheme are not specified)
// and the scheme is HTTP or HTTPS. This is useful for
// mailto:foo@example.com for example
func isRelativeLink(u url.URL) bool {
	return len(u.Host) == 0 && (len(u.Scheme) == 0 || strings.HasPrefix(u.Scheme, "http"))
}

// Create a new web fetcher.
func NewWeb(keepFragment, keepQuery, stayOnDomain bool) Web {
	return Web{
		KeepFragment: keepFragment,
		KeepQuery:    keepQuery,
		StayOnDomain: stayOnDomain,
	}
}
