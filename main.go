package main

import (
	"fmt"
	"net/http"

	"github.com/danielavshalumov/around/db"
	"github.com/danielavshalumov/around/handlers"
)

func main() {

	// Initiate Connection to db
	_, err := db.InitDB()
	if err != nil {
		fmt.Println(err)
	}
	// Set Up Endpoints
	http.HandleFunc("/back-link", handlers.GetBacklinks)

	fmt.Println("Server Listening on port 8080")
	http.ListenAndServe(":8080", nil)

}
