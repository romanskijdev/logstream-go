package logstream

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"runtime"
)

// Хук для вывода в консоль
type consoleHook struct {
	writer *bufferedWriter
}

func newConsoleHook(writer io.Writer) *consoleHook {
	return &consoleHook{writer: newBufferedWriter(writer)}
}

func (h *consoleHook) Fire(entry *logrus.Entry) error {
	formatter := &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", f.File // Возвращаем пустые строки для function и file
		},
	}

	line, err := formatter.Format(entry)
	if err != nil {
		fmt.Println("Failed to format log entry to Console")
		return err
	}
	_, err = h.writer.Write(line)

	if err != nil {
		return err
	}
	return h.writer.Flush() // Сбрасываем буфер после каждой записи
}
