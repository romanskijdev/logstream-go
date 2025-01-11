package logstream

import "github.com/sirupsen/logrus"

func (h *webSocketHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *fileHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
