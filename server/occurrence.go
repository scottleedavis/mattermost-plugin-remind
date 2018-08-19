package main

import (
	"fmt"
	"time"
	"encoding/json"
	"strings"
	"errors"
	"strconv"
	"github.com/google/uuid"
)

type Occurrence struct {
	Id string

	Username string

	ReminderId string

	Occurrence time.Time

	Snoozed time.Time

	Repeat string
}

func (p *Plugin) CreateOccurrences(request ReminderRequest) ([]Occurrence, error) {

	p.API.LogDebug("CreateOccurrences");

	if strings.HasPrefix(request.Reminder.When, "in") {

		p.API.LogDebug(request.Reminder.When)

		occurrences, inErr := p.in(request.Reminder.When)
		if inErr != nil {
			return []Occurrence{}, inErr
		}

		guid, gErr := uuid.NewRandom()
		if gErr != nil {
			p.API.LogError("failed to generate guid")
			return []Occurrence{}, gErr
		}

		for _, o := range occurrences {

			reminderOccurrence := Occurrence{guid.String(), request.Username, request.Reminder.Id, o, time.Time{}, ""}

			p.API.LogDebug("occurrence "+fmt.Sprintf("%v", reminderOccurrence))

			request.Reminder.Occurrences = append(request.Reminder.Occurrences, reminderOccurrence)
			p.upsertOccurrence(reminderOccurrence)

		}

		return request.Reminder.Occurrences, nil

	}

	// TODO handle the other when prefix's

	return []Occurrence{}, errors.New("unable to create occurrences")
}

func (p *Plugin) upsertOccurrence(reminderOccurrence Occurrence) {

	bytes, err := p.API.KVGet(string(fmt.Sprintf("%v", reminderOccurrence.Occurrence)))
	if err != nil {
		p.API.LogError("failed KVGet %s", err)
		return
	}

	var reminderOccurrences []Occurrence
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

func (p *Plugin) in(when string) (times []time.Time, err error) {

	whenSplit := strings.Split(when, " ")
	value := whenSplit[1]
	units := whenSplit[len(whenSplit)-1]

	p.API.LogDebug("whenSplit: "+fmt.Sprintf("%v",whenSplit))
	p.API.LogDebug("value: "+fmt.Sprintf("%v",value))
	p.API.LogDebug("units: "+fmt.Sprintf("%v",units))

	switch units {
	case "seconds", "second", "sec", "s":
		i, _ := strconv.Atoi(value)
		occurrence := time.Now().Round(time.Second).Add(time.Second * time.Duration(i))
		times = append(times, occurrence)
		p.API.LogDebug("occurrence: " + fmt.Sprintf("%v", occurrence))
		p.API.LogDebug("times: " + fmt.Sprintf("%v", times))

	//TODO handle the other units

	default:
		return nil, errors.New("could not format 'in'")
	}

	return times, nil

}
