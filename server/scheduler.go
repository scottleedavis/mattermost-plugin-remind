package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
)

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
								URL: fmt.Sprintf("/plugins/%s/delete/ephemeral", manifest.ID),
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
								URL: fmt.Sprintf("/plugins/%s/view/ephemeral", manifest.ID),
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
		URL:       fmt.Sprintf("/plugins/%s/dialog", manifest.ID),
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
							Text:  T("button.snooze.10min"),
							Value: "10min",
						},
						{
							Text:  T("button.snooze.20min"),
							Value: "20min",
						},
						{
							Text:  T("button.snooze.30min"),
							Value: "30min",
						},
						{
							Text:  T("button.snooze.40min"),
							Value: "40min",
						},
						{
							Text:  T("button.snooze.50min"),
							Value: "50min",
						},
						{
							Text:  T("button.snooze.1hr"),
							Value: "1hr",
						},
						{
							Text:  T("button.snooze.2hr"),
							Value: "2hr",
						},
						{
							Text:  T("button.snooze.3hr"),
							Value: "3hr",
						},
						{
							Text:  T("button.snooze.4hr"),
							Value: "4hr",
						},
						{
							Text:  T("button.snooze.5hr"),
							Value: "5hr",
						},
						{
							Text:  T("button.snooze.6hr"),
							Value: "6hr",
						},
						{
							Text:  T("button.snooze.7hr"),
							Value: "7hr",
						},
						{
							Text:  T("button.snooze.8hr"),
							Value: "8hr",
						},
						{
							Text:  T("button.snooze.9hr"),
							Value: "9hr",
						},
						{
							Text:  T("button.snooze.10hr"),
							Value: "10hr",
						},
						{
							Text:  T("button.snooze.11hr"),
							Value: "11hr",
						},
						{
							Text:  T("button.snooze.12hr"),
							Value: "12hr",
						},
						{
							Text:  T("button.snooze.13hr"),
							Value: "13hr",
						},
						{
							Text:  T("button.snooze.14hr"),
							Value: "14hr",
						},
						{
							Text:  T("button.snooze.15hr"),
							Value: "15hr",
						},
						{
							Text:  T("button.snooze.16hr"),
							Value: "16hr",
						},
						{
							Text:  T("button.snooze.17hr"),
							Value: "17hr",
						},
						{
							Text:  T("button.snooze.18hr"),
							Value: "18hr",
						},
						{
							Text:  T("button.snooze.19hr"),
							Value: "19hr",
						},
						{
							Text:  T("button.snooze.20hr"),
							Value: "20hr",
						},
						{
							Text:  T("button.snooze.21hr"),
							Value: "21hr",
						},
						{
							Text:  T("button.snooze.22hr"),
							Value: "22hr",
						},
						{
							Text:  T("button.snooze.23hr"),
							Value: "23hr",
						},
						{
							Text:  T("button.snooze.tomorrow"),
							Value: "tomorrow",
						},
						{
							Text:  T("button.snooze.1day"),
							Value: "1day",
						},
						{
							Text:  T("button.snooze.2day"),
							Value: "2day",
						},
						{
							Text:  T("button.snooze.3day"),
							Value: "3day",
						},
						{
							Text:  T("button.snooze.4day"),
							Value: "4day",
						},
						{
							Text:  T("button.snooze.5day"),
							Value: "5day",
						},
						{
							Text:  T("button.snooze.6day"),
							Value: "6day",
						},
						{
							Text:  T("button.snooze.7day"),
							Value: "7day",
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
		if !p.running {
			return
		}
		p.TriggerReminders()
		p.runner()
	}()
}
