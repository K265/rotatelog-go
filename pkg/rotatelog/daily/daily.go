package daily

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/K265/rotatelog-go/pkg/rotatelog"
	"github.com/pkg/errors"
)

const timeFormat = "2006-01-02"

type HandlerImpl struct {
	Prefix       string
	Ext          string
	MaxSize      int64
	KeepDays     int
	MaxIndex     int
	index        int
	lastBasename string
	notifier     rotatelog.Notifier
}

func New(prefix string, ext string, keepDays int, maxSize int64, notifier rotatelog.Notifier) *rotatelog.Writer {
	h := &HandlerImpl{
		Prefix:   prefix,
		Ext:      ext,
		KeepDays: keepDays,
		MaxSize:  maxSize,
		MaxIndex: 1000,
		notifier: notifier,
	}

	return rotatelog.New(h)
}

// GetFilename returns filename like `server.2021-09-07.0.log`
func (h *HandlerImpl) GetFilename(t time.Time) (string, error) {
	basename := h.Prefix + t.Format(timeFormat)
	if basename != h.lastBasename {
		h.lastBasename = basename
		h.index = 0
		go h.Prune(t)
	}

	var filename string
	for {
		filename = basename + "." + strconv.Itoa(h.index) + h.Ext
		fi, err := os.Stat(filename)
		if os.IsNotExist(err) {
			break
		}

		if err == nil && fi.Size() < h.MaxSize {
			break
		}

		h.index += 1
		if h.index > h.MaxIndex {
			return "", errors.Errorf("log rotation index: %v exceed MaxIndex: %v", h.index, h.MaxIndex)
		}
	}

	return filename, nil
}

func (h *HandlerImpl) Prune(now time.Time) {
	dir := path.Dir(h.Prefix)
	prefix := path.Base(h.Prefix)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, file := range files {
		filename := file.Name()
		timeStringStart := len(prefix)
		timeStringEnd := timeStringStart + len(timeFormat)
		if !strings.HasPrefix(filename, prefix) {
			continue
		}

		if len(filename) < timeStringEnd {
			continue
		}

		timeString := filename[timeStringStart:timeStringEnd]
		t, err := time.Parse(timeFormat, timeString)
		if err != nil {
			continue
		}

		if now.Sub(t).Hours() > float64(h.KeepDays)*24 {
			_ = os.Remove(path.Join(dir, filename))
		}

	}
}

func (h *HandlerImpl) OnOpenFile(file *os.File) {
	if h.notifier != nil {
		h.notifier.OnOpenFile(file)
	}
}
