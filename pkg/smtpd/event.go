package smtpd

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// EventLevel type gives the level of an event.
type EventLevel string

const (
	// EventLevelAll constant
	EventLevelAll EventLevel = "ALL"
	// EventLevelTrace constant
	EventLevelTrace EventLevel = "TRC"
	// EventLevelDebug constant
	EventLevelDebug EventLevel = "DBG"
	// EventLevelInfo constant
	EventLevelInfo EventLevel = "INF"
	// EventLevelWarn constant
	EventLevelWarn EventLevel = "WRN"
	// EventLevelError constant
	EventLevelError EventLevel = "ERR"
	// EventLevelFatal constant
	EventLevelFatal EventLevel = "FTL"
	// EventLevelNone constant
	EventLevelNone EventLevel = "NON"
)

// Event type holds information about an event that can be logged.
type Event map[string]interface{}

var (
	// EventFieldLevel is the name of the time field.
	EventFieldLevel = "level"
	// EventFieldTime is the name of the time field.
	EventFieldTime = "time"
	// EventFieldMessage is the name of the message field.
	EventFieldMessage = "message"
)

func event(level EventLevel, message string, base Event) Event {
	base[EventFieldLevel] = level
	base[EventFieldTime] = time.Now().Unix()
	base[EventFieldMessage] = message
	return base
}

type eventProducer func(string, Event) Event

func eTrace(message string, base Event) Event {
	return event(EventLevelTrace, message, base)
}

func eDebug(message string, base Event) Event {
	return event(EventLevelDebug, message, base)
}

func eInfo(message string, base Event) Event {
	return event(EventLevelInfo, message, base)
}

func eWarn(message string, base Event) Event {
	return event(EventLevelWarn, message, base)
}

func eError(message string, base Event) Event {
	return event(EventLevelError, message, base)
}

func eFatal(message string, base Event) Event {
	return event(EventLevelFatal, message, base)
}

func (e Event) String() string {
	var result strings.Builder
	fmt.Fprintf(&result, "%v %v %v", time.Unix(e[EventFieldTime].(int64), 0).Format(time.RFC3339), e[EventFieldLevel], e[EventFieldMessage])

	if len(e) > 3 {
		result.WriteString(" {")
	}

	keys := make([]string, len(e))

	i := 0
	for k := range e {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	for _, key := range keys {
		if key != EventFieldTime && key != EventFieldLevel && key != EventFieldMessage {
			fmt.Fprintf(&result, " %v=\"%v\"", key, e[key])
		}
	}

	if len(e) > 3 {
		result.WriteString(" }")
	}

	return result.String()
}
