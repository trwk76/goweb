package log

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func NewTextAccessWriter(w io.Writer) AccessWriter {
	if w == nil {
		return nil
	}

	return TextAccessWriter{w: w}
}

func NewTextEventWriter(w io.Writer) EventWriter {
	if w == nil {
		return nil
	}

	return TextEventWriter{w: w}
}

type (
	TextAccessWriter struct {
		w io.Writer
	}

	TextEventWriter struct {
		w io.Writer
	}
)

func (w TextAccessWriter) Write(tm time.Time, corrID string, r *http.Request, status int, respSize uint64, duration time.Duration) {
	fmt.Fprintf(
		w.w,
		"%s (%s) %s %s?%s: %d (%d bytes): %dmsecs\n",
		tm.Format(time.DateTime),
		corrID,
		r.Method,
		r.URL.Path,
		r.URL.Query().Encode(),
		status,
		respSize,
		uint64(duration/time.Millisecond),
	)
}

func (w TextEventWriter) Write(tm time.Time, context string, level EventLevel, text string) {
	fmt.Fprintf(
		w.w,
		"%s (%s) [%s]: %s\n",
		tm.Format(time.DateTime),
		context,
		level.String(),
		text,
	)
}

var (
	_ AccessWriter = TextAccessWriter{}
	_ EventWriter  = TextEventWriter{}
)
