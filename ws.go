package logstream

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

// Хук для отправки в WebSocket
type webSocketHook struct {
	broadcast chan<- string
}

func newWebSocketHook(broadcast chan<- string) *webSocketHook {
	return &webSocketHook{broadcast: broadcast}
}

func (h *webSocketHook) Fire(entry *logrus.Entry) error {
	formatter := &logrus.JSONFormatter{}
	line, err := formatter.Format(entry)
	if err != nil {
		fmt.Println("Failed to format log entry to JSON")

		return err // Возвращаем ошибку, если форматирование не удалось
	}

	h.broadcast <- string(line)
	return nil // Возвращаем nil только если все прошло успешно
}
