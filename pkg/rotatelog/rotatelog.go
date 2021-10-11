package rotatelog

import (
	"os"
	"sync"
	"time"
)

type Notifier interface {
	// OnOpenFile is called when opened new files
	OnOpenFile(file *os.File)
}

type Handler interface {
	// GetFilename returns filename by current time
	GetFilename(t time.Time) (string, error)
}

type Writer struct {
	mutex        sync.RWMutex
	file         *os.File
	lastFilename string
	handler      Handler
}

func New(handler Handler) *Writer {
	return &Writer{handler: handler}
}

func (w *Writer) Write(p []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	now := time.Now()
	filename, err := w.handler.GetFilename(now)
	if err != nil {
		return 0, err
	}

	if filename != w.lastFilename {
		if w.file != nil {
			if err := w.file.Close(); err != nil {
				return 0, err
			}
		}

		w.file, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return 0, err
		}

		if notifier, ok := w.handler.(Notifier); ok {
			notifier.OnOpenFile(w.file)
		}
	}

	return w.file.Write(p)
}
