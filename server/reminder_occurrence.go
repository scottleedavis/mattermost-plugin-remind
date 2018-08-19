package main

import (
	"fmt"
	"time"
	"encoding/json"
	"github.com/google/uuid"
)

type ReminderOccurrence struct {
	Id string

	Username string

	ReminderId string

	Occurrence time.Time

	Snoozed time.Time

	Repeat string
}

func (p *Plugin) CreateOccurrences(request ReminderRequest) ([]ReminderOccurrence) {

	var ReminderOccurrences []ReminderOccurrence

	// switch the when patterns

	// handle seconds as proof of concept

	guid, gerr := uuid.NewRandom()
	if gerr != nil {
		p.API.LogError("Failed to generate guid")
		return []ReminderOccurrence{}
	}

	occurrence := time.Now().Round(time.Second).Add(time.Second * time.Duration(5))
	reminderOccurrence := ReminderOccurrence{guid.String(), request.Username, request.Reminder.Id, occurrence, time.Time{}, ""}
	ReminderOccurrences = append(ReminderOccurrences, reminderOccurrence)

	p.upsertOccurrence(reminderOccurrence)

	return ReminderOccurrences
}

func (p *Plugin) upsertOccurrence(reminderOccurrence ReminderOccurrence) {

	bytes, err := p.API.KVGet(string(fmt.Sprintf("%v", reminderOccurrence.Occurrence)))
	if err != nil {
		p.API.LogError("failed KVGet %s", err)
		return
	}

	var reminderOccurrences []ReminderOccurrence

	roErr := json.Unmarshal(bytes, &reminderOccurrences)
	if roErr != nil {
		p.API.LogError("new occurrence " + string(fmt.Sprintf("%v", reminderOccurrence.Occurrence)))
	} else {
		p.API.LogError("existing " + fmt.Sprintf("%v", reminderOccurrences))
	}

	reminderOccurrences = append(reminderOccurrences, reminderOccurrence)
	ro, __ := json.Marshal(reminderOccurrences)

	if __ != nil {
		p.API.LogError("failed to marshal reminderOccurrences %s", reminderOccurrence.Id)
		return
	}

	p.API.KVSet(string(fmt.Sprintf("%v", reminderOccurrence.Occurrence)), ro)

}
