package models

type Backlink struct {
	Id       int64
	Source   string
	Link     string
	Dofollow bool
	Title    string
	Session  int64
}

type BacklinkSaveResponse struct {
	Response string
}
