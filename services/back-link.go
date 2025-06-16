package services

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"slices"

	"github.com/danielavshalumov/around/config"
	"github.com/danielavshalumov/around/models"
	"golang.org/x/net/html"
)

type CrawlerService struct {
	DB           *config.Db
	mu           sync.RWMutex
	wg           sync.WaitGroup
	semaphore    chan struct{}
	limitReached atomic.Bool
}

func NewCrawlerService(db *config.Db, maxThreads int) *CrawlerService {
	cs := &CrawlerService{
		DB:        db,
		semaphore: make(chan struct{}, maxThreads),
	}
	cs.limitReached.Store(false)
	return cs
}

func (cs *CrawlerService) StartCrawl(spider *models.Spider, ctx context.Context) (int32, map[string]string) {

	cs.wg.Add(1)
	go func() {
		defer cs.wg.Done()
		Crawl(cs, spider, ctx, spider.StartUrl, spider.MaxDepth)
	}()
	cs.wg.Wait()
	fmt.Println("Crawling finished")

	return 0, spider.Backlinks
}

func Crawl(cs *CrawlerService, s *models.Spider, ctx context.Context, current_url string, depth int) {
	if depth == 0 {
		return
	}
	cs.mu.Lock()
	switch {
	case s.Visited[current_url]:
		cs.mu.Unlock()
		return
	default:
		if cs.limitReached.Load() {
			cs.mu.Unlock()
			return
		}
	}
	cs.mu.Unlock()

	select {
	case cs.semaphore <- struct{}{}:
		defer func() { <-cs.semaphore }()
	case <-ctx.Done():
		return
	}

	cs.mu.Lock()
	s.Visited[current_url] = true
	cs.mu.Unlock()

	time.Sleep(2 * time.Second)
	fmt.Printf("Depth %d Crawling %s\n", depth, current_url)
	links := extractAnchorTags(current_url)
	var absolute, relative []string

	for link, rel := range links {
		if link == "" || strings.Contains(link, "feedspot") {
			continue
		}
		link = strings.Replace(link, "www.", "", 1)
		// Different Operations for Absolute and Relative links
		if strings.HasPrefix(link, "http") {
			cs.mu.Lock()
			if checkBacklink(link, current_url, s.CompDomains, s) != "" && depth != s.MaxDepth {
				cs.mu.Unlock()
				var dofollow bool
				if strings.Contains(rel, "nofollow") {
					dofollow = false
				} else {
					dofollow = true
				}
				cs.DB.InsertIntoBacklink(&models.Backlink{Source: current_url, Link: link, Dofollow: dofollow})
				cs.mu.Lock()
				s.Backlinks[link] = rel
				cs.mu.Unlock()
				continue
			}
			cs.mu.Unlock()
			absolute = append(absolute, link)

		} else {
			relative = append(relative, link)
		}

		var next_url string

		// Uses conditional for now, TODO will change to interface later
		if strings.HasPrefix(link, "//duckduckgo") && strings.Contains(link, "https") {
			link_mal := link[strings.Index(link, "https"):]
			next_url = link_mal[:strings.Index(link_mal, "&")]
		} else if !strings.Contains(link, "https") {
			next_url = current_url[:strings.Index(current_url, ".com")+4] + link
		}

		if depth < s.MaxDepth {
			parsed_link, err := url.Parse(link)
			if err != nil {
				fmt.Println(link)
				fmt.Println("Error parsing link")
				continue
			}
			path_link := parsed_link.Path

			if !strings.Contains(path_link, "discussions") && !strings.Contains(path_link, "thread") && !strings.Contains(path_link, "forum") && !strings.Contains(path_link, "threads") && !strings.Contains(path_link, "forums") && !strings.Contains(path_link, "comments") {
				continue
			}
			curr_parse, err := url.QueryUnescape(current_url)
			if path_link[0] != '/' {
				path_link = "/" + path_link
			}
			if !strings.Contains(curr_parse, ".com") {
				continue
			}
			next_url = "https://" + curr_parse[strings.Index(curr_parse, "https://")+8:strings.Index(curr_parse, ".com")+4] + path_link
		}

		cs.wg.Add(1)
		go func(next_url string) {
			defer cs.wg.Done()
			Crawl(cs, s, ctx, next_url, depth-1)
		}(next_url)
	}

	// Test Print

	// for _, name := range absolute {
	// 	fmt.Println(name)
	// }
	// fmt.Println("------------------------------")
	// for _, name := range relative {
	// 	fmt.Println(name)
	// }

}

func checkBacklink(link string, current_url string, filter []string, s *models.Spider) string {

	_parsed, err := url.QueryUnescape(current_url)
	if err != nil {
		fmt.Printf("Failed unescaping string or string does not need unescaping %s\n", current_url)
	}
	parsed, err := url.Parse(_parsed)
	if err != nil {
		fmt.Printf("Error parsing current_url %s", current_url)
		return ""
	}
	parsed_link, err := url.Parse(link)
	if err != nil {
		fmt.Printf("extractBackinks() ~ Error parsing result link from scraped links %s", link)
		return ""
	}

	// Checks domain for similarities
	for _, value := range strings.Split(parsed_link.Hostname(), ".") {

		if strings.Contains(value, strings.Replace(parsed.Hostname(), ".com", "", 1)) {
			return ""
		}
	}
	comp_flag := true
	if len(filter) == 0 {
		filter = []string{"youtube.com", "facebook.com", "twitter.com", "instagram.com", "pinterest.com", "google.com", "internetbrands.com", "xenforo.com", "wpforo.com"}
		comp_flag = false
	}

	if strings.Contains(link, "amazon.com/registry/") {
		return ""
	}

	parsed_link_host := parsed_link.Hostname()
	parsed_host := parsed.Hostname()

	for _, value := range strings.Split(parsed_link_host, ".") {
		split := strings.Split(parsed_host, ".")
		for _, word := range split {
			if strings.Contains(word, "com") {
				continue
			}
			if strings.Contains(value, word) {
				return ""
			}
		}
		// if slices.Contains(split, value) {
		// 	return ""
		// }
	}
	var backlinkCondition bool
	if comp_flag {
		backlinkCondition = slices.Contains(filter, parsed_link_host)
		fmt.Println("debugging: ", filter, parsed_link_host)
	} else {
		backlinkCondition = !slices.Contains(filter, parsed_link_host)
	}

	if backlinkCondition {
		fmt.Println("------------ Backlink Found ------------")
		fmt.Println(current_url + "->" + link)
		fmt.Println(parsed_host, "->", parsed_link_host)
		fmt.Println("----------------------------------------")
		return link
	}

	// if slices.Contains(s.CompDomains,parsed_link_host) {}

	return ""
}

func extractAnchorTags(page_url string) map[string]string {
	// Get HTML from Page URL
	page_html := func(page_url string) string {

		parsed_url, err := url.QueryUnescape(page_url)
		if err != nil {
			fmt.Printf("Error parsing url %s\n", page_url)
		}

		res, err := http.Get(parsed_url)
		if err != nil {
			fmt.Printf("Erorr %v making GET request to: %s\n", err, page_url)
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
				//TODO: Somehow find a way to make this more concise put in another method or something
				if attr.Key == "href" && (!strings.Contains(attr.Val, "post-") && !strings.Contains(attr.Val, "#post") && !strings.HasPrefix(attr.Val, "javascript") && !strings.HasPrefix(attr.Val, "data:") && !strings.HasPrefix(attr.Val, "#") && !strings.HasPrefix(attr.Val, "tel")) {
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
