package main

import (
	"fmt"
	"time"
	"strings"
	"errors"

	"github.com/google/uuid"
)

func (p *Plugin) Run() {

	if !p.running {
		p.running = true
		p.runner()
	}
}

func (p *Plugin) stop() {
	p.running = false
}

func (p *Plugin) runner() {

	go func() {
		<-time.NewTimer(time.Second).C
		p.TriggerReminders()
		if !p.running {
			return
		}
		p.runner()
	}()
}

func (p *Plugin) ScheduleReminder(request ReminderRequest) (string, error) {

	var when string
	var target string
	var message string
	var useTo bool
	useTo = false
	var useToString string
	if useTo {
		useToString = " to"
	} else {
		useToString = ""
	}

	guid, gerr := uuid.NewRandom()
	if gerr != nil {
		p.API.LogError("Failed to generate guid")
	}

	target, when, message, perr := p.parseRequest(request)
	if perr != nil {
		return ExceptionText, nil
	}

	request.Reminder.TeamId = request.TeamId
	request.Reminder.Id = guid.String()
	request.Reminder.Username = request.Username
	request.Reminder.Target = target
	request.Reminder.Message = message
	request.Reminder.When = when
	request.Reminder.Occurrences = p.CreateOccurrences(request)

	// TODO REMOVE THIS LATER
	p.API.KVDelete(request.Username)
	////////////

	p.UpsertReminder(request)

	response := ":thumbsup: I will remind " + target + useToString + " \"" + request.Reminder.Message + "\" " + when;
	return response, nil
}

func (p *Plugin) parseRequest(request ReminderRequest) (string, string, string, error) {

	commandSplit := strings.Split(request.Payload, " ")

	p.API.LogError("parseRequest " + fmt.Sprintf("%v", request))
	p.API.LogError(request.Payload)

	if strings.HasPrefix(request.Payload, "me") ||
		strings.HasPrefix(request.Payload, "~") ||
		strings.HasPrefix(request.Payload, "@") {

		p.API.LogError("found target")

		var message string
		var when string
		var firstIndex int
		var lastIndex int

		firstIndex = strings.Index(request.Payload, "\"")
		lastIndex = strings.LastIndex(request.Payload, "\"")

		if firstIndex > -1 && lastIndex > -1 && firstIndex != lastIndex {

			message := request.Payload[firstIndex:lastIndex]
			when = strings.Replace(request.Payload, message, "", -1)
			when = strings.Replace(when, commandSplit[1], "", -1)
			p.API.LogError("quotes when " + fmt.Sprintf("%v", firstIndex) + " " + fmt.Sprintf("%v", lastIndex) + " " + when)

			return commandSplit[0], when, message, nil
		}

		p.API.LogError("no quotes when " + fmt.Sprintf("%v", firstIndex) + " " + fmt.Sprintf("%v", lastIndex))

		message = "foo"

		////////
		// TODO determine when
		/////////

		return commandSplit[0], "in 2 seconds", message, nil
	}
	err := errors.New("Unrecognized Target")

	return "", "", "", err
}
