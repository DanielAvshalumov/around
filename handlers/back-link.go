package handlers

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/danielavshalumov/around/lib"
	"github.com/danielavshalumov/around/models"
	"github.com/danielavshalumov/around/services"
)

type BacklinkHandler struct {
	cs *services.CrawlerService
	bs *services.BacklinkService
}

func NewBacklinkHandler(cs *services.CrawlerService, bs *services.BacklinkService) *BacklinkHandler {
	return &BacklinkHandler{
		cs: cs,
		bs: bs,
	}
}

func (b *BacklinkHandler) GetUrl(w http.ResponseWriter, r *http.Request) {
	b.bs.TestUrl()
}

func (b *BacklinkHandler) GetBacklink(w http.ResponseWriter, r *http.Request) {
	reqBacklinkId := r.PathValue("id")
	backlinkId, err := strconv.Atoi(reqBacklinkId)
	if err != nil {
		fmt.Println("error converting request backlink ID from string to integer")
	}
	backlink, err := b.bs.GetBacklink(backlinkId)
	if err != nil {
		json.NewEncoder(w).Encode(models.SimpleError{
			Error: "Error getting backlink by ID from the handler",
		})
	}
	json.NewEncoder(w).Encode(backlink)
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

	ip, port, err := net.SplitHostPort(r.RemoteAddr)
	userId, err := lib.GetUserIdFromCookie(r)
	if err != nil {
		fmt.Println("Cookie error in bcaklink handler", err)
	}
	fmt.Printf("User Id in Backlink Handler %d\n", userId)
	if err != nil {
		fmt.Println("Error splitting ip and port")
	}
	fmt.Printf("HostPort %s:%s\n", ip, port)
	fmt.Println()
	// Acquire Payload
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.SimpleError{
			Error: "Invalid JSON",
		})
	}

	// "https://html.duckduckgo.com/html?q=\"" +
	// keywords := req.Keywords
	query := fmt.Sprintf("buying %s forum reccomendations (inurl:forum OR inurl:thread OR inurl:community inurl:discussion)", req.Industry)
	// query := fmt.Sprintf("%s forums inurl:reccomendation", req.Industry)
	// query := "https://html.duckduckgo.com/html?q=inanchor:" + strings.Join(keywords, "+") + " " + req.Industry + " %20forums"

	// comp_domains could be null
	comp_domain := req.Comp_domains
	browser := req.Browser
	fmt.Println(browser, query)
	spider := models.NewSpider(query, 4, comp_domain, req.Industry)
	goroutineCount := runtime.NumGoroutine()
	fmt.Printf("Number of go runtime that still exist before starting %d\n", goroutineCount)
	crawlJobId, prospects := b.cs.StartCrawl(spider, browser, userId, r.Context())
	fmt.Println("In the handler now after crawling")
	goroutineCount = runtime.NumGoroutine()
	fmt.Printf("Number of go runtime that still exist %d\n", goroutineCount)

	fmt.Println(crawlJobId)
	go func() {
		time.Sleep(time.Second * 10)
		goroutineCount = runtime.NumGoroutine()
		fmt.Printf("Number of go runtime that still exist %d\n", goroutineCount)
	}()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(prospects)
}
