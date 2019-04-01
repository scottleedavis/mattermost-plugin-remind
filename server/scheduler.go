package main

import (
	// "strings"
	"time"
	// "github.com/mattermost/mattermost-server/model"
)

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

	useTo := strings.HasPrefix(request.Reminder.Message, T("to"))
	var useToString string
	if useTo {
		useToString = " " + T("to")
	} else {
		useToString = ""
	}

	request.Reminder.Id = model.NewId()
	request.Reminder.TeamId = request.TeamId
	request.Reminder.Username = request.Username
	request.Reminder.Completed = p.emptyTime.Format(time.RFC3339)

	if cErr := p.CreateOccurrences(request); cErr != nil {
		p.API.LogError(cErr.Error())
		return T(model.REMIND_EXCEPTION_TEXT), nil
	}

	// p.UpsertReminder(request)

	// if target == "me" {
	// 	target = "you"
	// }

	// response := ":thumbsup: I will remind " + target + useToString + " \"" + request.Reminder.Message + "\" " + when
	// return response, nil

	// TODO ////////////////////////////////////////////////////////////////////////////////////

	// schan := a.Srv.Store.Remind().SaveReminder(&request.Reminder)
	// if result := <-schan; result.Err != nil {
	// 	mlog.Error(result.Err.Message)
	// 	return T(model.REMIND_EXCEPTION_TEXT), nil
	// }

	// if request.Reminder.Target == T("me") {
	// 	request.Reminder.Target = T("you")
	// }

	// var responseParameters = map[string]interface{}{
	// 	"Target":  request.Reminder.Target,
	// 	"UseTo":   useToString,
	// 	"Message": request.Reminder.Message,
	// 	"When":    a.formatWhen(request.UserId, request.Reminder.When, request.Occurrences[0].Occurrence, false),
	// }
	// response := T("response", responseParameters)

	// return response, nil

	return "this is a test", nil
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
