package main

import "time"

/*
import (
	// "fmt"
	// "github.com/google/uuid"
	// "time"
)
*/

// TODO
func (p *Plugin) ScheduleReminder(request *ReminderRequest) (string, error) {

	user, _ := p.API.GetUserByUsername(request.Username)
	T, _ := p.translation(user)

	if pErr := p.ParseRequest(request); pErr != nil {
		p.API.LogError(pErr.Error())
		return T("exception.response"), nil
	}

	// var when string
	// var target string
	// var message string
	// var useTo bool
	// useTo = false
	// var useToString string
	// if useTo {
	// 	useToString = " to"
	// } else {
	// 	useToString = ""
	// }

	// guid, gErr := uuid.NewRandom()
	// if gErr != nil {
	// 	p.API.LogError("Failed to generate guid")
	// 	return T("exception-response"), nil
	// }

	// target, when, message, pErr := p.ParseRequest(request)
	// if pErr != nil {
	// 	p.API.LogError("parse request failed: " + fmt.Sprintf("%v", pErr))
	// 	return T("exception-response"), nil
	// }

	// request.Reminder.TeamId = request.TeamId
	// request.Reminder.Id = guid.String()
	// request.Reminder.Username = request.Username
	// request.Reminder.Target = target
	// request.Reminder.Message = message
	// request.Reminder.When = when
	// request.Reminder.Occurrences, _ = p.CreateOccurrences(request)

	// p.API.LogDebug(fmt.Sprintf("%v", request.Reminder.Occurrences))

	// p.UpsertReminder(request)

	// if target == "me" {
	// 	target = "you"
	// }

	// response := ":thumbsup: I will remind " + target + useToString + " \"" + request.Reminder.Message + "\" " + when
	// return response, nil

	useTo := strings.HasPrefix(request.Reminder.Message, T("app.reminder.chrono.to"))
	var useToString string
	if useTo {
		useToString = " " + T("app.reminder.chrono.to")
	} else {
		useToString = ""
	}

	request.Reminder.Id = model.NewId()
	request.Reminder.TeamId = request.TeamId
	request.Reminder.UserId = request.UserId
	request.Reminder.Completed = emptyTime.Format(time.RFC3339)

	if cErr := a.createOccurrences(request); cErr != nil {
		mlog.Error(cErr.Error())
		return T(model.REMIND_EXCEPTION_TEXT), nil
	}

	schan := a.Srv.Store.Remind().SaveReminder(&request.Reminder)
	if result := <-schan; result.Err != nil {
		mlog.Error(result.Err.Message)
		return T(model.REMIND_EXCEPTION_TEXT), nil
	}

	if request.Reminder.Target == T("app.reminder.me") {
		request.Reminder.Target = T("app.reminder.you")
	}

	var responseParameters = map[string]interface{}{
		"Target":  request.Reminder.Target,
		"UseTo":   useToString,
		"Message": request.Reminder.Message,
		"When":    a.formatWhen(request.UserId, request.Reminder.When, request.Occurrences[0].Occurrence, false),
	}
	response := T("app.reminder.response", responseParameters)

	return response, nil
}

func (p *Plugin) Run() {

	if !p.running {
		p.running = true
		p.runner()
	}
}

func (p *Plugin) Stop() {
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
