package config

import (
	"net/http"
	"time"
)

type Http struct {
	Client  *http.Client
	Proxy   string
	Headers map[string]string
}

func NewHttpClient() *Http {
	return &Http{
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
		Proxy:   "",
		Headers: make(map[string]string),
	}
}

func NewHttpProxyClient(proxy string) *Http {
	return &Http{
		Proxy: proxy,
	}
}
