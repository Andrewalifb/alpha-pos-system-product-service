package dto

import (
	"encoding/json"
	"time"
)

const (
	MESSAGE_FAILED_GET_DATA_FROM_BODY = "failed get data from body"
)

type JSONTime struct {
	time.Time
}

func (t *JSONTime) UnmarshalJSON(b []byte) error {
	var timestamp struct {
		Seconds int64 `json:"seconds"`
		Nanos   int64 `json:"nanos"`
	}
	if err := json.Unmarshal(b, &timestamp); err != nil {
		return err
	}
	t.Time = time.Unix(timestamp.Seconds, timestamp.Nanos)
	return nil
}

func (t JSONTime) MarshalJSON() ([]byte, error) {
	timestamp := struct {
		Seconds int64 `json:"seconds"`
		Nanos   int64 `json:"nanos"`
	}{
		Seconds: t.Unix(),
		Nanos:   int64(t.Nanosecond()),
	}
	return json.Marshal(timestamp)
}
