package models

type BacklinkResponse struct {
	Id       int64
	Source   string
	Backlink string
	Dofollow bool
}
