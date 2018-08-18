package main

import (
	"time"
)

type Reminder struct {

    Id string

	Username string

	Target string

	Message string

	Occurrences []ReminderOccurrence

	Completed time.Time
}

type ReminderRequest struct {
	Username string

	Payload string

	Reminder Reminder
}
