package models

import (
	//"sync"
	//"time"
	//"github.com/elastic/beats/libbeat/logp"
)

 
type States struct {
	// states store
	states []State
}

// NewStates generates a new states registry.
func NewStates() *States {
	return &States{
		states: nil,
	}
}
 