package main

import (
	"fmt"
	"net/http"

	"github.com/danielavshalumov/around/config"
	"github.com/danielavshalumov/around/handlers"
	"github.com/danielavshalumov/around/services"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {

	//Load Environment Variables
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading env variables")
	}
	// Initiate Connection to db
	db, err := config.InitDB()
	if err != nil {
		fmt.Println(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Define Services
	UserService := services.CreateUserService(db)
	CrawlerService := services.NewCrawlerService(db, rdb)
	BacklinkService := services.NewBacklinkService(db)

	// Define Handlers
	AuthHandler := handlers.NewAuthHandler(UserService)
	BacklinkHandler := handlers.NewBacklinkHandler(CrawlerService, BacklinkService)
	UserHandler := handlers.NewUserHandler(UserService)

	// Set Up Endpoints
	http.Handle("/api/user/backlinks", config.CORS("http://localhost:3000")(http.HandlerFunc(UserHandler.GetUserBacklinks)))
	http.Handle("/api/user/backlink/{id}", config.CORS("http://localhost:3000")(http.HandlerFunc(UserHandler.SaveBacklink)))
	http.Handle("/api/auth/me", config.CORS("http://localhost:3000")(http.HandlerFunc(AuthHandler.HandleMe)))
	http.Handle("/api/auth/verify-user", config.CORS("http://localhost:3000")(http.HandlerFunc(AuthHandler.HandleAuthCallback)))
	http.Handle("/forum-scrape", config.CORS("http://localhost:3000")(http.HandlerFunc(BacklinkHandler.GetBacklinks)))
	http.Handle("/back-link/{id}", config.CORS("http://localhost:3000")(http.HandlerFunc(BacklinkHandler.GetBacklink)))

	fmt.Println("Server Listening on port 8081")
	http.ListenAndServe(":8081", nil)

}
