package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/danielavshalumov/around/config"
	"github.com/danielavshalumov/around/models"
	"github.com/danielavshalumov/around/models/auth"
)

type UserHandler struct {
	DB *config.Db
}

func HandleUserAccess(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.SimpleError{
			Error: "Method Not Allowed",
		})
		return
	}

	var req auth.UserDTO
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.SimpleError{
			Error: "Invalid JSON",
		})
	}

}
