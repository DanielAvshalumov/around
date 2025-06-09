package models

type Spider struct {
	Visited     map[string]bool
	Backlinks   map[string]string
	UserAgent   string
	MaxDepth    int
	StartUrl    string
	CompDomains []string
}

func NewSpider(startUrl string, maxDepth int, compDomains []string) *Spider {
	return &Spider{
		Visited:     make(map[string]bool),
		Backlinks:   make(map[string]string),
		CompDomains: compDomains,
		StartUrl:    startUrl,
		MaxDepth:    maxDepth,
		UserAgent:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/11.1.2 Safari/605.1.15",
	}
}
