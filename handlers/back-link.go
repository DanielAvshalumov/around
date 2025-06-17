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

	// keywords := req.Keywords

	query := "https://html.duckduckgo.com/html?q=" + req.Industry
	// query := "https://html.duckduckgo.com/html?q=inanchor:" + strings.Join(keywords, "+") + " " + req.Industry + " %20forums"

	comp_domain := req.Comp_domains

	fmt.Println("comp_domains", comp_domain)

	spider := models.NewSpider(query, 5, comp_domain)
	fmt.Println("competitor Domains", spider)

	crawlJobId, prospects := b.crawlerService.StartCrawl(spider, r.Context())
	fmt.Println(crawlJobId)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(prospects)
}
