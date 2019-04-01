package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/model"
)

type Occurrence struct {
	Id string

	Username string

	ReminderId string

	Occurrence time.Time

	Snoozed time.Time

	Repeat string
}

func (p *Plugin) CreateOccurrences(request *ReminderRequest) error {

	user, _ := p.API.GetUserByUsername(request.Username)
	_, locale := p.translation(user)

	switch locale {
	case "en":
		return p.createOccurrencesEN(request)
	default:
		return p.createOccurrencesEN(request)
	}

}

func (p *Plugin) createOccurrencesEN(request *ReminderRequest) error {

	user, _ := p.API.GetUserByUsername(request.Username)
	T, _ := p.translation(user)

	if strings.HasPrefix(request.Reminder.When, T("in")) {
		if occurrences, inErr := p.in(request.Reminder.When, user); inErr != nil {
			return inErr
		} else {
			return p.addOccurrences(request, occurrences)
		}
	}

	if strings.HasPrefix(request.Reminder.When, T("at")) {
		if occurrences, inErr := p.at(request.Reminder.When, user); inErr != nil {
			return inErr
		} else {
			return p.addOccurrences(request, occurrences)
		}
	}

	if strings.HasPrefix(request.Reminder.When, T("on")) {
		if occurrences, inErr := p.on(request.Reminder.When, user); inErr != nil {
			return inErr
		} else {
			return p.addOccurrences(request, occurrences)
		}
	}

	if strings.HasPrefix(request.Reminder.When, T("every")) {
		if occurrences, inErr := p.every(request.Reminder.When, user); inErr != nil {
			return inErr
		} else {
			return p.addOccurrences(request, occurrences)
		}
	}

	if occurrences, freeErr := p.freeForm(request.Reminder.When, user); freeErr != nil {
		return freeErr
	} else {
		return p.addOccurrences(request, occurrences)
	}

}

func (p *Plugin) addOccurrences(request *ReminderRequest, occurrences []time.Time) error {

	user, _ := p.API.GetUserByUsername(request.Username)
	T, _ := p.translation(user)

	for _, o := range occurrences {

		repeat := ""

		if p.isRepeating(request) {
			repeat = request.Reminder.When
			if strings.HasPrefix(request.Reminder.Target, "@") &&
				request.Reminder.Target != T("app.reminder.me") {

				rUser, _ := p.API.GetUserByUsername(request.Username)

				if tUser, tErr := p.API.GetUserByUsername(request.Reminder.Target[1:]); tErr != nil {
					return tErr
				} else {
					if rUser.Id != tUser.Id {
						return errors.New("repeating reminders for another user not permitted")
					}
				}

			}
		}

		occurrence := &Occurrence{
			Id:         model.NewId(),
			Username:   request.Username,
			ReminderId: request.Reminder.Id,
			Repeat:     repeat,
			Occurrence: o,
			Snoozed:    p.emptyTime,
		}

		request.Request.Occurrences = p.upsertOccurrence(occurrence)
	}

	return nil
}

func (p *Plugin) isRepeating(request *ReminderRequest) bool {

	user, _ := p.API.GetUserByUsername(request.Username)
	T, _ := p.translation(user)

	return strings.Contains(request.Reminder.When, T("every")) ||
		strings.Contains(request.Reminder.When, T("sundays")) ||
		strings.Contains(request.Reminder.When, T("mondays")) ||
		strings.Contains(request.Reminder.When, T("tuesdays")) ||
		strings.Contains(request.Reminder.When, T("wednesdays")) ||
		strings.Contains(request.Reminder.When, T("thursdays")) ||
		strings.Contains(request.Reminder.When, T("fridays")) ||
		strings.Contains(request.Reminder.When, T("saturdays"))

}

func (p *Plugin) upsertOccurrence(occurrence *Occurrence) []Occurrence {

	bytes, err := p.API.KVGet(string(fmt.Sprintf("%v", occurrence.Occurrence)))
	if err != nil {
		p.API.LogError("failed KVGet %s", err)
		return nil
	}

	var occurrences []Occurrence
	roErr := json.Unmarshal(bytes, &occurrences)
	if roErr != nil {
		p.API.LogDebug("new occurrence " + string(fmt.Sprintf("%v", occurrence.Occurrence)))
	} else {
		p.API.LogDebug("existing " + fmt.Sprintf("%v", occurrences))
	}

	occurrences = append(occurrences, occurrence)
	ro, __ := json.Marshal(occurrences)

	if __ != nil {
		p.API.LogError("failed to marshal reminderOccurrences %s", occurrence.Id)
		return
	}

	p.API.KVSet(string(fmt.Sprintf("%v", occurrence.Occurrence)), ro)

	return occurrences

}
