package models

type BacklinkRequest struct {
	Browser      string   `json:"browser"`
	Industry     string   `json:"industry"`
	Keywords     []string `json:"keywords"`
	Comp_domains []string `json:"comp_domains"`
}
