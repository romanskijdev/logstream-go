package logstream

import "github.com/sirupsen/logrus"

func (h *webSocketHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *consoleHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
