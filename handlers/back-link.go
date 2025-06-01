package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/danielavshalumov/around/models"
)

func GetBacklinks(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.SimpleError{
			Error: "Method Not Allowed",
		})
		return
	}

	query := r.URL.Query().Get("q")
	fmt.Println(query)

}
