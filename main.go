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

	// Define Handlers
	AuthHandler := handlers.NewAuthHandler(UserService)
	BacklinkHandler := handlers.NewBacklinkHandler(db, rdb)

	// Set Up Endpoints
	http.Handle("/api/auth/me", config.CORS("http://localhost:3000")(http.HandlerFunc(AuthHandler.HandleMe)))
	http.Handle("/api/auth/verify-user", config.CORS("http://localhost:3000")(http.HandlerFunc(AuthHandler.HandleAuthCallback)))
	http.Handle("/back-link", config.CORS("http://localhost:3000")(http.HandlerFunc(BacklinkHandler.GetBacklinks)))

	fmt.Println("Server Listening on port 8080")
	http.ListenAndServe(":8080", nil)

}
