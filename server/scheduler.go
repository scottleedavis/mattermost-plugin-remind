package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/model"
)

const TriggerHostName = "__TRIGGERHOST__"

func (p *Plugin) ScheduleReminder(request *ReminderRequest, channelId string) (*model.Post, error) {

	user, uErr := p.API.GetUserByUsername(request.Username)
	if uErr != nil {
		p.API.LogError(uErr.Error())
		return nil, uErr
	}
	T, _ := p.translation(user)
	location := p.location(user)

	if pErr := p.ParseRequest(request); pErr != nil {
		p.API.LogError(pErr.Error())
		return nil, pErr
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
		return nil, cErr
	}

	if rErr := p.UpsertReminder(request); rErr != nil {
		p.API.LogError(rErr.Error())
		return nil, rErr
	}

	if request.Reminder.Target == T("me") {
		request.Reminder.Target = T("you")
	}

	t := ""
	if len(request.Reminder.Occurrences) > 0 {
		t = request.Reminder.Occurrences[0].Occurrence.In(location).Format(time.RFC3339)
	}
	var responseParameters = map[string]interface{}{
		"Target":  request.Reminder.Target,
		"UseTo":   useToString,
		"Message": request.Reminder.Message,
		"When": p.formatWhen(
			request.Username,
			request.Reminder.When,
			t,
			false,
		),
	}

	return &model.Post{
		ChannelId: channelId,
		UserId:    p.remindUserId,
		Props: model.StringInterface{
			"attachments": []*model.SlackAttachment{
				{
					Text: T("schedule.response", responseParameters),
					Actions: []*model.PostAction{
						{
							Id: model.NewId(),
							Integration: &model.PostActionIntegration{
								Context: model.StringInterface{
									"reminder_id":   request.Reminder.Id,
									"occurrence_id": request.Reminder.Occurrences[0].Id,
									"action":        "delete/ephemeral",
								},
								URL: fmt.Sprintf("%s/plugins/%s/delete/ephemeral", p.URL, manifest.Id),
							},
							Type: model.POST_ACTION_TYPE_BUTTON,
							Name: T("button.delete"),
						},
						{
							Id: model.NewId(),
							Integration: &model.PostActionIntegration{
								Context: model.StringInterface{
									"reminder_id":   request.Reminder.Id,
									"occurrence_id": request.Reminder.Occurrences[0].Id,
									"action":        "view/ephemeral",
								},
								URL: fmt.Sprintf("%s/plugins/%s/view/ephemeral", p.URL, manifest.Id),
							},
							Type: model.POST_ACTION_TYPE_BUTTON,
							Name: T("button.view.reminders"),
						},
					},
				},
			},
		},
	}, nil

}

func (p *Plugin) InteractiveSchedule(triggerId string, user *model.User) {

	T, _ := p.translation(user)

	dialogRequest := model.OpenDialogRequest{
		TriggerId: triggerId,
		URL:       fmt.Sprintf("%s/plugins/%s/dialog", p.URL, manifest.Id),
		Dialog: model.Dialog{
			Title:       T("schedule.reminder"),
			CallbackId:  model.NewId(),
			SubmitLabel: T("button.schedule"),
			Elements: []model.DialogElement{
				{
					DisplayName: T("schedule.message"),
					Name:        "message",
					Type:        "text",
					SubType:     "text",
				},
				{
					DisplayName: T("schedule.target"),
					Name:        "target",
					HelpText:    T("schedule.target.help"),
					Placeholder: "me",
					Type:        "text",
					SubType:     "text",
					Optional:    true,
				},
				{
					DisplayName: T("schedule.time"),
					Name:        "time",
					Type:        "select",
					SubType:     "select",
					Options: []*model.PostActionOptions{
						{
							Text:  T("button.snooze.20min"),
							Value: "20min",
						},
						{
							Text:  T("button.snooze.1hr"),
							Value: "1hr",
						},
						{
							Text:  T("button.snooze.3hr"),
							Value: "3hr",
						},
						{
							Text:  T("button.snooze.tomorrow"),
							Value: "tomorrow",
						},
						{
							Text:  T("button.snooze.nextweek"),
							Value: "nextweek",
						},
					},
				},
			},
		},
	}
	if pErr := p.API.OpenInteractiveDialog(dialogRequest); pErr != nil {
		p.API.LogError("Failed opening interactive dialog " + pErr.Error())
	}
}

func (p *Plugin) Run() {
	p.Stop()
	p.getAndSetLock()
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
		if !p.running && !p.getAndSetLock() {
			return
		}
		p.TriggerReminders()
		p.runner()
	}()
}

func (p *Plugin) getAndSetLock() bool {
	hostname, _ := os.Hostname()
	bytes, bErr := p.API.KVGet(TriggerHostName)
	if bErr != nil {
		p.API.LogError("failed KVGet %s", bErr)
		return false
	}
	if string(bytes) != "" && string(bytes) != hostname {
		return false
	} else if string(bytes) == hostname {
		return true
	}
	p.API.KVSet(TriggerHostName, []byte(hostname))
	return true

}
