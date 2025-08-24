package services

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/proxy"

	"github.com/danielavshalumov/around/config"
	"github.com/danielavshalumov/around/models"
	"golang.org/x/net/html"
)

type Browser interface {
	GetQuery(query string) string
	CrawlSerp(link string, current_url string) string
}

type BrowserFactory struct {
}

func (bf *BrowserFactory) build(browser string) Browser {
	switch browser {
	case "google":
		return NewGoogle()
	case "duckduckgo":
		return NewDuckDuckGo()
	default:
		return nil
	}
}

type Google struct {
	StartUrl string
}

func NewGoogle() *Google {
	return &Google{
		StartUrl: "https://google.com/search?hl=en&q=",
	}
}

type DuckDuckGo struct {
	StartUrl string
}

func NewDuckDuckGo() *DuckDuckGo {
	return &DuckDuckGo{
		StartUrl: "https://html.duckduckgo.com/html?q=",
	}
}

type CrawlerService struct {
	browser      Browser
	DB           *config.Db
	mu           sync.RWMutex
	wg           sync.WaitGroup
	semaphore    chan struct{}
	count        int32
	limitReached atomic.Bool
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewCrawlerService(db *config.Db, maxThreads int) *CrawlerService {

	cs := &CrawlerService{
		DB:        db,
		semaphore: make(chan struct{}, maxThreads),
	}
	cs.limitReached.Store(false)
	return cs
}

func (g *Google) CrawlSerp(link string, current_url string) string {
	// if strings.Contains()
	fmt.Println(link, current_url)
	return ""
}

func (g *Google) GetQuery(query string) string {
	fmt.Printf(fmt.Sprintf("%s%s", g.StartUrl, url.QueryEscape(query)))
	return fmt.Sprintf("%s%s", g.StartUrl, url.QueryEscape(query))
}

func (b *DuckDuckGo) CrawlSerp(link string, current_url string) string {
	var next_url string
	if strings.HasPrefix(link, "//duckduckgo") && strings.Contains(link, "https") {
		link_mal := link[strings.Index(link, "https"):]
		next_url = link_mal[:strings.Index(link_mal, "&")]
	} else if !strings.Contains(link, "https") {
		next_url = current_url[:strings.Index(current_url, ".com")+4] + link
	}
	new_next_url, err := url.PathUnescape(next_url)
	if err != nil {
		fmt.Println("error unescaping path from duckduck go impl of CrawlSerp")
	}
	return new_next_url
}

func (b *DuckDuckGo) GetQuery(query string) string {
	EscapedQuery := url.PathEscape(query)
	return fmt.Sprintf("%s%s", b.StartUrl, EscapedQuery)
}

func (cs *CrawlerService) StartCrawl(spider *models.Spider, browser string, parentCtx context.Context) (int32, []models.BacklinkResponse) {

	spider.SetUserAgent()
	fmt.Println("user agent", spider.UserAgent)
	cs.count = 0
	ctx, cancel := context.WithCancel(parentCtx)
	cs.ctx = ctx
	cs.cancel = cancel
	bf := BrowserFactory{}
	cs.browser = bf.build(browser)
	fmt.Println(cs.browser)
	fmt.Println("comp_domains", spider.CompDomains)
	cs.wg.Add(1)
	go func() {
		defer cs.wg.Done()
		cs.Crawl(spider, cs.browser.GetQuery(spider.Query), spider.MaxDepth)
	}()
	cs.wg.Wait()
	fmt.Println("Crawling finished")

	var res []models.BacklinkResponse
	for source, target := range spider.Backlinks {
		res = append(res, models.BacklinkResponse{Source: source, Backlink: target})
	}

	return 0, res
}

func (cs *CrawlerService) Crawl(s *models.Spider, current_url string, depth int) {
	if depth == 0 {
		return
	}
	time.Sleep((time.Millisecond * 1200))
	curr_parse, err := url.QueryUnescape(current_url)
	if err != nil {
		fmt.Println("Error unescaping url")
	}
	cs.mu.Lock()
	currentCount := atomic.LoadInt32(&cs.count)
	if cs.limitReached.Load() || currentCount >= 10 {
		fmt.Println("limit reached")
		cs.cancel()
		cs.mu.Unlock()
		return
	}
	switch {
	case s.Visited[curr_parse]:
		cs.mu.Unlock()
		return
	case s.Backlinks[curr_parse] != "":
		cs.mu.Unlock()
		return
	}
	cs.mu.Unlock()

	select {
	case cs.semaphore <- struct{}{}:
		defer func() { <-cs.semaphore }()
	case <-cs.ctx.Done():
		return
	}

	cs.mu.Lock()
	s.Visited[curr_parse] = true
	cs.mu.Unlock()

	time.Sleep(2 * time.Second)
	fmt.Printf("Depth %d Crawling %s\n", depth, current_url)
	// Separate here

	// links := extractAnchorTags(curr_parse, (depth == s.MaxDepth))
	links := extractAnchorTags(current_url, true, s)

	var absolute, relative []string
	for link, rel := range links {

		newCurrentCount := atomic.LoadInt32(&cs.count)
		if newCurrentCount >= 10 || cs.limitReached.Load() {
			fmt.Println("limit reached")
			cs.cancel()
			return
		}

		if s.Visited[link] {
			continue
		}

		if link == "" || strings.Contains(link, "feedspot") || strings.Contains(link, "feedburner") {
			continue
		}

		var next_url string

		if depth == s.MaxDepth {
			// Uses conditional for now, TODO will change to interface later
			next_url = cs.browser.CrawlSerp(link, curr_parse)
			fmt.Println("next url", next_url)
		}

		link = strings.Replace(link, "www.", "", 1)
		// Different Operations for Absolute and Relative links
		if depth < s.MaxDepth {

			if strings.HasPrefix(link, "https") {
				if depth < s.MaxDepth-1 {
					cs.mu.Lock()

					if checkBacklink(link, curr_parse, s.CompDomains) != "" && depth != s.MaxDepth {
						cs.mu.Unlock()
						var dofollow bool
						if strings.Contains(rel, "nofollow") {
							dofollow = false
						} else {
							dofollow = true
						}
						cs.DB.InsertIntoBacklink(&models.Backlink{Source: curr_parse, Link: link, Dofollow: dofollow})
						// Was thinkgin to making the value into an array, but this is probably and the top switch case is the reason for dupes
						cs.mu.Lock()
						s.Backlinks[link] = curr_parse
						atomic.AddInt32(&cs.count, 1)
						fmt.Println("The count is now", atomic.LoadInt32(&cs.count))
						if atomic.LoadInt32(&cs.count) > 9 {
							cs.mu.Unlock()
							return
						}
					}
					cs.mu.Unlock()
				}
				absolute = append(absolute, link)

			} else {
				relative = append(relative, link)
			}

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

			if path_link[0] != '/' {
				path_link = "/" + path_link
			}
			// TODO add more TLD functionality
			if !strings.Contains(curr_parse, ".com") {
				continue
			}
			next_url = "https://" + curr_parse[strings.Index(curr_parse, "https://")+8:strings.Index(curr_parse, ".com")+4] + path_link
		}

		cs.wg.Add(1)
		go func(next_url string) {
			defer cs.wg.Done()
			cs.Crawl(s, next_url, depth-1)
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

func checkBacklink(link string, current_url string, filter []string) string {

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

	// Checks domain for similarities
	for _, value := range strings.Split(parsed_link.Hostname(), ".") {

		if strings.Contains(value, strings.Replace(parsed.Hostname(), ".com", "", 1)) {
			return ""
		}
	}
	comp_flag := true
	if len(filter) == 0 {
		filter = []string{"youtube.com", "facebook.com", "twitter.com", "instagram.com", "pinterest.com", "google.com", "internetbrands.com", "xenforo.com", "wpforo.com", "futureplc.com", "tiktok.com", "linkedin.com", "vbulletin.com"}
		comp_flag = false
	}

	if strings.Contains(link, "amazon.com/registry/") || strings.Contains(link, "utm") {
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
		// fmt.Println("debugging: ", filter, parsed_link_host)
	} else {
		backlinkCondition = !slices.Contains(filter, parsed_link_host)
	}

	if backlinkCondition && (strings.Contains(link, "/p/") || strings.Contains(link, "collections") || strings.Contains(link, "product")) {
		fmt.Println("------------ Backlink Found ------------")
		fmt.Println(current_url + "->" + link)
		fmt.Print(parsed_host, "->", parsed_link_host)
		fmt.Println("----------------------------------------")
		return link
	}

	// if slices.Contains(s.CompDomains,parsed_link_host) {}

	return ""
}

func extractAnchorTags(page_url string, proxyFlag bool, s *models.Spider) map[string]string {
	// Get HTML from Page URL
	torProxy := "127.0.0.1:9050"
	page_html := func(page_url string) string {
		// Make the Request
		var cli *http.Client
		if proxyFlag {
			dialer, err := proxy.SOCKS5("tcp", torProxy, nil, proxy.Direct)
			if err != nil {
				fmt.Println("Error with Tor Proxy")
			}
			transport := &http.Transport{
				Dial: dialer.Dial,
			}
			cli = &http.Client{
				Transport: transport,
				Timeout:   30 * time.Second,
			}
		} else {
			cli = &http.Client{}
		}

		req, err := http.NewRequest("GET", page_url, nil)
		req.Header.Set("User-Agent", s.UserAgent)
		res, err := cli.Do(req)
		if err != nil {
			fmt.Printf("Erorr %v making GET request to: %s\n", err, page_url)
			return ""
		}
		// Return Body
		defer res.Body.Close()
		fmt.Printf("GET - %s - Status code %d\n", page_url, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		return string(body)
	}(page_url)

	// Read Body
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
			var href string
			rel := "none"
			for _, attr := range attrs {
				//TODO: Somehow find a way to make this more concise put in another method or something
				if attr.Key == "href" && (!strings.Contains(attr.Val, "latest") && !strings.Contains(attr.Val, "page-") && !strings.Contains(attr.Val, "post-") && !strings.Contains(attr.Val, "#post") && !strings.HasPrefix(attr.Val, "javascript") && !strings.HasPrefix(attr.Val, "data:") && !strings.HasPrefix(attr.Val, "#") && !strings.HasPrefix(attr.Val, "tel")) {
					href = attr.Val
				}
				if attr.Key == "rel" && attr.Val != "" {
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
