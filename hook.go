package logstream

import "github.com/sirupsen/logrus"

// fileHook to duplicate messages in the file recorder
type fileHook struct {
	logger *logrus.Logger
}

// WebSocketHook to duplicate messages in the websocket recorder
type webSocketHook struct {
	logger *logrus.Logger
}
