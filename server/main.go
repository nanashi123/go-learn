package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type HealthzHandler struct{}

func (v *HealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("X-My-Server", "Go-Learn")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	server := http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      &HealthzHandler{},
	}
	fmt.Println("ブラウザで http://localhost:8080/ を開いてください。")

	if err := server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			log.Fatalf("サーバーが異常終了しました: %v", err)
		}
	}
}
