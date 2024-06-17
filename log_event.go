package web

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

func LogError(tm time.Time, context string, format string, args ...any) {
	LogEvent(tm, context, EventError, format, args...)
}

func LogWarning(tm time.Time, context string, format string, args ...any) {
	LogEvent(tm, context, EventWarning, format, args...)
}

func LogInfo(tm time.Time, context string, format string, args ...any) {
	LogEvent(tm, context, EventInfo, format, args...)
}

func LogTrace(tm time.Time, context string, format string, args ...any) {
	LogEvent(tm, context, EventTrace, format, args...)
}

func LogEvent(tm time.Time, context string, typ EventType, format string, args ...any) {
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
