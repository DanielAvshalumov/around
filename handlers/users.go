package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/danielavshalumov/around/lib"
	"github.com/danielavshalumov/around/models"
	"github.com/danielavshalumov/around/services"
)

type UserHandler struct {
	UserService *services.UserService
}

func NewUserHandler(userService *services.UserService) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}

func (u UserHandler) SaveBacklink(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(models.SimpleError{
			Error: "Method Not Allowed",
		})
		return
	}

	userId, err := lib.GetUserIdFromCookie(r)
	if err != nil {
		fmt.Println("error getting userId from cookie")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.SimpleError{
			Error: "Unauthorized JWT",
		})
		return
	}

	backlinkIdPath := r.PathValue("id")
	backlinkId, err := strconv.Atoi(backlinkIdPath)
	if err != nil {
		fmt.Println("Error with path variable", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.SimpleError{
			Error: "Bad Request",
		})
		return
	}

	var response *models.BacklinkSaveResponse
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		fmt.Println("Bady Body Request")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.SimpleError{
			Error: "Method not allowed",
		})
		return
	}
	responseInsert := response.Response
	rowsAffected, err := u.UserService.SaveBacklink(backlinkId, userId, responseInsert)
	if err != nil {
		fmt.Println("Error wth saving backlink")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.SimpleError{
			Error: "Bad Request",
		})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(rowsAffected)

}

func (u UserHandler) GetUserBacklinks(w http.ResponseWriter, r *http.Request) {
	userId, err := lib.GetUserIdFromCookie(r)
	if err != nil {
		fmt.Println("error getting user from cookie", err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.SimpleError{
			Error: "Unauthorized",
		})
		return
	}

	userBacklinks, err := u.UserService.GetUserBacklinks(userId)
	if err != nil {
		fmt.Println("User Handler - Error getting data from the service layer", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.SimpleError{
			Error: "Bad request",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userBacklinks)

}
