package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/model"
)

type Reminder struct {
	TeamId string

	Id string

	Username string

	Target string

	Message string

	When string

	Occurrences []Occurrence

	Completed time.Time
}

type ReminderRequest struct {
	TeamId string

	Username string

	Payload string

	Reminder Reminder
}

func (p *Plugin) GetReminders(username string) []Reminder {

	user, uErr := p.API.GetUserByUsername(username)

	if uErr != nil {
		p.API.LogError("failed to query user %s", username)
		return []Reminder{}
	}

	bytes, bErr := p.API.KVGet(user.Username)
	if bErr != nil {
		p.API.LogError("failed KVGet %s", bErr)
		return []Reminder{}
	}

	var reminders []Reminder
	err := json.Unmarshal(bytes, &reminders)

	if err != nil {
		p.API.LogError("new reminder " + user.Username)
	} else {
		p.API.LogDebug("existing " + fmt.Sprintf("%v", reminders))
	}

	return reminders
}

func (p *Plugin) UpsertReminder(request *ReminderRequest) error {

	user, uErr := p.API.GetUserByUsername(request.Username)

	if uErr != nil {
		p.API.LogError("failed to query user %s", request.Username)
		return uErr
	}

	bytes, bErr := p.API.KVGet(user.Username)
	if bErr != nil {
		p.API.LogError("failed KVGet %s", bErr)
		return bErr
	}

	var reminders []Reminder
	err := json.Unmarshal(bytes, &reminders)

	if err != nil {
		p.API.LogDebug("new reminder " + user.Username)
	} else {
		p.API.LogDebug("existing " + fmt.Sprintf("%v", reminders))
	}

	reminders = append(reminders, request.Reminder)
	ro, rErr := json.Marshal(reminders)

	if rErr != nil {
		p.API.LogError("failed to marshal reminders %s", user.Username)
		return rErr
	}

	p.API.KVSet(user.Username, ro)

	return nil
}

func (p *Plugin) TriggerReminders() {

	bytes, err := p.API.KVGet(string(fmt.Sprintf("%v", time.Now().Round(time.Second))))

	p.API.LogDebug("*")

	if err != nil {
		p.API.LogError("failed KVGet %s", err)
	} else if string(bytes[:]) == "" {
	} else {

		p.API.LogDebug(string(bytes[:]))

		var reminderOccurrences []Occurrence
		roErr := json.Unmarshal(bytes, &reminderOccurrences)
		if roErr != nil {
			p.API.LogError("Failed to unmarshal reminder occurrences " + fmt.Sprintf("%v", roErr))
			return
		}

		p.API.LogDebug(fmt.Sprintf("%v", reminderOccurrences))

		for _, ReminderOccurrence := range reminderOccurrences {

			user, err := p.API.GetUserByUsername(ReminderOccurrence.Username)

			if err != nil {
				p.API.LogError("failed to query user %s", user.Id)
				continue
			}

			bytes, b_err := p.API.KVGet(user.Username)
			if b_err != nil {
				p.API.LogError("failed KVGet %s", b_err)
				return
			}

			var reminders []Reminder
			uerr := json.Unmarshal(bytes, &reminders)

			if uerr != nil {
				continue
			}

			reminder := p.findReminder(reminders, ReminderOccurrence)

			p.API.LogDebug(fmt.Sprintf("%v", reminder))

			if strings.HasPrefix(reminder.Target, "@") || strings.HasPrefix(reminder.Target, "me") { //@user

				p.API.LogDebug("DM: " + fmt.Sprintf("%v", p.remindUserId) + "__" + fmt.Sprintf("%v", user.Id))
				channel, cErr := p.API.GetDirectChannel(p.remindUserId, user.Id)

				if cErr != nil {
					p.API.LogError("fail to get direct channel ", fmt.Sprintf("%v", cErr))
				} else {
					p.API.LogDebug("got direct channel " + fmt.Sprintf("%v", channel))

					var finalTarget string
					finalTarget = reminder.Target
					if finalTarget == "me" {
						finalTarget = "You"
					} else {
						finalTarget = "@" + user.Username
					}

					if _, err = p.API.CreatePost(&model.Post{
						UserId:    p.remindUserId,
						ChannelId: channel.Id,
						Message:   fmt.Sprintf(finalTarget + " asked me to remind you \"" + reminder.Message + "\"."),
					}); err != nil {
						p.API.LogError(
							"failed to post DM message",
							"user_id", user.Id,
							"error", err.Error(),
						)
					}
				}

			} else { //~ channel

				channel, cErr := p.API.GetChannelByName(reminder.TeamId, strings.Replace(reminder.Target, "~", "", -1), false)

				if cErr != nil {
					p.API.LogError("fail to get channel " + fmt.Sprintf("%v", cErr))
				} else {
					p.API.LogDebug("got channel " + fmt.Sprintf("%v", channel))

					if _, err = p.API.CreatePost(&model.Post{
						UserId:    p.remindUserId,
						ChannelId: channel.Id,
						Message:   fmt.Sprintf("@" + user.Username + " asked me to remind you \"" + reminder.Message + "\"."),
					}); err != nil {
						p.API.LogError(
							"failed to post DM message",
							"user_id", user.Id,
							"error", err.Error(),
						)
					}
				}
			}

		}

	}


	//TODO replace/merge the above with the following code.  ///////////////////////////////

	// t := time.Now().Round(time.Second).Format(time.RFC3339)
	// schan := a.Srv.Store.Remind().GetByTime(t)

	// if result := <-schan; result.Err != nil {
	// 	mlog.Error(result.Err.Message)
	// } else {
	// 	occurrences := result.Data.(model.Occurrences)

	// 	if len(occurrences) == 0 {
	// 		return
	// 	}

	// 	for _, occurrence := range occurrences {

	// 		reminder := model.Reminder{}

	// 		schan = a.Srv.Store.Remind().GetReminder(occurrence.ReminderId)
	// 		if result := <-schan; result.Err != nil {
	// 			continue
	// 		} else {
	// 			reminder = result.Data.(model.Reminder)
	// 		}

	// 		user, _ := a.GetUser(reminder.UserId)
	// 		T, _ := a.translation(user)

	// 		if strings.HasPrefix(reminder.Target, "@") || strings.HasPrefix(reminder.Target, T("app.reminder.me")) {

	// 			channel, cErr := a.GetOrCreateDirectChannel(remindUser.Id, user.Id)
	// 			if cErr != nil {
	// 				continue
	// 			}

	// 			finalTarget := reminder.Target
	// 			if finalTarget == T("app.reminder.me") {
	// 				finalTarget = T("app.reminder.you")
	// 			} else {
	// 				finalTarget = "@" + user.Username
	// 			}

	// 			var messageParameters = map[string]interface{}{
	// 				"FinalTarget": finalTarget,
	// 				"Message":     reminder.Message,
	// 			}

	// 			interactivePost := model.Post{
	// 				ChannelId:     channel.Id,
	// 				PendingPostId: model.NewId() + ":" + fmt.Sprint(model.GetMillis()),
	// 				UserId:        remindUser.Id,
	// 				Message:       T("app.reminder.message", messageParameters),
	// 				Props: model.StringInterface{
	// 					"attachments": []*model.SlackAttachment{
	// 						{
	// 							Actions: []*model.PostAction{
	// 								{
	// 									Integration: &model.PostActionIntegration{
	// 										Context: model.StringInterface{
	// 											"reminderId":   reminder.Id,
	// 											"occurrenceId": occurrence.Id,
	// 											"action":       "complete",
	// 										},
	// 										URL: "mattermost://remind",
	// 									},
	// 									Name: T("app.reminder.update.button.complete"),
	// 									Type: "action",
	// 								},
	// 								{
	// 									Integration: &model.PostActionIntegration{
	// 										Context: model.StringInterface{
	// 											"reminderId":   reminder.Id,
	// 											"occurrenceId": occurrence.Id,
	// 											"action":       "delete",
	// 										},
	// 										URL: "mattermost://remind",
	// 									},
	// 									Name: T("app.reminder.update.button.delete"),
	// 									Type: "action",
	// 								},
	// 								{
	// 									Integration: &model.PostActionIntegration{
	// 										Context: model.StringInterface{
	// 											"reminderId":   reminder.Id,
	// 											"occurrenceId": occurrence.Id,
	// 											"action":       "snooze",
	// 										},
	// 										URL: "mattermost://remind",
	// 									},
	// 									Name: T("app.reminder.update.button.snooze"),
	// 									Type: "select",
	// 									Options: []*model.PostActionOptions{
	// 										{
	// 											Text:  T("app.reminder.update.button.snooze.20min"),
	// 											Value: "20min",
	// 										},
	// 										{
	// 											Text:  T("app.reminder.update.button.snooze.1hr"),
	// 											Value: "1hr",
	// 										},
	// 										{
	// 											Text:  T("app.reminder.update.button.snooze.3hr"),
	// 											Value: "3hrs",
	// 										},
	// 										{
	// 											Text:  T("app.reminder.update.button.snooze.tomorrow"),
	// 											Value: "tomorrow",
	// 										},
	// 										{
	// 											Text:  T("app.reminder.update.button.snooze.nextweek"),
	// 											Value: "nextweek",
	// 										},
	// 									},
	// 								},
	// 							},
	// 						},
	// 					},
	// 				},
	// 			}

	// 			if _, pErr := a.CreatePostAsUser(&interactivePost, false); pErr != nil {
	// 				mlog.Error(fmt.Sprintf("%v", pErr))
	// 			}

	// 			if occurrence.Repeat != "" {
	// 				a.RescheduleOccurrence(&occurrence)
	// 			}

	// 		} else if strings.HasPrefix(reminder.Target, "~") {

	// 			channel, cErr := a.GetChannelByName(
	// 				strings.Replace(reminder.Target, "~", "", -1),
	// 				reminder.TeamId,
	// 				false,
	// 			)

	// 			if cErr != nil {
	// 				mlog.Error(cErr.Message)
	// 				continue
	// 			}

	// 			var messageParameters = map[string]interface{}{
	// 				"FinalTarget": "@" + user.Username,
	// 				"Message":     reminder.Message,
	// 			}

	// 			interactivePost := model.Post{
	// 				ChannelId:     channel.Id,
	// 				PendingPostId: model.NewId() + ":" + fmt.Sprint(model.GetMillis()),
	// 				UserId:        remindUser.Id,
	// 				Message:       T("app.reminder.message", messageParameters),
	// 				Props:         model.StringInterface{},
	// 			}

	// 			if _, pErr := a.CreatePostAsUser(&interactivePost, false); pErr != nil {
	// 				mlog.Error(fmt.Sprintf("%v", pErr))
	// 			}

	// 			if occurrence.Repeat != "" {
	// 				a.RescheduleOccurrence(&occurrence)
	// 			}

	// 		}

	// 	}

	// }

}

func (p *Plugin) DeleteReminders(user *model.User) string {
	T, _ := p.translation(user)
	dErr := p.API.KVDelete(user.Username)
	if dErr != nil {
		p.API.LogError("failed KVDelete %s", dErr)
		return T("exception.response")
	}
	return T("clear.response")
}

func (p *Plugin) findReminder(reminders []Reminder, reminderOccurrence Occurrence) Reminder {
	for _, reminder := range reminders {
		if reminder.Id == reminderOccurrence.ReminderId {
			return reminder
		}
	}
	return Reminder{}
}
