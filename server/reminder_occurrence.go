package main

import (
	"time"
)

type ReminderOccurrence struct {

	Id string

	Username string 
	
	ReminderId string

	Occurrence time.Time

	Snoozed time.Time

	Repeat string
}
