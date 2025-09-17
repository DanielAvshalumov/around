package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/danielavshalumov/around/config"
	"github.com/danielavshalumov/around/models"
	"github.com/danielavshalumov/around/services"
	"github.com/redis/go-redis/v9"
)

type BacklinkHandler struct {
	DB  *config.Db
	RDB *redis.Client
}

func NewBacklinkHandler(db *config.Db, rdb *redis.Client) *BacklinkHandler {
	return &BacklinkHandler{
		DB:  db,
		RDB: rdb,
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

	CrawlerService := services.NewCrawlerService(b.DB, b.RDB, 100)

	// Acquire Payload
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.SimpleError{
			Error: "Invalid JSON",
		})
	}
	// "https://html.duckduckgo.com/html?q=\"" +
	// keywords := req.Keywords
	query := fmt.Sprintf("%s forums (inurl:forum OR inurl:thread OR inurl:community inurl:discussion)", req.Industry)
	// query := "https://html.duckduckgo.com/html?q=inanchor:" + strings.Join(keywords, "+") + " " + req.Industry + " %20forums"

	// comp_domains could be null
	comp_domain := req.Comp_domains
	browser := req.Browser
	fmt.Println(browser, query)
	spider := models.NewSpider(query, 4, comp_domain)

	crawlJobId, prospects := CrawlerService.StartCrawl(spider, browser, r.Context())
	fmt.Println(crawlJobId)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(prospects)
}
