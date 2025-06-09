package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/danielavshalumov/around/models"
	"github.com/danielavshalumov/around/services"
)

type BacklinkHandler struct {
	crawlerService *services.CrawlerService
}

func NewBacklinkHandler(crawler *services.CrawlerService) *BacklinkHandler {
	return &BacklinkHandler{
		crawlerService: crawler,
	}
}

func (b *BacklinkHandler) GetBacklinks(w http.ResponseWriter, r *http.Request) {

	var req models.BacklinkRequest

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.SimpleError{
			Error: "Method Not Allowed",
		})
		return
	}

	// Acquire Payload
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.SimpleError{
			Error: "Invalid JSON",
		})
	}

	query := "https://html.duckduckgo.com/html?q=" + req.Industry + "%20forums"
	comp_domain := req.Comp_domains

	spider := models.NewSpider(query, 3, comp_domain)
	fmt.Println("competitor Domains", spider)

	crawlJobId := b.crawlerService.StartCrawl(spider, r.Context())
	fmt.Println(crawlJobId)
}
