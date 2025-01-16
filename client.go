package logstream

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

// upgrader for websockets
var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// clients connected to websockets
var clients = make(map[*websocket.Conn]bool)

// channel for transmitting log messages
var broadcast = make(chan map[string]interface{})

// LoggerConfig structure for storing logger configuration
type LoggerConfig struct {
	DateFormat string
	FilePath   *string
}

func InitLoggerClient(dateFormat *string, filePath *string) *LoggerConfig {
	if dateFormat == nil {
		defaultDate := "2006-01-02"
		dateFormat = &defaultDate
	}
	return &LoggerConfig{
		DateFormat: *dateFormat,
		FilePath:   filePath,
	}
}

// InitLogger logger init
func (config *LoggerConfig) InitLogger() error {
	// install a standard logger for console output
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: fmt.Sprintf("%s 15:04:05", config.DateFormat),
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", f.File
		},
	})
	logrus.SetOutput(os.Stdout)

	// create a WebSocket logger
	webSocketLogger := logrus.New()
	webSocketLogger.SetFormatter(&logrus.JSONFormatter{})
	webSocketLogger.SetOutput(newWebSocketWriter(broadcast))

	// installing a hook for duplicating messages in the WebSocket logger
	logrus.AddHook(&webSocketHook{logger: webSocketLogger})

	if config.FilePath != nil {
		// ensure directory exists
		err := os.MkdirAll(filepath.Dir(*config.FilePath), 0755)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to create log directory: %v", err))
		}

		// creating logs file (if not exist)
		fileLogger := logrus.New()
		fileLogger.SetFormatter(&logrus.JSONFormatter{})
		file, err := os.OpenFile(*config.FilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return errors.New(fmt.Sprintf("failed to open log file: %s. %v", *config.FilePath, err))
		}
		fileLogger.SetOutput(newFileWriter(file, config.DateFormat))

		// install a hook for duplicating messages into a logger file
		logrus.AddHook(&fileHook{logger: fileLogger})
	}

	// running a log message handler in a separate goroutine
	go handleLogMessages()
	return nil
}

// HandleConnections handling connections
func (config *LoggerConfig) HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		logrus.Errorf("❌ Failed to upgrade to websocket: %+v", err)
		return
	}
	defer ws.Close()

	clients[ws] = true

	// sending all previous logs to the new client
	if config.FilePath != nil {
		err := config.sendPreviousLogs(ws, *config.FilePath)
		if err != nil {
			logrus.Error(err)
			return
		}
	}

	for {
		var msg map[string]interface{}
		err := ws.ReadJSON(&msg)
		if err != nil {
			delete(clients, ws)
			break
		}
	}
}
