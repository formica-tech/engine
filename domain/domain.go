package domain

import (
	"encoding/json"
	"time"
)

type Interval struct {
	Start time.Time
	End   time.Time
}

type SignalPayload map[string]interface{}

type EntitySignal struct {
	ID        string
	EntityID  string
	Event     string
	Payload   SignalPayload
	Timestamp time.Time
}

func (ds *EntitySignal) GetPayload() string {
	b, err := json.Marshal(ds.Payload)
	if err != nil {
		return "{}"
	}

	return string(b)
}

type EntitySignalRepository interface {
	BatchSave(signals []EntitySignal) error
}

type EntitySignalFilter struct {
	EntityID []string
	SourceID []string
	Event    []string
	Interval Interval
}
