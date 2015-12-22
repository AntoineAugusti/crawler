package processors

import (
	"fmt"
	"html"
	"io/ioutil"
	"net/url"
	"regexp"
	"strings"

	"github.com/antoineaugusti/crawler/strip"
)

// A saver processor handles web pages by saving their content to the disk
// after getting rid of HTML tags, multiple whitespaces and so on.
type Saver struct {
	basePath string
}

// Process a fetched page with its content, title, its children URLs
// and an optional error
func (s Saver) Process(url, content, title string, urls []string, previousErr error) {
	path := s.constructPath(url)

	fmt.Printf("On URL: %s ; title: %q saved to %s\n", url, removeMultipleWhiteSpaces(title), path)

	err := ioutil.WriteFile(path, cleanContent(content), 0644)
	if err != nil {
		panic(err)
	}
}

// Determine where we are going to save the content of the fetched resource
func (s Saver) constructPath(resourceURL string) string {
	urlStruct, err := url.Parse(resourceURL)
	if err != nil {
		panic(err)
	}
	urlStruct.Scheme = ""

	cleanedURL := strings.Replace(urlStruct.String(), "/", "-", -1)[2:]
	path := s.basePath + cleanedURL + ".txt"

	return path
}

// Remove HTML tags, the doctype and multiple whitespaces from a HTML page
func cleanContent(content string) []byte {
	// Remove HTML tags
	content = strip.StripTags(content)
	// Remove doctype
	content = strings.Replace(content, "<!DOCTYPE html>", "", -1)
	content = removeMultipleWhiteSpaces(content)
	content = html.UnescapeString(content)

	return []byte(content)
}

// Replace multiple whitespaces in a string by a single whitespace
func removeMultipleWhiteSpaces(str string) string {
	reg, err := regexp.Compile("[\\s]+")
	if err != nil {
		panic(err)
	}
	return reg.ReplaceAllString(str, " ")
}

// Create a new saver processor.
func NewSaver(base string) Saver {
	return Saver{
		basePath: base,
	}
}
