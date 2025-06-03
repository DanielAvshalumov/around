package services

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/danielavshalumov/around/models"
	"golang.org/x/net/html"
)

type CrawlerService struct {
	DB *sql.DB
}

func NewCrawlerService(db *sql.DB) *CrawlerService {
	return &CrawlerService{
		DB: db,
	}
}

func (cs *CrawlerService) StartCrawl(spider *models.Spider, ctx context.Context) int32 {

	Crawl(spider.StartUrl, spider.MaxDepth)

	return 0
}

func Crawl(next_url string, depth int) {

	links := extractAnchorTags("https://www.houzz.com/discussions/6494587/big-box-or-building-supply-store-near-sf-w-floating-vanities-in-stock#n=4")
	fmt.Println(links)

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
		fmt.Printf("Status code %d\n", res.StatusCode)
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
			fmt.Println(href, rel)
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			trav(c)
		}
	}
	trav(doc)
	return res
}
