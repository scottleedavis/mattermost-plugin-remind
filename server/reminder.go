package main

import (
	"time"
)

type Reminder struct {

	UserId string

	Target string

	Username string

	Message string

	Occurrences []string

	Completed time.Time
}
