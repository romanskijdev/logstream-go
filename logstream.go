package logstream

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
)

// Upgrader для веб-сокетов
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Клиенты, подключенные к веб-сокетам
var clients = make(map[*websocket.Conn]bool)

// Канал для передачи сообщений логов
var broadcast = make(chan string)

// Инициализация логирования
func InitLogger() {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Создаем многозадачный Writer для перенаправления логов в WebSocket и консоль
	multiWriter := io.MultiWriter(os.Stdout, httpResponseWriter{w: broadcast})
	logrus.SetOutput(multiWriter)

	go handleLogMessages()
}

// Обработчик соединений
func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Errorf("❌ Failed to upgrade to websocket: %+v", err)
		return
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var msg string
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(clients, ws)
			break
		}
	}
}

// Обработчик сообщений логов
func handleLogMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				logrus.Printf("❌ Failed to write message: %+v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

// httpResponseWriter позволяет передавать логи в канал
type httpResponseWriter struct {
	w chan string
}

func (w httpResponseWriter) Write(p []byte) (int, error) {
	w.w <- string(p)
	return len(p), nil
}
