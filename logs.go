package logstream

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
	"time"
)

// fileWriter to write logs to a file
type fileWriter struct {
	file       *os.File
	mutex      sync.Mutex
	dateFormat string
}

func newFileWriter(file *os.File, dateFormat string) *fileWriter {
	return &fileWriter{
		file:       file,
		dateFormat: dateFormat,
	}
}

// sending all previous logs for today to a new client
func (config *LoggerConfig) sendPreviousLogs(ws *websocket.Conn, filePath string) error {
	logsByDate, err := readLogsFromFile(filePath)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to read log file: %+v", err))
	}

	for _, logs := range logsByDate {
		for _, logEntry := range logs {
			err := ws.WriteJSON(logEntry)
			if err != nil {
				return errors.New(fmt.Sprintf("failed to write log to WebSocket: %+v", err))
			}
		}
	}

	return nil
}

// reading logs from a file and grouping by dates
func readLogsFromFile(filepath string) (map[string][]map[string]interface{}, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var logsByDate = make(map[string][]map[string]interface{})
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&logsByDate)
	if err != nil && err != io.EOF {
		return nil, err
	}

	return logsByDate, nil
}

// log message handler
func handleLogMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				logrus.Errorf("failed to write message: %+v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func (w *fileWriter) Write(p []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// reading existing logs
	logsByDate, err := readLogsFromFile(w.file.Name())
	if err != nil {
		return 0, err
	}

	// adding a new log to the corresponding array by date
	var newLog map[string]interface{}
	err = json.Unmarshal(p, &newLog)
	if err != nil {
		return 0, err
	}

	date := time.Now().Format(w.dateFormat) // use the date format passed to InitLogger
	logsByDate[date] = append(logsByDate[date], newLog)

	// overwriting a file with updated logs
	w.file.Seek(0, 0)
	err = w.file.Truncate(0)
	if err != nil {
		return 0, err
	}
	err = json.NewEncoder(w.file).Encode(logsByDate)
	if err != nil {
		return 0, err
	}

	return len(p), nil
}

func (h *fileHook) Fire(entry *logrus.Entry) error {
	h.logger.WithFields(entry.Data).Log(entry.Level, entry.Message)
	return nil
}
