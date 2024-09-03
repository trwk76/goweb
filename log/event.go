package log

import (
	"fmt"
	"time"
)

func Fatal(tm time.Time, context string, format string, args ...any) {
	Event(tm, context, EventFatal, format, args...)
}

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

func Event(tm time.Time, context string, level EventLevel, format string, args ...any) {
	if len(eventWriters) < 1 {
		return
	}

	text := fmt.Sprintf(format, args...)

	for _, w := range eventWriters {
		w.Write(tm, context, level, text)
	}
}

func AddEventWriter(maxLevel EventLevel, writer EventWriter) {
	if writer == nil {
		return
	}

	eventWriters = append(eventWriters, eventWriter{
		maxLevel: maxLevel,
		writer:   writer,
	})
}

type (
	EventWriter interface {
		Write(tm time.Time, context string, level EventLevel, text string)
	}

	EventLevel uint8

	eventWriter struct {
		maxLevel EventLevel
		writer   EventWriter
	}
)

const (
	EventFatal EventLevel = iota
	EventError
	EventWarning
	EventInfo
	EventTrace
)

func (l EventLevel) String() string {
	switch l {
	case EventFatal:
		return "FATAL"
	case EventError:
		return "ERROR"
	case EventWarning:
		return "WARNING"
	case EventInfo:
		return "INFO"
	}

	return "TRACE"
}

func (e eventWriter) Write(tm time.Time, context string, level EventLevel, text string) {
	if level > e.maxLevel {
		return
	}

	defer recover()

	e.writer.Write(tm, context, level, text)
}

var eventWriters []eventWriter
