package logstream

import (
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	path := "logs.json"
	client := InitLoggerClient(nil, &path)

	client.InitLogger()

	// endpoint config
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", client.HandleConnections)
	log.SetOutput(io.Discard)

	// panic recover middleware (example)
	wrappedMux := panicRecoveryMiddleware(mux)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { // Ваш код
		log.Println("This is a test message")
		panic("TEST PANIC")
	})

	// server start
	logrus.Info("starting server on :8080")
	go func() {
		err := http.ListenAndServe(":8080", wrappedMux)
		if err != nil {
			logrus.Fatalf("failed to start server: %+v", err)
		}
	}()

	// sending logs for every second (for example to usage)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		logrus.Error("This is a ERR log message")
		logrus.Warning("This is a Warning log message")
		logrus.Info("This is a Info log message")
		logrus.Debug("This is a Debug log message")

		ticker.Reset(60 * time.Second)
	}
}

// panic recovery middleware
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
