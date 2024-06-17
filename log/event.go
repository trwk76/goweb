package log

import "time"

type (
	EventWriter interface {
		Event(tm time.Time, context string, typ EventType, format string, args ...any)
	}

	EventType uint8
)

const (
	EventError EventType = iota
	EventWarning
	EventInfo
	EventTrace
)

func Error(tm time.Time, context string, format string, args ...any) {
	Event(tm, context, EventError, format, args...)
}

func Warning(tm time.Time, context string, format string, args ...any) {
	Event(tm, context, EventWarning, format, args...)
}

func Info(tm time.Time, context string, format string, args ...any) {
	Event(tm, context, EventInfo, format, args...)
}

func Trace(tm time.Time, context string, format string, args ...any) {
	Event(tm, context, EventTrace, format, args...)
}

func Event(tm time.Time, context string, typ EventType, format string, args ...any) {
	if evtWriter == nil {
		return
	}

	defer recover()
	evtWriter.Event(tm, context, typ, format, args...)
}

func SetEventWriter(w EventWriter) {
	evtWriter = w
}

var evtWriter EventWriter
