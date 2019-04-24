package main

import (
	"os"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/model"
)

const TriggerHostName = "__TRIGGERHOST__"

func (p *Plugin) ScheduleReminder(request *ReminderRequest) (string, error) {

	user, uErr := p.API.GetUserByUsername(request.Username)
	if uErr != nil {
		p.API.LogError(uErr.Error())
		return "", uErr
	}
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

	if cErr := p.CreateOccurrences(request); cErr != nil {
		p.API.LogError(cErr.Error())
		return T("exception.response"), nil
	}

	if rErr := p.UpsertReminder(request); rErr != nil {
		p.API.LogError(rErr.Error())
		return T("exception.response"), nil
	}

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

	hostname, _ := os.Hostname()
	bytes, bErr := p.API.KVGet(TriggerHostName)
	if bErr != nil {
		p.API.LogError("failed KVGet %s", bErr)
		return
	}
	if string(bytes) != "" && string(bytes) != hostname {
		return
	}
	p.API.KVSet(TriggerHostName, []byte(hostname))

	if !p.running {
		p.running = true
		p.runner()
	}

}

func (p *Plugin) Stop() {
	p.API.KVSet(TriggerHostName, []byte(""))
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
