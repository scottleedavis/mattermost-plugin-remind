package main

import (
	"time"
)

type ReminderOccurrence struct {

	Id string

	ReminderId string

	Occurrence time.Time

	Snoozed time.Time

	Repeat string
}
