/*
Fetches page data, converts the HTML into Urls, and formats the URLs
*/
package fetcher

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
)

type Page struct {
	Url   *url.URL
	Depth int
	Links []*Page
}

func FetchUrls(currentUrl *url.URL) (urls []*url.URL, httpError error, docError error) {
	resp, httpError := http.Get(currentUrl.String())
	if httpError != nil {
		log.Printf("failed to get URL %s: %v", currentUrl.String(), httpError)
		return
	}
	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") { // Check if HTML file
		return
	}
	defer resp.Body.Close() // Closes response body FetchUrls function is done.
	doc, docError := goquery.NewDocumentFromReader(resp.Body)
	if docError != nil {
		log.Printf("failed to parse HTML: %v", docError)
		return
	}
	localProcessed := make(map[string]struct{}) // Ensures we don't store the same Url twice and
	// end up spawning 2 goroutines for same result
	doc.Find("a").Each(func(index int, item *goquery.Selection) {
		href, ok := item.Attr("href")
		if IsRelativeUrl(href) && ok && href != "" {
			absoluteUrl := ParseRelativeUrl(currentUrl, strings.TrimRight(href, "/")) // Standardises URL
			_, isPresent := localProcessed[absoluteUrl.Path]
			if !isPresent {
				localProcessed[absoluteUrl.Path] = struct{}{}
				urls = append(urls, absoluteUrl)
			}
		}
	})
	return
}

func ParseRelativeUrl(rootUrl *url.URL, relativeUrl string) (absoluteUrl *url.URL) {
	absoluteUrl, err := url.Parse(rootUrl.Scheme + "://" + rootUrl.Host + path.Clean("/"+relativeUrl))
	if err != nil {
		return nil
	}
	absoluteUrl.Fragment = "" // Removes '#' identifiers from Url
	return
}

func IsRelativeUrl(href string) bool {
	match, _ := regexp.MatchString("^(?:[a-zA-Z]+:)?//", href)
	return !match
}
