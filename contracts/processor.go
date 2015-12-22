package contracts

type Processor interface {
	// Process a fetched resource with its content, title, its children URLs
	// and an optional error
	Process(url, content, title string, urls []string, err error)
}
