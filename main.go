package main

import (
	"github.com/sirupsen/logrus"
	"log"
	"logstream/logstream"
	"net/http"
)

// Middleware для обработки паник
func panicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("😱 Panic recovered: %v", r)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func main() {
	logstream.InitLogger()

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Ваш код, который может вызвать панику
		panic("Test panic!")
	})
	mux.HandleFunc("/ws", logstream.HandleConnections)

	// Оборачиваем все обработчики в middleware
	wrappedMux := panicRecoveryMiddleware(mux)

	log.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", wrappedMux)
	if err != nil {
		log.Fatalf("❌ Failed to start server: %+v", err)
	}
}
