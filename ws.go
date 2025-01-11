package logstream

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

// WebSocketWriter for sending messages to web sockets
type webSocketWriter struct {
	broadcast chan<- map[string]interface{}
}

func newWebSocketWriter(broadcast chan<- map[string]interface{}) *webSocketWriter {
	return &webSocketWriter{broadcast: broadcast}
}

func (w *webSocketWriter) Write(p []byte) (int, error) {
	var logEntry map[string]interface{}
	if err := json.Unmarshal(p, &logEntry); err != nil {
		return 0, err
	}
	w.broadcast <- logEntry
	return len(p), nil
}

func (h *webSocketHook) Fire(entry *logrus.Entry) error {
	h.logger.WithFields(entry.Data).Log(entry.Level, entry.Message)
	return nil
}
