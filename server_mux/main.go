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
		w.Write([]byte("/ OK"))
	})

	// ワイルドカード
	mux.HandleFunc("GET /content/{name}", func(w http.ResponseWriter, r *http.Request) {
		contentName := r.PathValue("name") // URL の {name} を取得
		fmt.Fprintf(w, "/content/%s OK", contentName)
	})

	// ルーターの階層化（blue/green）
	blueMux := http.NewServeMux()
	blueMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/blue/healthz OK"))
	})
	greenMux := http.NewServeMux()
	greenMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/green/healthz OK"))
	})

	mux.Handle("/blue/", http.StripPrefix("/blue", blueMux))    // "/blue/"のパターンはusersMuxが処理する
	mux.Handle("/green/", http.StripPrefix("/green", greenMux)) // "/green/"のパターンはusersMuxが処理する

	// ルーターの階層化（users）
	usersMux := http.NewServeMux()
	usersMux.HandleFunc("/mypage", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/users/mypage OK"))
	})
	usersMux.HandleFunc("/settings", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/users/settings OK"))
	})

	mux.Handle("/users/", http.StripPrefix("/users", usersMux)) // "/user/"のパターンはusersMuxが処理する

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
