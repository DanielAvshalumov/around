package models

type Spider struct {
	Query       string
	Visited     map[string]bool
	Backlinks   map[string]string
	UserAgent   string
	MaxDepth    int
	CompDomains []string
}

func NewSpider(startUrl string, maxDepth int, compDomains []string) *Spider {
	return &Spider{
		Visited:     make(map[string]bool),
		Backlinks:   make(map[string]string),
		CompDomains: compDomains,
		MaxDepth:    maxDepth,
		// UserAgent:   "AroundBot/1.0",
		Query:     startUrl,
		UserAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.6422.61 Safari/537.36",
	}
}
