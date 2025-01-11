package logstream

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
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
	logrus.SetReportCaller(false)
	logrus.StandardLogger().ReplaceHooks(make(logrus.LevelHooks))

	// Создаем хуки для разных форматов
	consoleHook := newConsoleHook(os.Stdout)
	webSocketHook := newWebSocketHook(broadcast)

	// Добавляем хуки к логгеру
	logrus.AddHook(consoleHook)
	logrus.AddHook(webSocketHook)

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
				logrus.Errorf("❌ Failed to write message: %+v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
