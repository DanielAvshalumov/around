package models

type Spider struct {
	Visited   map[string]bool
	Backlinks map[string]string
	UserAgent string
	MaxDepth  int
	StartUrl  string
}

func NewSpider(startUrl string, maxDepth int) *Spider {
	return &Spider{
		Visited:   make(map[string]bool),
		Backlinks: make(map[string]string),
		StartUrl:  startUrl,
		MaxDepth:  maxDepth,
		UserAgent: "AroundBot/1.0",
	}
}
