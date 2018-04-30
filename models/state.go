package models

import "time"

// State is used to communicate the reading state of a file
type State struct {
	Key         string        `json:"key"` 
	Timestamp   time.Time     `json:"timestamp"`
}
