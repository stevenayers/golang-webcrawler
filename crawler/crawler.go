package crawler

import (
	"golang-webcrawler/fetcher"
	"net/url"
	"sync"
)

type Crawled struct {
	Uris  map[string]struct{}
	Mutex sync.Mutex
}

func Crawl(page *fetcher.Page, alreadyCrawled *Crawled) {
	if page.Depth <= 0 {
		return
	}
	if hasAlreadyCrawled(page.URL, alreadyCrawled) {
		return
	}
	uris := fetcher.FetchURLs(page.URL)
	wg := sync.WaitGroup{}
	childPagesChan := make(chan *fetcher.Page)
	for _, uri := range uris {
		wg.Add(1)
		go func(uri *url.URL) {
			defer wg.Done()
			childPage := fetcher.Page{URL: uri, Depth: page.Depth - 1}
			Crawl(&childPage, alreadyCrawled)
			childPagesChan <- &childPage
		}(uri)
	}
	go func() {
		wg.Wait()
		close(childPagesChan)
	}()
	for childPages := range childPagesChan {
		page.Links = append(page.Links, childPages)
	}
}

func hasAlreadyCrawled(uri *url.URL, alreadyCrawled *Crawled) bool {
	alreadyCrawled.Mutex.Lock()
	_, isPresent := alreadyCrawled.Uris[uri.String()]
	if isPresent {
		alreadyCrawled.Mutex.Unlock()
		return true
	}
	alreadyCrawled.Uris[uri.String()] = struct{}{} //add url to list of those seen
	alreadyCrawled.Mutex.Unlock()
	return false
}
