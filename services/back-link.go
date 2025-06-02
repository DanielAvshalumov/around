package services

import (
	"database/sql"

	"github.com/danielavshalumov/around/models"
)

type CrawlerService struct {
	db *sql.DB
}

func (cs *CrawlerService) StartCrawl(spider *models.Spider) int32 {

}

func Crawl() {

}
