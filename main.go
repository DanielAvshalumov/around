package main

import (
	"fmt"
	"net/http"

	"github.com/danielavshalumov/around/config"
	"github.com/danielavshalumov/around/handlers"
	"github.com/danielavshalumov/around/services"
)

func main() {

	fmt.Println("works")
	// Initiate Connection to db
	db, err := config.InitDB()
	if err != nil {
		fmt.Println(err)
	}

	CrawlerService := services.NewCrawlerService(db)

	BacklinkHandler := handlers.NewBacklinkHandler(CrawlerService)

	// Set Up Endpoints
	http.HandleFunc("/back-link", BacklinkHandler.GetBacklinks)

	fmt.Println("Server Listening on port 8080")
	http.ListenAndServe(":8080", nil)

}
