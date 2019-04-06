package main

import (
	"github.com/mattermost/mattermost-server/model"
	"strings"
	"time"
)

func (p *Plugin) ScheduleReminder(request *ReminderRequest) (string, error) {

	user, _ := p.API.GetUserByUsername(request.Username)
	T, _ := p.translation(user)
	location := p.location(user)

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
	request.Reminder.Completed = p.emptyTime

	p.API.LogInfo("1)=====================================>")
	if cErr := p.CreateOccurrences(request); cErr != nil {
		p.API.LogError(cErr.Error())
		return T("exception.response"), nil
	}
	p.API.LogInfo("2)=====================================>")

	if rErr := p.UpsertReminder(request); rErr != nil {
		p.API.LogError(rErr.Error())
		return T("exception.response"), nil
	}
	p.API.LogInfo("3)=====================================>")

	if request.Reminder.Target == T("me") {
		request.Reminder.Target = T("you")
	}

	var responseParameters = map[string]interface{}{
		"Target":  request.Reminder.Target,
		"UseTo":   useToString,
		"Message": request.Reminder.Message,
		"When": p.formatWhen(
			request.Username,
			request.Reminder.When,
			request.Reminder.Occurrences[0].Occurrence.In(location).Format(time.RFC3339),
			false,
		),
	}

	return T("schedule.response", responseParameters), nil
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
