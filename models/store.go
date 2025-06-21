package models

import "time"

type Value struct {
	Value  string
	Expiry time.Time
}
