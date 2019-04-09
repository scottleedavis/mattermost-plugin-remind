package main

import (
	"encoding/json"
	"fmt"
	// "io/ioutil"
	"net/http"
	// "net/http/httputil"
	// "bytes"
	"time"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

// ActionContext passed from action buttons
type ActionContext struct {
	ReminderID     string `json:"reminder_id"`
	OccurrenceID   string `json:"occurrence_id"`
	Action         string `json:"action"`
	SelectedOption string `json:"selected_option"`
}

// Action type for decoding action buttons
type Action struct {
	UserID  string         `json:"user_id"`
	PostID  string         `json:"post_id"`
	Context *ActionContext `json:"context"`
}

// TODO use ephemeral posts when they are no longer not experimental
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {

	var action *Action
	json.NewDecoder(r.Body).Decode(&action)

	switch action.Context.Action {
	case "complete":
		p.handleComplete(w, r, action)
	case "complete/list":
		p.handleCompleteList(w, r, action)
	case "view/complete/list":
		p.handleViewCompleteList(w, r, action)
	case "delete":
		p.handleDelete(w, r, action)
	case "delete/list":
		p.handleDeleteList(w, r, action)
	case "delete/complete/list":
		p.handleDeleteCompleteList(w, r, action)
	case "snooze":
		p.handleSnooze(w, r, action)
	case "snooze/list":
		p.handleSnoozeList(w, r, action)
	case "close/list":
		p.handleCloseList(w, r, action)
	default:
		response := &model.PostActionIntegrationResponse{}
		writePostActionIntegrationResponseError(w, response)
	}
}

func (p *Plugin) handleComplete(w http.ResponseWriter, r *http.Request, action *Action) {

	reminder := p.GetReminder(action.UserID, action.Context.ReminderID)

	for _, occurrence := range reminder.Occurrences {
		p.ClearScheduledOccurrence(reminder, occurrence)
	}

	reminder.Completed = time.Now().UTC()
	p.UpdateReminder(action.UserID, reminder)

	if post, pErr := p.API.GetPost(action.PostID); pErr != nil {
		p.API.LogError("unable to get post " + pErr.Error())
		writePostActionIntegrationResponseError(w, &model.PostActionIntegrationResponse{})
	} else {

		user, uError := p.API.GetUser(action.UserID)
		if uError != nil {
			p.API.LogError(uError.Error())
			return
		}
		finalTarget := reminder.Target
		if finalTarget == T("me") {
			finalTarget = T("you")
		} else {
			finalTarget = "@" + user.Username
		}

		messageParameters := map[string]interface{}{
			"FinalTarget": finalTarget,
			"Message":     reminder.Message,
		}

		var updateParameters = map[string]interface{}{
			"Message": reminder.Message,
		}

		post.Message = "~~" + T("reminder.message", messageParameters) + "~~\n" + T("action.complete", updateParameters)
		post.Props = model.StringInterface{}
		p.API.UpdatePost(post)
		writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})
	}

}

func (p *Plugin) handleCompleteList(w http.ResponseWriter, r *http.Request, action *Action) {

	reminder := p.GetReminder(action.UserID, action.Context.ReminderID)

	for _, occurrence := range reminder.Occurrences {
		p.ClearScheduledOccurrence(reminder, occurrence)
	}

	reminder.Completed = time.Now().UTC()
	p.UpdateReminder(action.UserID, reminder)
	p.UpdateListReminders(action.UserID, action.PostID)
	writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})
}

func (p *Plugin) handleViewCompleteList(w http.ResponseWriter, r *http.Request, action *Action) {
	p.ListCompletedReminders(action.UserID, action.PostID)
}

func (p *Plugin) handleDelete(w http.ResponseWriter, r *http.Request, action *Action) {

	reminder := p.GetReminder(action.UserID, action.Context.ReminderID)

	for _, occurrence := range reminder.Occurrences {
		p.ClearScheduledOccurrence(reminder, occurrence)
	}

	message := reminder.Message
	p.DeleteReminder(action.UserID, reminder)

	if post, pErr := p.API.GetPost(action.PostID); pErr != nil {
		p.API.LogError("unable to get post " + pErr.Error())
		writePostActionIntegrationResponseError(w, &model.PostActionIntegrationResponse{})
	} else {
		p.API.LogInfo(fmt.Sprintf("%v", post.Props))
		var deleteParameters = map[string]interface{}{
			"Message": message,
		}
		post.Message = T("action.delete", deleteParameters)
		post.Props = model.StringInterface{}
		p.API.UpdatePost(post)
		writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})
	}

}

func (p *Plugin) handleDeleteList(w http.ResponseWriter, r *http.Request, action *Action) {

	reminder := p.GetReminder(action.UserID, action.Context.ReminderID)

	for _, occurrence := range reminder.Occurrences {
		p.ClearScheduledOccurrence(reminder, occurrence)
	}

	p.DeleteReminder(action.UserID, reminder)
	p.UpdateListReminders(action.UserID, action.PostID)
	writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})
}

func (p *Plugin) handleDeleteCompleteList(w http.ResponseWriter, r *http.Request, action *Action) {

	p.DeleteCompletedReminders(action.UserID)
	p.UpdateListReminders(action.UserID, action.PostID)
	writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})
}

func (p *Plugin) handleSnooze(w http.ResponseWriter, r *http.Request, action *Action) {

	reminder := p.GetReminder(action.UserID, action.Context.ReminderID)

	for _, occurrence := range reminder.Occurrences {
		if occurrence.Id == action.Context.OccurrenceID {
			p.ClearScheduledOccurrence(reminder, occurrence)
		}
	}

	if post, pErr := p.API.GetPost(action.PostID); pErr != nil {
		p.API.LogError("unable to get post " + pErr.Error())
		writePostActionIntegrationResponseError(w, &model.PostActionIntegrationResponse{})
	} else {
		p.API.LogInfo(fmt.Sprintf("%v", post.Props))
		var snoozeParameters = map[string]interface{}{
			"Message": reminder.Message,
		}

		switch action.Context.SelectedOption {
		case "20min":
			for _, occurrence := range reminder.Occurrences {
				if occurrence.Id == action.Context.OccurrenceID {
					occurrence.Snoozed = time.Now().UTC().Round(time.Second).Add(time.Minute * time.Duration(20))
					p.UpdateReminder(action.UserID, reminder)
					p.upsertSnoozedOccurrence(&occurrence)
					post.Message = T("action.snooze.20min", snoozeParameters)
					break
				}
			}
		case "1hr":
			for _, occurrence := range reminder.Occurrences {
				if occurrence.Id == action.Context.OccurrenceID {
					occurrence.Snoozed = time.Now().UTC().Round(time.Second).Add(time.Hour * time.Duration(1))
					p.UpdateReminder(action.UserID, reminder)
					p.upsertSnoozedOccurrence(&occurrence)
					post.Message = T("action.snooze.1hr", snoozeParameters)
					break
				}
			}
		case "3hrs":
			for _, occurrence := range reminder.Occurrences {
				if occurrence.Id == action.Context.OccurrenceID {
					occurrence.Snoozed = time.Now().UTC().Round(time.Second).Add(time.Hour * time.Duration(3))
					p.UpdateReminder(action.UserID, reminder)
					p.upsertSnoozedOccurrence(&occurrence)
					post.Message = T("action.snooze.3hr", snoozeParameters)
					break
				}
			}
		case "tomorrow":
			for _, occurrence := range reminder.Occurrences {
				if occurrence.Id == action.Context.OccurrenceID {

					if user, uErr := p.API.GetUser(action.UserID); uErr != nil {
						p.API.LogError(uErr.Error())
						return
					} else {
						location := p.location(user)
						tt := time.Now().In(location).Add(time.Hour * time.Duration(24))
						occurrence.Snoozed = time.Date(tt.Year(), tt.Month(), tt.Day(), 9, 0, 0, 0, location).UTC()
						p.UpdateReminder(action.UserID, reminder)
						p.upsertSnoozedOccurrence(&occurrence)
						post.Message = T("action.snooze.tomorrow", snoozeParameters)
						break
					}
				}
			}
		case "nextweek":
			for _, occurrence := range reminder.Occurrences {
				if occurrence.Id == action.Context.OccurrenceID {

					if user, uErr := p.API.GetUser(action.UserID); uErr != nil {
						p.API.LogError(uErr.Error())
						return
					} else {
						location := p.location(user)

						todayWeekDayNum := int(time.Now().In(location).Weekday())
						weekDayNum := 1
						day := 0

						if weekDayNum < todayWeekDayNum {
							day = 7 - (todayWeekDayNum - weekDayNum)
						} else if weekDayNum >= todayWeekDayNum {
							day = 7 + (weekDayNum - todayWeekDayNum)
						}

						tt := time.Now().In(location).Add(time.Hour * time.Duration(24))
						occurrence.Snoozed = time.Date(tt.Year(), tt.Month(), tt.Day(), 9, 0, 0, 0, location).AddDate(0, 0, day).UTC()
						p.UpdateReminder(action.UserID, reminder)
						p.upsertSnoozedOccurrence(&occurrence)
						post.Message = T("action.snooze.nextweek", snoozeParameters)
						break
					}
				}
			}
		}

		post.Props = model.StringInterface{}
		p.API.UpdatePost(post)
		writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})
	}
}

func (p *Plugin) handleSnoozeList(w http.ResponseWriter, r *http.Request, action *Action) {

	reminder := p.GetReminder(action.UserID, action.Context.ReminderID)

	for _, occurrence := range reminder.Occurrences {
		if occurrence.Id == action.Context.OccurrenceID {
			p.ClearScheduledOccurrence(reminder, occurrence)
		}
	}

	if _, pErr := p.API.GetPost(action.PostID); pErr != nil {
		p.API.LogError("unable to get post " + pErr.Error())
		writePostActionIntegrationResponseError(w, &model.PostActionIntegrationResponse{})
	} else {
		switch action.Context.SelectedOption {
		case "20min":
			for _, occurrence := range reminder.Occurrences {
				if occurrence.Id == action.Context.OccurrenceID {
					occurrence.Snoozed = time.Now().UTC().Round(time.Second).Add(time.Minute * time.Duration(20))
					p.UpdateReminder(action.UserID, reminder)
					p.upsertSnoozedOccurrence(&occurrence)
					break
				}
			}
		case "1hr":
			for _, occurrence := range reminder.Occurrences {
				if occurrence.Id == action.Context.OccurrenceID {
					occurrence.Snoozed = time.Now().UTC().Round(time.Second).Add(time.Hour * time.Duration(1))
					p.UpdateReminder(action.UserID, reminder)
					p.upsertSnoozedOccurrence(&occurrence)
					break
				}
			}
		case "3hrs":
			for _, occurrence := range reminder.Occurrences {
				if occurrence.Id == action.Context.OccurrenceID {
					occurrence.Snoozed = time.Now().UTC().Round(time.Second).Add(time.Hour * time.Duration(3))
					p.UpdateReminder(action.UserID, reminder)
					p.upsertSnoozedOccurrence(&occurrence)
					break
				}
			}
		case "tomorrow":
			for _, occurrence := range reminder.Occurrences {
				if occurrence.Id == action.Context.OccurrenceID {

					if user, uErr := p.API.GetUser(action.UserID); uErr != nil {
						p.API.LogError(uErr.Error())
						return
					} else {
						location := p.location(user)
						tt := time.Now().In(location).Add(time.Hour * time.Duration(24))
						occurrence.Snoozed = time.Date(tt.Year(), tt.Month(), tt.Day(), 9, 0, 0, 0, location).UTC()
						p.UpdateReminder(action.UserID, reminder)
						p.upsertSnoozedOccurrence(&occurrence)
						break
					}
				}
			}
		case "nextweek":
			for _, occurrence := range reminder.Occurrences {
				if occurrence.Id == action.Context.OccurrenceID {

					if user, uErr := p.API.GetUser(action.UserID); uErr != nil {
						p.API.LogError(uErr.Error())
						return
					} else {
						location := p.location(user)

						todayWeekDayNum := int(time.Now().In(location).Weekday())
						weekDayNum := 1
						day := 0

						if weekDayNum < todayWeekDayNum {
							day = 7 - (todayWeekDayNum - weekDayNum)
						} else if weekDayNum >= todayWeekDayNum {
							day = 7 + (weekDayNum - todayWeekDayNum)
						}

						tt := time.Now().In(location).Add(time.Hour * time.Duration(24))
						occurrence.Snoozed = time.Date(tt.Year(), tt.Month(), tt.Day(), 9, 0, 0, 0, location).AddDate(0, 0, day).UTC()
						p.UpdateReminder(action.UserID, reminder)
						p.upsertSnoozedOccurrence(&occurrence)
						break
					}
				}
			}
		}

		p.UpdateListReminders(action.UserID, action.PostID)
		writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})
	}
}

func (p *Plugin) handleCloseList(w http.ResponseWriter, r *http.Request, action *Action) {
	p.API.DeletePost(action.PostID)
}

func writePostActionIntegrationResponseOk(w http.ResponseWriter, response *model.PostActionIntegrationResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response.ToJson())
}

func writePostActionIntegrationResponseError(w http.ResponseWriter, response *model.PostActionIntegrationResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write(response.ToJson())
}
