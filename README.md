# WebCrawler

Web Crawler, recurring through links to n depth.


## Requirements

- Go 1.11+
- goquery package intalled in `$GOPATH`
    ```bash
    go get github.com/PuerkitoBio/goquery
    ```

## Usage

To build, run `go build main.go` in `golang-webcrawler` directory.

```bash
main --Url https://google.com --Depth 5
```
 