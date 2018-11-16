package main

import (
	"io"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	// 出力メッセージは環境変数があればそれを使い、なければ"Hello, World!"を使う
	message := os.Getenv("CUSTOM_MESSAGE")
	if message == "" {
		message = "Hello, World!"
	}
	io.WriteString(w, message)
}
