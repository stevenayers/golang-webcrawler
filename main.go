package main

import (
	"flag"
	"fmt"
	"golang-webcrawler/crawl"
	"golang-webcrawler/fetch"
	"net/url"
	"strings"
	"time"
)

var (
	Url   = flag.String("Url", "https://google.com", "Root URL of website to crawl.")
	Depth = flag.Int("Depth", 2, "Depth of crawling.")
)

func main() {
	flag.Parse()

	rootUrl, _ := url.Parse(*Url)
	rootPage := fetch.Page{Url: rootUrl, Depth: *Depth}
	crawler := crawl.Crawler{AlreadyCrawled: make(map[string]struct{})}
	startTime := time.Now()
	crawler.Crawl(&rootPage)
	duration := time.Since(startTime)

	printSiteMap(&rootPage, 0)
	fmt.Printf("Completed in %dms\n", duration.Nanoseconds()/1000000)
	fmt.Printf("Root URL: %s, Depth: %d\n", *Url, *Depth)
}

func printSiteMap(page *fetch.Page, indent int) {
	formattedUrl := strings.Repeat(" ", indent) + page.Url.String()
	fmt.Println(formattedUrl)
	if len(page.Links) > 0 {
		for _, childPage := range page.Links {
			printSiteMap(childPage, indent+2)
		}
	}
}
