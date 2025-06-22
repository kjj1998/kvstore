package models

import "time"

type Value struct {
	Value  string
	Expiry time.Time
}

type WALEntry struct {
	Command string
	Key     string
	Value   string
	Expiry  string
}
