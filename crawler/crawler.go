/*
Controls the crawling through a website's structure, also manages the crawler state.
*/
package crawler

import (
	"golang-webcrawler/fetcher"
	"net/url"
	"strings"
	"sync"
)

type Crawled struct { // Struct to manage Crawl state in one place.
	Urls  map[string]struct{}
	Mutex sync.Mutex
}

func Crawl(page *fetcher.Page, alreadyCrawled *Crawled) {
	if page.Depth <= 0 {
		return
	}
	if hasAlreadyCrawled(page.Url, alreadyCrawled) {
		return
	}
	wg := sync.WaitGroup{}
	childPagesChan := make(chan *fetcher.Page)
	childUrls, _, _ := fetcher.FetchUrls(page.Url)
	for _, childUrl := range childUrls { // Iterate through links found on page
		wg.Add(1)
		go func(childUrl *url.URL) { // create goroutines for each link found and crawl the child page
			defer wg.Done()
			childPage := fetcher.Page{Url: childUrl, Depth: page.Depth - 1}
			Crawl(&childPage, alreadyCrawled)
			childPagesChan <- &childPage
		}(childUrl)
	}
	go func() { // Close channel when direct child pages have returned
		wg.Wait()
		close(childPagesChan)
	}()
	for childPages := range childPagesChan { // Feed channel values into slice, possibly performance inefficient.
		page.Links = append(page.Links, childPages)
	}
}

func hasAlreadyCrawled(Url *url.URL, alreadyCrawled *Crawled) bool {
	/*
		Locks alreadyCrawled, then returns true/false dependent on Url being in map.
		If false, we store the Url.
	*/
	cleanUrl := strings.TrimRight(Url.String(), "/")
	alreadyCrawled.Mutex.Lock()
	_, isPresent := alreadyCrawled.Urls[cleanUrl]
	if isPresent {
		alreadyCrawled.Mutex.Unlock()
		return true
	}
	alreadyCrawled.Urls[cleanUrl] = struct{}{}
	alreadyCrawled.Mutex.Unlock()
	return false
}
