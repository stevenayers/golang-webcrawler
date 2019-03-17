package fetcher

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Page struct {
	URL   *url.URL
	Depth int
	Links []*Page
}

func FetchURLs(currentURL *url.URL) (urls []*url.URL) {
	doc := fetchPage(currentURL)
	localProcessed := make(map[string]struct{})
	doc.Find("a").Each(func(index int, item *goquery.Selection) {
		href, ok := item.Attr("href")
		if isInternalLink(href) && ok && href != "" {
			uri := parseURI(currentURL, href)
			_, isPresent := localProcessed[uri.Path]
			if !isPresent {
				localProcessed[uri.Path] = struct{}{}
				urls = append(urls, uri)
			}
		}
	})
	return
}

func fetchPage(currentURL *url.URL) (doc *goquery.Document) {
	resp, err := http.Get(currentURL.String())
	if err != nil {
		log.Printf("failed to get URL %s: %v", currentURL.String(), err)
		return
	}
	defer resp.Body.Close()
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") {
		return
	}
	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("failed to parse HTML: %v", err)
		return
	}
	return
}

func parseURI(rootURL *url.URL, uri string) (absoluteURL *url.URL) {
	absoluteURL, err := url.Parse(rootURL.Scheme + "://" + rootURL.Host + uri)
	if err != nil {
		return nil
	}
	absoluteURL.Fragment = ""
	return
}

func isInternalLink(href string) bool {
	return !strings.Contains(href, "http") && strings.HasPrefix(href, "/")
}
