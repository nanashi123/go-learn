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
	w.Write([]byte("/healthz OK"))
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/healthz", &HealthzHandler{})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Write([]byte("/ OK"))
	})
	mux.HandleFunc("GET /blue/{healthz}", func(w http.ResponseWriter, r *http.Request) {
		healthz := r.PathValue("healthz") // URL の {healthz} を取得
		fmt.Fprintf(w, "/blue/%s OK", healthz)
	})

	server := http.Server{
		Addr:         ":8080",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
	fmt.Println("ブラウザで http://localhost:8080/ を開いてください。")

	if err := server.ListenAndServe(); err != nil {
		if err != http.ErrServerClosed {
			log.Fatalf("サーバーが異常終了しました: %v", err)
		}
	}
}
