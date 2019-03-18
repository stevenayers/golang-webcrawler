# Notes
## Why Golang?
I first looked at Python AsyncIO, because I'm more confident writing production code in Python than Go, but as AsyncIO is a library, it felt clunky in comparison to the native concurrency of Go. Never touching Go and having written lots of Python was outweighed by the suitability and usability of the language in this scenario.

## Design
- Initially started drafting in one file until it became clear the logic could be split.
- Package for crawl and recursive logic, package for pages and processing the data.
- Decided to go for a nested structure. This more reflects the nature of a site map.
- Goroutines spawned in a recursive loop, initiating a goroutine to add recursive elements (and so on), then passing the value down the channel seems sensible.
- Flat data structure I think would be slower, and when tried caused lots of timeout errors.

## Gotchas
- Mutex on `alreadyCrawled` because of concurrent read write error.
- `fetchPage` returned on error but this didn't cause `fetchUrls` to return, causing nil pointer exception on `doc.NewDocumentFromReader(resp.Body)`. I ended up merging `fetchPage` into `FetchUrls`, having the http request in it's own function didn't seem necessary (however `FetchUrls` is too big for my liking, and given more time would break down to make it cleaner and for easier testing).
- I didn’t pass in `childlUrl` into `go func(){}()` inside the `childUrls` for loop, which somehow meant it did not always use the instance of `childUrl` currently being iterated over.
- `wg.Wait(); close(childPagesChan)` was hanging if not inside a `go func(){}()` block. Explained below:
   ```
   1. Initiates waitgroup
   2. Increments waitgroup & spawns goroutines, all block on [channel <- value]
   3. Spawns `wg.Wait(); close(childPagesChan)` goroutine, if not in goroutine
      it would hang on wg.Wait() and not continue.
   4. Hits range iterating over channel which unblocks [channel <- value]
   5. defer wg.Done() runs, wg.Wait() unblocks and channel closes.
   6. while 5 has been going on, range has consumed the channel values and exits.
   ```

## Issues
- `Page` struct could be better suited in `crawler.go` . It's never referenced inside `fetcher.go` (lack of cohesion). Coupled fetcher & crawler together.
- `Crawled` struct could be `Crawler` and have `haslreadyCrawled` and `Crawl` as methods. Same with `Page` struct and `FetchUrls` function.
- Tests maybe should have used `assert.Equal` instead of `t.Fatalf()`, however `t.Fatalf()` allows custom messages (I could have written better test error messages).
- Doesn’t retry links if get request throws an error/timeout.
- Data structure doesn’t truly reflect a site map, which can have referencing loops. This many to many relationship data structure could be reflected if the data was stored flat in something like JSON, with reference IDs attached to each page object, and these reference IDs stored under `Page.Links` rather than nesting, and then could be rendered using something like D3js.

## Research
- [Medium - Webl: A simple web crawler written in Go](https://medium.com/@a4word/webl-a-simple-web-crawler-written-in-go-c1ce50b4f687)
- [Golang Source Code - fmt tests](https://github.com/golang/go/blob/master/src/fmt/fmt_test.go)
- Google IO conference concurrency videos
- Sentdex & Jake Wright Golang videos
- Golang official docs
- StackOverflow
