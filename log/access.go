package log

import (
	"net/http"
	"time"
)

func Access(tm time.Time, corrID string, req *http.Request, status int, respSize uint64, duration time.Duration) {
	for _, w := range accessWriters {
		writeAccess(w, tm, corrID, req, status, respSize, duration)
	}
}

func AddAccessWriter(w AccessWriter) {
	if w == nil {
		return
	}

	accessWriters = append(accessWriters, w)
}

type (
	AccessWriter interface {
		Write(tm time.Time, corrID string, req *http.Request, status int, respSize uint64, duration time.Duration)
	}
)

func writeAccess(w AccessWriter, tm time.Time, corrID string, req *http.Request, status int, respSize uint64, duration time.Duration) {
	defer recover()

	w.Write(tm, corrID, req, status, respSize, duration)
}

var accessWriters []AccessWriter
