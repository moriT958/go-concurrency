package exercise

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher) {
	var (
		cache map[string]bool
		mu    sync.Mutex
		wg    sync.WaitGroup
		crawl func(url string, depth int)
	)

	cache = make(map[string]bool)
	crawl = func(url string, depth int) {
		defer wg.Done()
		if depth <= 0 {
			return
		}

		mu.Lock()
		if cache[url] {
			mu.Unlock()
			return
		}
		cache[url] = true
		mu.Unlock()

		body, urls, err := fetcher.Fetch(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("found: %s %q\n", url, body)
		for _, u := range urls {
			wg.Add(1)
			go crawl(u, depth-1)
		}
	}

	wg.Add(1)
	go crawl(url, depth)
	wg.Wait()
}

// this is not be parallelized.
func CrawlSync(url string, depth int, fetcher Fetcher) {
	var (
		cache map[string]bool
		crawl func(url string, depth int)
	)

	cache = make(map[string]bool)
	crawl = func(url string, depth int) {
		if depth <= 0 {
			return
		}

		if cache[url] {
			return
		}
		cache[url] = true

		body, urls, err := fetcher.Fetch(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("found: %s %q\n", url, body)
		for _, u := range urls {
			Crawl(u, depth-1, fetcher)
		}
	}

	crawl(url, depth)
}
