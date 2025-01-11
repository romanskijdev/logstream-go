package logstream

import (
	"github.com/sirupsen/logrus"
	"io"
)

// Хук для вывода в консоль
type consoleHook struct {
	writer io.Writer
}

func newConsoleHook(writer io.Writer) *consoleHook {
	return &consoleHook{writer: writer}
}

func (h *consoleHook) Fire(entry *logrus.Entry) error {
	formatter := &logrus.TextFormatter{} // Используем TextFormatter
	line, err := formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = h.writer.Write(line)
	return err
}

// Хук для отправки в WebSocket
type webSocketHook struct {
	broadcast chan<- string
}

func newWebSocketHook(broadcast chan<- string) *webSocketHook {
	return &webSocketHook{broadcast: broadcast}
}

func (h *webSocketHook) Fire(entry *logrus.Entry) error {
	formatter := &logrus.JSONFormatter{} // Используем JSONFormatter
	line, err := formatter.Format(entry)
	if err != nil {
		return err
	}

	h.broadcast <- string(line) // Отправляем JSON строку
	return nil
}
