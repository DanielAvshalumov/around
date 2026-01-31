package services

import (
	"fmt"
	"io"
	"net/http"

	"github.com/danielavshalumov/around/config"
	"github.com/danielavshalumov/around/models"
)

type BacklinkService struct {
	DB *config.Db
}

func NewBacklinkService(db *config.Db) *BacklinkService {
	return &BacklinkService{
		DB: db,
	}
}

func (bs *BacklinkService) GetBacklink(backlinkId int) (*models.Backlink, error) {
	backlink, err := bs.DB.GetBacklinkById(backlinkId)
	if err != nil {
		fmt.Println("error getting backlink by id from DB", err.Error())
		return nil, err
	}
	return backlink, nil
}

func (bs *BacklinkService) TestUrl() (string, error) {
	url := "https://www.snowboardingforum.com/threads/k2-simple-pleasures.267228/"
	httpClient := http.Client{}
	fmt.Println("Making call")
	req, err := httpClient.Get(url)
	if err != nil {
		fmt.Println("error testing url", err.Error())
		return "", err
	}
	defer req.Body.Close()
	fmt.Println("Opening Body")
	body, err := io.ReadAll(req.Body)
	stringBody := string(body)
	fmt.Println("Printing Body")
	fmt.Println(stringBody)
	fmt.Println("Printing res")
	fmt.Println(body)
	fmt.Println(req)
	return stringBody, nil
}
