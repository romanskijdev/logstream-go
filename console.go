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
		logrus.WithError(err).Error("Failed to format log entry to Console")
		return err
	}
	_, err = h.writer.Write(line)
	return err
}
