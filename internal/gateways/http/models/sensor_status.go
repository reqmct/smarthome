package models

import "time"

type SensorStatus struct {
	Timestamp time.Time `json:"timestamp"`
	Payload   int64     `json:"payload"`
}
