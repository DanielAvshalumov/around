package main

import (
	"fmt"
	"net/http"

	"github.com/danielavshalumov/around/config"
	"github.com/danielavshalumov/around/handlers"
	"github.com/redis/go-redis/v9"
)

func main() {

	// Initiate Connection to db
	db, err := config.InitDB()
	if err != nil {
		fmt.Println(err)
	}
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Define Handlers
	BacklinkHandler := handlers.NewBacklinkHandler(db, rdb)

	// Set Up Endpoints

	http.Handle("/back-link", config.CORS("http://localhost:3000")(http.HandlerFunc(BacklinkHandler.GetBacklinks)))

	fmt.Println("Server Listening on port 8080")
	http.ListenAndServe(":8080", nil)

}
