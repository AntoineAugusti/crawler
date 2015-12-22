[![Software License](https://img.shields.io/badge/License-MIT-orange.svg?style=flat-square)](https://github.com/AntoineAugusti/crawler/blob/master/LICENSE.md)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/AntoineAugusti/crawler)

# A golang crawler

## Requirements
- A working Go installation. Binaries can be downloaded from https://golang.org/dl/ and installation instructions are available here: https://golang.org/doc/install
- The GOPATH environment variable must be set

## Getting started
You can grab this package with the following command:
```
go get gopkg.in/antoineaugusti/crawler.v0
```

And then build it:
```
cd ${GOPATH%/}/src/github.com/antoineaugusti/crawler
go build
```

## Usage
A [main file](main.go) is provided to suit most needs. This is a web crawler that saves crawled HTML content to the disk, after removing the HTML tags (almost all of them) and keeping only the "meaningful" content. I highly advise you to custom this main file or to create your own fetcher or processor to suit your need. You can then reuse the [crawler](crawlers/crawler.go) to do the work concurrently.

From the `-h` flag:
```
Usage of ./crawler:
  -start string
        address to start from (default "http://www.lemonde.fr/")
  -levels int
        depth of the web crawl (default 50)
  -savePath string
        where to saved crawled resources (default "/tmp/crawl/")
  -concurrentFetchers int
        number of fetchers to run concurrently (default 30)
  -fragment
        keep the fragment part of an URL. Example #top
  -query
        keep the query part of an URL. Example ?foo=bar
  -stayOnDomain
        do not crawl resources that are stored on another domain. (default true)
```

## Interfaces
Some components are abstracted behind interfaces, or contracts as I like to call them. We have got 2 of them:
- A [Fetcher](contracts/fetcher.go) is in charge of fetching a resource and finding children resources to explore after. The Fetcher is also in charge of determining if a children resource should be explored or not
- A [Processor](contracts/processor.go) is in charge of processing a fetched resource

If you want to customise things, you will need to implement at least one of these 2 interfaces and write a custom main file. Of course you can use the provided main file and implementations to get you started.
