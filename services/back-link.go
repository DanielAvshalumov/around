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

	Crawl(cs, spider, ctx, spider.StartUrl, spider.MaxDepth)

	return 0
}

func Crawl(cs *CrawlerService, s *models.Spider, ctx context.Context, next_url string, depth int) {

	cs.mu.RLock()
	switch {
	case s.Visited[next_url]:
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
	s.Visited[next_url] = true
	cs.mu.Unlock()

	links := extractAnchorTags("https://www.houzz.com/discussions/6494587/big-box-or-building-supply-store-near-sf-w-floating-vanities-in-stock#n=4")
	var absolute, relative, backlinks []string
	for link := range links {
		if strings.HasPrefix(link, "http") {
			absolute = append(absolute, link)
		} else {
			relative = append(relative, link)
		}
	}

	for _, name := range absolute {
		fmt.Println(name)
	}
	fmt.Println("------------------------------")
	for _, name := range relative {
		fmt.Println(name)
	}

}

func extractBacklinks(links map[string]string, current_url string) (map[string]string, error) {
	res := make(map[string]string)
	parsed, err := url.Parse(current_url)
	if err != nil {
		fmt.Printf("Error parsing current_url %s", current_url)
		return nil, err
	}
	for key, val := range links {
		parsed_link, err := url.Parse(key)
		if err != nil {
			fmt.Printf("extractBackinks() ~ Error parsing result link from scraped links %s", link)
			continue
		}
		if parsed_link.Hostname() != parsed.Hostname() {
			res[key] = val
			fmt.Println("------------ Backlink Found ------------")
			fmt.Println(current_url + "->" + key)
			fmt.Println("----------------------------------------")
		}
	}
	return res, nil
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
