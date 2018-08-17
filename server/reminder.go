package main

import (
	"time"
)

type Reminder struct {

	Username string

	Target string

	Message string

	Occurrences []time.Time
	
	Completed time.Time
}

type ReminderRequest struct {
	Username string

	Payload string
}
