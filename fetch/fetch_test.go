package fetch_test

import (
	"golang-webcrawler/fetch"
	"net/url"
	"testing"
)

type FetchUrlTest struct {
	Url       string
	httpError bool
}

var FetchUrlTests = []FetchUrlTest{
	{"http://example.com", false},
	{"HTTP://EXAMPLE.COM", false},
	{"https://www.exmaple.com", true},
	{"ftp://example.com/file.txt", true},
	{"//cdn.example.com/lib.js", true},
	{"/myfolder/test.txt", true},
	{"test", true},
}

type RelativeUrlTest struct {
	Url        string
	IsRelative bool
}

var RelativeUrlTests = []RelativeUrlTest{
	{"http://example.com", false},
	{"HTTP://EXAMPLE.COM", false},
	{"https://www.exmaple.com", false},
	{"ftp://example.com/file.txt", false},
	{"//cdn.example.com/lib.js", false},
	{"/myfolder/test.txt", true},
	{"test", true},
}

type ParseUrlTest struct {
	Url         string
	ExpectedUrl string
}

var ParseUrlTests = []ParseUrlTest{
	{"/myfolder/test", "http://example.com/myfolder/test"},
	{"test", "http://example.com/test"},
	{"test/", "http://example.com/test"},
	{"test#jg380gj39v", "http://example.com/test"},
}

func TestFetchUrlsHttpError(t *testing.T) {
	for _, test := range FetchUrlTests {
		Url, _ := url.Parse(test.Url)
		page := fetch.Page{Url: Url, Depth: 1}
		_, err := page.FetchUrls()
		if (err != nil) != test.httpError {
			t.Fatalf("%s returned error: %t (expected %t)", test.Url, !test.httpError, test.httpError)
		}
	}
}

// If I had more time, I could also simulate a page with a given number of links, and check that the number of links
// on the page reflect the number of links returned.
// Another test case is checking correct errors from parseDoc
// Would also test IsRelativeHtml regexs (very important to test Regex)

func TestIsRelativeUrl(t *testing.T) {
	for _, test := range RelativeUrlTests {
		if fetch.IsRelativeUrl(test.Url) != test.IsRelative {
			t.Fatalf("URL %s did not return %t", test.Url, test.IsRelative)
		}
	}
}

func TestParseRelativeUrl(t *testing.T) {
	rootUrl, _ := url.Parse("http://example.com")
	for _, test := range ParseUrlTests {
		absoluteUrl := fetch.ParseRelativeUrl(rootUrl, test.Url)
		if absoluteUrl.String() != test.ExpectedUrl {
			t.Fatalf("Relative URL %s did not match %s when parsed: %s", test.Url, test.ExpectedUrl, absoluteUrl)
		}
	}
}
