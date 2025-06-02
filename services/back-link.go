package services

import (
	"database/sql"

	"github.com/danielavshalumov/around/models"
)

type CrawlerService struct {
	DB *sql.DB
}

func NewCrawlerService(db *sql.DB) *CrawlerService {
	return &CrawlerService{
		DB: db,
	}
}

func (cs *CrawlerService) StartCrawl(spider *models.Spider) int32 {

	return 0
}

func Crawl() {

}
