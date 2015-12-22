package processors

import (
	"fmt"
)

// A simple processor that just outputs basic
// information to the standard output
type Output struct {
}

// Process a fetched resource with its content, title, its children URLs
// and an optional error
func (o Output) Process(url, content, title string, urls []string, err error) {
	fmt.Printf("On URL: %s ; title: %q\n", url, title)
}
