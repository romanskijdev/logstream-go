package logstream

import (
	"encoding/json"
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
	// Создаем стандартный форматтер для вывода в консоль
	logrus.SetFormatter(&logrus.TextFormatter{})

	// Создаем многозадачный Writer для вывода в консоль и кастомный Writer для WebSocket
	multiWriter := io.MultiWriter(os.Stdout, newWebSocketWriter())
	logrus.SetOutput(multiWriter)

	go handleLogMessages()
}

// Кастомный Writer для WebSocket, который форматирует логи в JSON
type webSocketWriter struct{}

func newWebSocketWriter() *webSocketWriter {
	return &webSocketWriter{}
}

func (w *webSocketWriter) Write(p []byte) (int, error) {
	var logEntry map[string]interface{}
	err := json.Unmarshal(p, &logEntry)
	if err != nil {
		return 0, err
	}
	jsonLog, err := json.Marshal(logEntry)
	if err != nil {
		return 0, err
	}
	broadcast <- string(jsonLog)
	return len(p), nil
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
