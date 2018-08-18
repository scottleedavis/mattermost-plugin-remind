package main

import (
	"time"
	// "strings"

	"github.com/google/uuid"
)

func (p *Plugin) runner() {

    go func() {
		<-time.NewTimer(time.Second).C
		p.triggerReminders()
		if !p.running {
			return
		}
		p.runner()
	}()
}

func (p *Plugin) run() {
	
	if !p.running {
		p.running = true
		p.runner()
	}
}

func (p *Plugin) stop() {
	p.running = false
}


func (p *Plugin) parseRequest(request ReminderRequest) (string, string, string, error) {
	return "me", "in 2 seconds", "this is an end to end test.", nil
}

func (p *Plugin) scheduleReminder(request ReminderRequest) (string, error) {

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
	request.Reminder.Occurrences = p.createOccurrences(request)

	// TODO REMOVE THIS LATER
	p.API.KVDelete(request.Username)
	////////////

	p.upsertReminder(request)

	response := ":thumbsup: I will remind " + target + useToString + " \"" + request.Reminder.Message + "\" " + when;
   
	return response, nil
}