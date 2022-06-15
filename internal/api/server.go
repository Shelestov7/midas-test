package api

import (
	"net/http"
	"time"
)

func InitHTTPServer(router http.Handler) *http.Server {
	return &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
}
