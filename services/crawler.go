package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"

	// "log"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/proxy"

	"github.com/danielavshalumov/around/config"
	"github.com/danielavshalumov/around/lib"
	"github.com/danielavshalumov/around/models"

	"github.com/chromedp/chromedp"
	"github.com/redis/go-redis/v9"
)

type Browser interface {
	GetQuery(query string, params bool) string
	CrawlSerp(link string, current_url string) string
	GetReferer() string
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
	Referrer string
}

func NewDuckDuckGo() *DuckDuckGo {
	return &DuckDuckGo{
		StartUrl: "https://html.duckduckgo.com/html?q=",
		Referrer: "https://duckduckgo.com",
	}
}

type CrawlerService struct {
	DB  *config.Db
	RDB *redis.Client
}

func NewCrawlerService(db *config.Db, rdb *redis.Client) *CrawlerService {
	return &CrawlerService{
		DB:  db,
		RDB: rdb,
	}
}

type CrawlerWorker struct {
	browser Browser

	mu               sync.RWMutex
	wg               sync.WaitGroup
	semaphore        chan struct{}
	chromedpThrottle chan struct{}
	count            int32
	threadCount      int32
	limitReached     atomic.Bool
	browserCtx       context.Context
	ctx              context.Context
	browserCancel    context.CancelFunc
	cancel           context.CancelFunc
	crawlmu          sync.RWMutex
	tabs             map[int]context.CancelFunc
	tabCounter       int
}

func NewCrawlerWorker(maxThreads int) *CrawlerWorker {

	cw := &CrawlerWorker{
		semaphore:        make(chan struct{}, maxThreads),
		chromedpThrottle: make(chan struct{}, maxThreads/5),
	}
	cw.limitReached.Store(false)
	return cw
}

func (g *Google) CrawlSerp(link string, current_url string) string {
	fmt.Println(link, current_url)
	return ""
}

func (g *Google) GetQuery(query string, params bool) string {
	fmt.Printf(fmt.Sprintf("%s%s", g.StartUrl, url.QueryEscape(query)))
	return fmt.Sprintf("%s%s", g.StartUrl, url.QueryEscape(query))
}

func (b *DuckDuckGo) CrawlSerp(link string, current_url string) string {
	var next_url string
	if !strings.Contains(link, "https") {
		return ""
	}
	if strings.Contains(link, "duckduckgo") && strings.Contains(link, "https") && strings.Contains(link, "&") {
		link_mal := link[strings.Index(link, "https"):]
		next_url = link_mal[:strings.Index(link_mal, "&")]
	} else {
		next_url = link
	}
	new_next_url, err := url.PathUnescape(next_url)
	if err != nil {
		fmt.Println("error unescaping path from duckduck go impl of CrawlSerp")
	}
	return new_next_url
}

func (b *DuckDuckGo) GetQuery(query string, params bool) string {
	if params {
		EscapedQuery := url.PathEscape(query)
		return fmt.Sprintf("%s%s", b.StartUrl, EscapedQuery)
	} else {
		res := b.StartUrl + query + "&"
		return res
	}
}

func (cw *CrawlerWorker) getNextPageParams(current_url string, s *models.Spider) (url.Values, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", current_url, nil)
	req.Header.Set("User-Agent", s.UserAgent)
	req.Header.Set("Referer", cw.browser.GetReferer())
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error in getNextPageParams()")
		return nil, errors.New("Error calling http req to SERP")
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	stringBody := string(body)

	reader := strings.NewReader(stringBody)
	doc, err := html.Parse(reader)
	if err != nil {
		fmt.Printf("Error parsing HTML from page %s", current_url)
	}

	var trav func(node *html.Node)
	params := url.Values{}
	trav = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "input" {
			attrs := node.Attr
			for _, attr := range attrs {
				if attr.Key == "name" {
					switch attr.Val {
					case "s":
						params.Set("s", node.Attr[2].Val)
					case "v":
						params.Set("v", node.Attr[2].Val)
					case "o":
						params.Set("o", node.Attr[2].Val)
					case "dc":
						params.Set("dc", node.Attr[2].Val)
					case "api":
						params.Set("api", node.Attr[2].Val)
					case "vqd":
						params.Set("vqd", node.Attr[2].Val)
					case "kl":
						params.Set("kl", node.Attr[2].Val)
					}
				}
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			trav(c)
		}
	}
	trav(doc)
	return params, nil
	// "q" : ,
	// "s" : ,
	// "v" : ,
	// "o" : ,
	// "dv": ,
	// "api": ,
	// "vqd": ,
	// "kl": ,
}

func (b *DuckDuckGo) GetReferer() string {
	return "https://duckduckgo.com"
}

func (b *Google) GetReferer() string {
	return "https://google.com"
}

func (cs *CrawlerService) StartCrawl(spider *models.Spider, browser string, userId int64, parentCtx context.Context) (int32, []models.BacklinkResponse) {

	worker := NewCrawlerWorker(50)
	spider.SetUserAgent()

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("exclude-switches", "enable-automation"),
		chromedp.Flag("disable-extensions", false),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-features", "IsolateOrigins,site-per-process"),
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("renderer-process-limit", "1"),
		chromedp.Flag("process-per-site", true),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("disable-breakpad", true),
		chromedp.Flag("disable-client-side-phishing-detection", true),
		chromedp.Flag("disable-component-update", true),
		chromedp.Flag("disable-hang-monitor", true),
		chromedp.Flag("disable-ipc-flooding-protection", true),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("disable-prompt-on-repost", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("disable-translate", true),
		chromedp.Flag("metrics-recording-only", true),
		chromedp.Flag("mute-audio", true),
		chromedp.Flag("max-old-space-size", "256"),
		chromedp.UserAgent(spider.UserAgent),
		chromedp.WindowSize(1920, 1080),
	)

	alloCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	// defer cancel()
	// defer alloCancel()
	browserCtx, browserCancel := chromedp.NewContext(alloCtx)
	err := chromedp.Run(browserCtx,
		chromedp.Navigate("https://html.duckduckgo.com"))
	if err != nil {
		fmt.Println("Error")
		fmt.Println(err)
	}
	// defer browserCancel()
	_tabs := make(map[int]context.CancelFunc)
	worker.tabs = _tabs
	fmt.Println("user agent", spider.UserAgent)
	worker.count = 0
	worker.threadCount = 0
	requestCtx, requestCancel := context.WithCancel(parentCtx)
	worker.browserCtx = browserCtx
	worker.browserCancel = browserCancel
	worker.ctx = requestCtx
	worker.cancel = requestCancel
	bf := BrowserFactory{}
	worker.browser = bf.build(browser)
	fmt.Println(worker.browser)
	fmt.Println("comp_domains", spider.CompDomains)

	if err := lib.RenewIp(); err != nil {
		fmt.Printf("Failed to renew IP\nError: %v\n", err)
	}
	time.Sleep(time.Second * 3)
	pages := 3
	worker.wg.Add(1)
	go func() {
		defer worker.wg.Done()
		worker.Crawl(spider, worker.browser.GetQuery(spider.Query, true), spider.MaxDepth, pages-1)
		fmt.Println("Crawling finished inside go function")
	}()
	worker.wg.Wait()
	fmt.Println("Crawling finished")
	// fmt.Println(c1.Browser.LostConnection)
	// c1.Allocator.Wait()

	session_id, err := cs.DB.CreateNewCrawlSession(userId)

	var res []models.BacklinkResponse
	for source, target := range spider.Backlinks {
		b_id, err := cs.DB.InsertIntoBacklink(&models.Backlink{Source: source, Link: target[0], Dofollow: false, Session: session_id, Title: target[2]})
		if err != nil {
			fmt.Println("Error insterting backlink from service", err.Error())
		}
		res = append(res, models.BacklinkResponse{Source: source, Backlink: target[0], Id: b_id})
		cs.AddToMLQueue(source, target[0], target[1], parentCtx)
	}

	chromedp.Cancel(browserCtx)
	chromedp.Cancel(alloCtx)
	// browserCancel()
	// cancel()
	fmt.Println("After the cancel")
	fmt.Println("Open threads")
	fmt.Println(worker.tabs)
	fmt.Println("closed threads")
	fmt.Println(spider.Backlinks)

	// for id, tabCancel := range cw.tabs {
	// 	fmt.Printf("cancelling context %d", id)
	// 	tabCancel()
	// }

	select {
	case <-browserCtx.Done():
		er := browserCtx.Err()
		fmt.Printf("browser Cancelled erro: %v\n", er)
	default:
		fmt.Println("browser Context has not been cancelled and is still runnning")
		fmt.Println(browserCtx)
	}

	select {
	case <-alloCtx.Done():
		er := alloCtx.Err()
		fmt.Printf("Allo Cancelled erro: %v\n", er)
	default:
		fmt.Println("Allo Context has not been cancelled and is still runnning")

		chromedp.Cancel(alloCtx)
		fmt.Println(alloCtx.Err())
	}

	return 0, res
}

func (cw *CrawlerWorker) Crawl(s *models.Spider, current_url string, depth int, pages int) {

	atomic.AddInt32(&cw.threadCount, 1)
	defer atomic.AddInt32(&cw.threadCount, -1)

	currentCount := atomic.LoadInt32(&cw.count)
	if cw.limitReached.Load() || currentCount >= 10 {
		fmt.Println("limit reached")
		fmt.Printf("thread count %d\n", atomic.LoadInt32(&cw.threadCount))
		cw.cancel()
		return
	}

	if depth == 0 {
		return
	}

	time.Sleep((time.Millisecond * 1200))
	curr_parse, err := url.QueryUnescape(current_url)

	if err != nil {
		fmt.Println("Error unescaping url")
	}

	cw.mu.Lock()
	switch {
	case s.Visited[curr_parse] == true:
		// fmt.Printf("FOUND VISITED: %s\n", curr_parse)
		cw.mu.Unlock()
		return
	case s.Backlinks[curr_parse][0] != "":
		fmt.Println("Backlinks Reaached")
		cw.mu.Unlock()
		return
	}
	cw.mu.Unlock()

	select {
	case cw.semaphore <- struct{}{}:
		defer func() { <-cw.semaphore }()
	case <-cw.ctx.Done():
		return
	}

	cw.mu.Lock()
	s.Visited[curr_parse] = true
	cw.mu.Unlock()

	failed := true
	var links map[string]string
	var page_html string
	var statusCode int
	var title string
	if depth == s.MaxDepth {

		for failed {
			_links, _page_html, _statusCode, _title := cw.extractAnchorTags(current_url, true, s)
			statusCode = _statusCode
			linkCount := 0
			for _link := range _links {
				fmt.Println(_link, _title)
				if len(_link) != 0 {
					linkCount++
				}
			}
			fmt.Printf("SERP COUNT %d", linkCount)
			if statusCode != 200 || linkCount < 4 {
				fmt.Println("FAILED SERP SEARCH")
				if err := lib.RenewIp(); err != nil {
					fmt.Printf("Failed to renew IP\nError: %v\n", err)
				}
				cw.mu.Lock()
				s.SetUserAgent()
				cw.mu.Unlock()
				time.Sleep(time.Second * 1)
			} else {
				links = _links
				page_html = _page_html
				failed = false
			}
		}

	} else {
		links, page_html, statusCode, title = cw.extractAnchorTags(current_url, true, s)
	}

	fmt.Printf("Page %d Depth %d Crawling %s - %s \n", pages, depth, current_url, curr_parse)

	var backlinkFound bool
	var absolute, relative []string
	for link := range links {

		newCurrentCount := atomic.LoadInt32(&cw.count)
		if newCurrentCount >= 10 || cw.limitReached.Load() {
			fmt.Println("limit reached")
			fmt.Printf("thread count %d\n", atomic.LoadInt32(&cw.threadCount))
			cw.cancel()
			return
		}
		if err != nil {
			fmt.Println("Error unescaping url")
		}

		if link == "" || strings.Contains(link, "feedspot") || strings.Contains(link, "feedburner") {
			continue
		}

		var next_url string

		if depth == s.MaxDepth {
			if err != nil {
				fmt.Println(err)
			}
			next_url = cw.browser.CrawlSerp(link, curr_parse)
			fmt.Println("next url", next_url)
		}

		link = strings.Replace(link, "www.", "", 1)
		// Different Operations for Absolute and Relative links
		if depth < s.MaxDepth {
			if strings.HasPrefix(link, "https") {
				if depth < s.MaxDepth-1 {
					cw.mu.Lock()

					if checkBacklink(link, curr_parse, s.CompDomains, s) != "" && depth != s.MaxDepth && !backlinkFound && s.Backlinks[link][0] == "" {
						fmt.Println("------------ Backlink Found ------------")
						fmt.Println(current_url + "->" + link)
						// fmt.Print(parsed_host, "->", parsed_link_host)
						fmt.Print("Visited", s.Visited[link])
						fmt.Println("----------------------------------------")
						backlinkFound = true
						backlinkValue := [3]string{curr_parse, page_html, title}
						// fmt.Printf("PAGE HTML, %s", page_html)
						s.Backlinks[link] = backlinkValue
						cw.mu.Unlock()

						// Was thinkgin to making the value into an array, but this is probably and the top switch case is the reason for dupes
						cw.mu.Lock()
						atomic.AddInt32(&cw.count, 1)
						fmt.Println("The count is now", atomic.LoadInt32(&cw.count))
						cw.mu.Unlock()
						return
					}
					cw.mu.Unlock()
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

			if strings.Contains(path_link, "/reply") {
				continue
			}

			if !strings.Contains(path_link, "discussion") && !strings.Contains(path_link, "thread") && !strings.Contains(path_link, "forum") && !strings.Contains(path_link, "comment") && !strings.Contains(path_link, "/t") && !strings.Contains(path_link, "view") {
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

		currentCount := atomic.LoadInt32(&cw.count)
		if cw.limitReached.Load() || currentCount >= 10 {
			fmt.Println("limit reached")
			fmt.Printf("thread count %d\n", atomic.LoadInt32(&cw.threadCount))
			cw.cancel()
			return
		}

		cw.wg.Add(1)
		go func(next_url string) {
			defer cw.wg.Done()
			cw.Crawl(s, next_url, depth-1, pages)
		}(next_url)

	}

	if depth == s.MaxDepth && pages > 0 {
		time.Sleep(time.Second * 10)
		params, err := cw.getNextPageParams(current_url, s)
		if err != nil {
			fmt.Println("failed going to the next page")
		}
		next_page_url := strings.ReplaceAll(cw.browser.GetQuery(s.Query, false), " ", "%20")
		next_page_url += params.Encode()
		fmt.Println("Next parm url", next_page_url)

		cw.wg.Add(1)
		go func(next_page_url string) {
			defer cw.wg.Done()
			cw.Crawl(s, next_page_url, depth, pages-1)
		}(next_page_url)
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
	if len(filter) <= 1 {
		filter = []string{"capterra.com", "flic.kr", "youtube.com", "facebook.com", "twitter.com", "instagram.com", "pinterest.com", "google.com", "internetbrands.com", "xenforo.com", "wpforo.com", "futureplc.com", "tiktok.com", "linkedin.com", "vbulletin.com", "twitch"}
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
		backlinkCondition = !slices.Contains(filter, parsed_link_host) && (strings.Contains(link, "/p/") || strings.Contains(link, "/collection/") || strings.Contains(link, "/product/") || strings.Contains(link, "/collections/") || strings.Contains(link, "/cgi-bin/") || strings.Contains(link, "sku"))
	}
	// if strings.Contains(link, "houzz") {
	// 	if !strings.Contains(link, "vr~") {
	// 		return ""
	// 	}
	// 	return link
	// }

	// Check title of product to see if its relevant
	if backlinkCondition {
		// res, err := http.Get(link)
		// if err != nil {
		// 	return ""
		// }
		// defer res.Body.Close()
		// body, err := io.ReadAll(res.Body)
		// if err != nil {
		// 	return ""
		// }
		// HtmlBody := strings.ToLower(string(body))
		// sIndex := strings.Index(BodyLower,"<title>")
		// fmt.Println(HtmlBody)
		// fmt.Printf("%s\n%s\n", link, HtmlBody)

		// title := HtmlBody[strings.Index(HtmlBody, "<title>")+7 : strings.Index(HtmlBody, "</title>")]
		// fmt.Printf("The title is %s", title)

		return link
	}

	// if slices.Contains(s.CompDomains,parsed_link_host) {}

	return ""
}

func (cs *CrawlerService) AddToMLQueue(curr_parse string, link string, page_html string, ctx context.Context) {
	// Redis Publish and Cleanup

	// var cleanPageHtml string
	lines := strings.Split(page_html, "\n")
	cleaned := make([]string, 0, len(lines))
	cleaned = append(cleaned, link)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			cleaned = append(cleaned, line)
		}
	}
	// cleanPageHtml = strings.Join(cleaned, "\n")

	// pubMessage := fmt.Sprintf("%s:%s", curr_parse, cleanPageHtml)
	// err := cs.RDB.Publish(ctx, "around-channel", pubMessage).Err()
	// if err != nil {
	// 	message := fmt.Sprintf("Link: %s\nCurr_parse:%s\nError:%v", link, curr_parse, err)
	// 	panic(message)
	// }
	// fmt.Printf("Message Published from %s", curr_parse)

	// Set in Cache
	err := cs.RDB.Set(ctx, curr_parse, page_html, 30*time.Minute).Err()
	if err != nil {
		fmt.Println("Error seting cache")
	}
}

func (cw *CrawlerWorker) extractAnchorTags(page_url string, proxyFlag bool, s *models.Spider) (map[string]string, string, int, string) {
	// Get HTML from Page URL
	// fmt.Printf("before page_html %s\n", page_url)
	test := false
	torProxy := "127.0.0.1:9050"
	page_html, statusCode := func(page_url string) (string, int) {
		// Make the Request
		var cli *http.Client
		if proxyFlag {
			dialer, err := proxy.SOCKS5("tcp", torProxy, nil, proxy.Direct)
			if err != nil {
				fmt.Println("Error with Tor Proxy")
			}
			transport := &http.Transport{
				Dial:                dialer.Dial,
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
			}
			cli = &http.Client{
				Transport: transport,
				Timeout:   15 * time.Second,
			}
		} else {
			cli = &http.Client{}
		}
		// fmt.Printf("before http call %s\n", page_url)
		req, err := http.NewRequest("GET", page_url, nil)
		cw.mu.Lock()
		req.Header.Set("User-Agent", s.UserAgent)
		cw.mu.Unlock()
		req.Header.Set("Referer", cw.browser.GetReferer())
		res, err := cli.Do(req)
		// fmt.Printf("after http call %s\n", page_url)
		if err != nil {
			fmt.Printf("Erorr %v making GET request to: %s\n", err, page_url)
			return "", 0
		}

		if res.StatusCode != 200 && !strings.Contains(page_url, "duckduckgo") {
			fmt.Printf("Page url: %s, failed with status code %d -> Chromedp time\n", page_url, res.StatusCode)
			browserErr := cw.ctx.Err()

			if browserErr != nil {
				return "", 400
			}

			var htmlContent string

			tabContext, cancel := chromedp.NewContext(cw.browserCtx)
			tabContext, cancel = context.WithTimeout(tabContext, time.Second*90)
			defer cancel()
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("Panic Recovered: %v deleting context\n", r)
					chromedp.Cancel(tabContext)
					time.Sleep(1 * time.Second)
				}
			}()

			cw.crawlmu.Lock()
			id := cw.tabCounter
			cw.tabCounter++
			fmt.Printf("Chrome Tab # %d\n", cw.tabCounter)
			cw.tabs[id] = cancel
			cw.crawlmu.Unlock()
			defer func(id int) {
				cw.crawlmu.Lock()
				delete(cw.tabs, id)
				cw.crawlmu.Unlock()
			}(id)
			fmt.Println("Run starting")
			err := chromedp.Run(tabContext,
				chromedp.Navigate(page_url),
				// 		chromedp.ActionFunc(func(ctx context.Context) error {
				// 			_err := chromedp.Evaluate(`
				//     // Override the navigator.webdriver property
				//     Object.defineProperty(navigator, 'webdriver', {
				//         get: () => false
				//     });

				// `, nil).Do(ctx)
				// 			return _err
				// 		}),
				// chromedp.WaitVisible("body"),
				chromedp.OuterHTML("body", &htmlContent, chromedp.ByQuery),
			)
			browserErr = cw.ctx.Err()
			if err != nil || browserErr != nil {
				fmt.Printf("Chromedp faield for %s - %v\n", page_url, err)
				chromedp.Cancel(tabContext)
				fmt.Println("run Completed")
				return "", 400
			}

			// fmt.Println("Browser pages")
			// targets, err := target.GetTargets().Do(cw.ctx)
			// if err != nil {
			// 	log.Fatal("failed to get targets:", err)
			// }

			// fmt.Println("Open tabs:")
			// for _, t := range targets {
			// 	if t.Type == "page" {
			// 		fmt.Printf(" - %s (%s)\n", t.Title, t.URL)
			// 	}
			// }

			// fmt.Println(htmlContent)
			fmt.Printf("GET - chromedp - %s\n", page_url)
			chromedp.Cancel(tabContext)
			select {
			case <-tabContext.Done():
				er := tabContext.Err()
				fmt.Printf("Tab Cancelled erro: %v\n", er)
			default:
				fmt.Println("Tab Context has not been cancelled and is still runnning")
			}

			// go cancel()
			// fmt.Println(htmlContent)
			test = true
			fmt.Println("run Completed")

			return htmlContent, 200
		}
		// Return Body
		defer res.Body.Close()
		fmt.Printf("GET - %s - Status code %d\n", page_url, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		return string(body), res.StatusCode
	}(page_url)

	if test == true {
		fmt.Println("Chromedp go thread - still inside function but outside anonymous part 1")
	}

	// Read Body
	reader := strings.NewReader(page_html)
	doc, err := html.Parse(reader)
	if err != nil {
		fmt.Printf("Error parsing HTML from page %s", page_url)
	}
	res := make(map[string]string)
	var trav func(node *html.Node)
	var maxNode []*html.Node
	var title string
	trav = func(node *html.Node) {
		var href string
		rel := "none"
		if node.Data == "title" {
			fmt.Println("found the title")
			if node.FirstChild != nil {
				title = node.FirstChild.Data
			}
		}
		for _, attr := range node.Attr {

			if node.Type == html.ElementNode && node.Data == "a" {
				if node.Data == "a" {
					//TODO: Somehow find a way to make this more concise put in another method or something
					cw.mu.Lock()
					if s.Visited[attr.Val] == false && attr.Key == "href" && (!strings.Contains(attr.Val, "latest") && !strings.Contains(attr.Val, "page-") && !strings.Contains(attr.Val, "post-") && !strings.Contains(attr.Val, "#post") && !strings.HasPrefix(attr.Val, "javascript") && !strings.HasPrefix(attr.Val, "data:") && !strings.HasPrefix(attr.Val, "#") && !strings.HasPrefix(attr.Val, "tel")) {
						href = attr.Val
					}
					cw.mu.Unlock()
					if attr.Key == "rel" && attr.Val != "" {
						rel = attr.Val
					}
				}
			} else {
				if (attr.Key == "class" || attr.Key == "id") && (strings.Contains(attr.Val, "message") && (strings.Contains(attr.Val, "post") || strings.Contains(attr.Val, "body"))) {
					maxNode = append(maxNode, node)
				}
			}
		}

		res[href] = rel
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			trav(c)
		}
	}

	trav(doc)
	var buf bytes.Buffer
	for _, node := range maxNode {
		if err := html.Render(&buf, node); err != nil {
			panic(err)
		}
	}
	output := buf.String()
	if test {
		fmt.Println("Chromedp go thread - still inside function but outside anonymous - right before return function")
	}
	return res, output, statusCode, title
}
