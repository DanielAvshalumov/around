package services

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/danielavshalumov/around/models"
	"golang.org/x/net/html"
)

type CrawlerService struct {
	DB           *sql.DB
	mu           sync.RWMutex
	wg           sync.WaitGroup
	semaphore    chan struct{}
	limitReached atomic.Bool
}

func NewCrawlerService(db *sql.DB, maxThreads int) *CrawlerService {
	cs := &CrawlerService{
		DB:        db,
		semaphore: make(chan struct{}, maxThreads),
	}
	cs.limitReached.Store(false)
	return cs
}

func (cs *CrawlerService) StartCrawl(spider *models.Spider, ctx context.Context) int32 {

	cs.wg.Add(1)
	go func() {
		defer cs.wg.Done()
		Crawl(cs, spider, ctx, spider.StartUrl, spider.MaxDepth)
		cs.wg.Wait()
	}()

	fmt.Println("Backlinks found:", spider.Backlinks)
	fmt.Println("Crawling finished")

	return 0
}

func Crawl(cs *CrawlerService, s *models.Spider, ctx context.Context, current_url string, depth int) {

	cs.mu.RLock()
	switch {
	case s.Visited[current_url]:
		cs.mu.RUnlock()
		return
	default:
		if cs.limitReached.Load() {
			cs.mu.RUnlock()
			return
		}
	}
	cs.mu.RUnlock()

	select {
	case cs.semaphore <- struct{}{}:
		defer func() { <-cs.semaphore }()
	case <-ctx.Done():
		return
	}

	cs.mu.Lock()
	s.Visited[current_url] = true
	cs.mu.Unlock()

	cs.mu.RLock()
	defer cs.mu.RUnlock()

	fmt.Println("Crawling %s", current_url)
	links := extractAnchorTags(current_url)
	var absolute, relative []string

	for link, rel := range links {
		if strings.HasPrefix(link, "http") {
			if depth != s.MaxDepth && checkBacklink(link, current_url) != "" {
				s.Backlinks[link] = rel
				continue
			}
			absolute = append(absolute, link)

		} else {
			relative = append(relative, link)
		}

		cs.wg.Add(1)
		go func(link string) {
			defer cs.wg.Done()
			Crawl(cs, s, ctx, link, depth-1)
		}(link)
	}

	// Printing for testing purposes
	for _, name := range absolute {
		fmt.Println(name)
	}
	fmt.Println("------------------------------")
	for _, name := range relative {
		fmt.Println(name)
	}

}

func checkBacklink(link string, current_url string) string {
	parsed, err := url.Parse(current_url)
	if err != nil {
		fmt.Printf("Error parsing current_url %s", current_url)
		return ""
	}
	parsed_link, err := url.Parse(link)
	if err != nil {
		fmt.Printf("extractBackinks() ~ Error parsing result link from scraped links %s", link)
		return ""
	}

	if parsed_link.Hostname() != parsed.Hostname() {
		fmt.Println("------------ Backlink Found ------------")
		fmt.Println(current_url + "->" + link)
		fmt.Println("----------------------------------------")
	}

	return link
}

func extractAnchorTags(page_url string) map[string]string {
	// Get HTML from Page URL
	page_html := func(page_url string) string {
		res, err := http.Get(page_url)
		if err != nil {
			fmt.Printf("Erorr making GET request to: %s", page_url)
			return ""
		}
		defer res.Body.Close()
		fmt.Printf("GET - %s - Status code %d\n", page_url, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		return string(body)
	}(page_url)

	reader := strings.NewReader(page_html)
	doc, err := html.Parse(reader)
	if err != nil {
		fmt.Printf("Error parsing HTML from page %s", page_url)
	}

	res := make(map[string]string)
	var trav func(node *html.Node)
	trav = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			attrs := node.Attr
			var href, rel string
			for _, attr := range attrs {
				if attr.Key == "href" && (!strings.HasPrefix(attr.Val, "javascript") && !strings.HasPrefix(attr.Val, "data:") && !strings.HasPrefix(attr.Val, "#") && !strings.HasPrefix(attr.Val, "tel")) {
					href = attr.Val
				}
				if attr.Key == "rel" {
					rel = attr.Val
				}
			}
			res[href] = rel
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			trav(c)
		}
	}
	trav(doc)
	return res
}
