package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
)

func (p *Plugin) InitAPI() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/dialog", p.handleDialog).Methods("POST")

	r.HandleFunc("/view/ephemeral", p.handleViewEphemeral).Methods("POST")
	r.HandleFunc("/view/complete/list", p.handleViewCompleteList).Methods("POST")

	r.HandleFunc("/complete", p.handleComplete).Methods("POST")
	r.HandleFunc("/complete/list", p.handleCompleteList).Methods("POST")

	r.HandleFunc("/delete", p.handleDelete).Methods("POST")
	r.HandleFunc("/delete/ephemeral", p.handleDeleteEphemeral).Methods("POST")
	r.HandleFunc("/delete/list", p.handleDeleteList).Methods("POST")
	r.HandleFunc("/delete/complete/list", p.handleDeleteCompleteList).Methods("POST")

	r.HandleFunc("/snooze", p.handleSnooze).Methods("POST")
	r.HandleFunc("/snooze/list", p.handleSnoozeList).Methods("POST")

	r.HandleFunc("/close/list", p.handleCloseList).Methods("POST")

	r.HandleFunc("/next/reminders", p.handleNextReminders).Methods("POST")

	return r
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}

func (p *Plugin) handleDialog(w http.ResponseWriter, req *http.Request) {

	body, _ := io.ReadAll(req.Body)
	defer req.Body.Close()

	var request *model.SubmitDialogRequest
	_ = json.Unmarshal(body, &request)

	user, uErr := p.API.GetUser(request.UserId)
	if uErr != nil {
		p.API.LogError(uErr.Error())
		return
	}

	T, _ := p.translation(user)
	location := p.location(user)

	message := request.Submission["message"]
	target := request.Submission["target"]
	ttime := request.Submission["time"]

	if target == nil {
		target = T("me")
	}
	if target != T("me") &&
		!strings.HasPrefix(target.(string), "@") &&
		!strings.HasPrefix(target.(string), "~") {
		target = "@" + target.(string)
	}

	var when string
	if ttime.(string) == "unit.test" {
		when = "in 20 minutes"
	} else {
		when = T("in") + " " + T("button.snooze."+ttime.(string))
		switch ttime.(string) {
		case "tomorrow":
			when = T("tomorrow")
		case "nextweek":
			when = T("monday")
		}
	}

	r := &ReminderRequest{
		TeamId:   request.TeamId,
		Username: user.Username,
		Payload:  message.(string),
		Reminder: Reminder{
			Id:        model.NewId(),
			TeamId:    request.TeamId,
			Username:  user.Username,
			Message:   message.(string),
			Completed: p.emptyTime,
			Target:    target.(string),
			When:      when,
		},
	}

	if cErr := p.CreateOccurrences(r); cErr != nil {
		p.API.LogError(cErr.Error())
		return
	}

	if rErr := p.UpsertReminder(r); rErr != nil {
		p.API.LogError(rErr.Error())
		return
	}

	if r.Reminder.Target == T("me") {
		r.Reminder.Target = T("you")
	}

	useTo := strings.HasPrefix(r.Reminder.Message, T("to"))
	var useToString string
	if useTo {
		useToString = " " + T("to")
	} else {
		useToString = ""
	}

	t := ""
	if len(r.Reminder.Occurrences) > 0 {
		t = r.Reminder.Occurrences[0].Occurrence.In(location).Format(time.RFC3339)
	}
	var responseParameters = map[string]interface{}{
		"Target":  r.Reminder.Target,
		"UseTo":   useToString,
		"Message": r.Reminder.Message,
		"When": p.formatWhen(
			r.Username,
			r.Reminder.When,
			t,
			false,
		),
	}

	reminder := &model.Post{
		ChannelId: request.ChannelId,
		UserId:    p.botUserId,
		Props: model.StringInterface{
			"attachments": []*model.SlackAttachment{
				{
					Text: T("schedule.response", responseParameters),
					Actions: []*model.PostAction{
						{
							Integration: &model.PostActionIntegration{
								Context: model.StringInterface{
									"reminder_id":   r.Reminder.Id,
									"occurrence_id": r.Reminder.Occurrences[0].Id,
									"action":        "delete/ephemeral",
								},
								URL: fmt.Sprintf("/plugins/%s/delete/ephemeral", manifest.ID),
							},
							Type: model.PostActionTypeButton,
							Name: T("button.delete"),
						},
						{
							Integration: &model.PostActionIntegration{
								Context: model.StringInterface{
									"reminder_id":   r.Reminder.Id,
									"occurrence_id": r.Reminder.Occurrences[0].Id,
									"action":        "view/ephemeral",
								},
								URL: fmt.Sprintf("/plugins/%s/view/ephemeral", manifest.ID),
							},
							Type: model.PostActionTypeButton,
							Name: T("button.view.reminders"),
						},
					},
				},
			},
		},
	}
	p.API.SendEphemeralPost(user.Id, reminder)

}

func (p *Plugin) handleViewEphemeral(w http.ResponseWriter, r *http.Request) {

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var request *model.PostActionIntegrationRequest
	_ = json.Unmarshal(body, &request)

	user, uErr := p.API.GetUser(request.UserId)
	if uErr != nil {
		p.API.LogError(uErr.Error())
		writePostActionIntegrationResponseError(w, &model.PostActionIntegrationResponse{})
		return
	}
	p.API.SendEphemeralPost(user.Id, p.ListReminders(user, request.ChannelId))

	writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})

}

func (p *Plugin) handleComplete(w http.ResponseWriter, r *http.Request) {

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var request *model.PostActionIntegrationRequest
	_ = json.Unmarshal(body, &request)

	reminder := p.GetReminder(request.Context["orig_user_id"].(string), request.Context["reminder_id"].(string))
	user, uErr := p.API.GetUser(request.UserId)
	if uErr != nil {
		p.API.LogError(uErr.Error())
		writePostActionIntegrationResponseError(w, &model.PostActionIntegrationResponse{})
		return
	}
	T, _ := p.translation(user)

	for _, occurrence := range reminder.Occurrences {
		p.ClearScheduledOccurrence(reminder, occurrence)
	}

	reminder.Completed = time.Now().UTC()
	urErr := p.UpdateReminder(request.Context["orig_user_id"].(string), reminder)
	if urErr != nil {
		p.API.LogError("failed to update reminder %s", urErr)
	}

	if post, pErr := p.API.GetPost(request.PostId); pErr != nil {
		p.API.LogError("unable to get post " + pErr.Error())
		writePostActionIntegrationResponseError(w, &model.PostActionIntegrationResponse{})
	} else {

		user, uError := p.API.GetUser(request.UserId)
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
		_, upErr := p.API.UpdatePost(post)
		if upErr != nil {
			p.API.LogError("failed to update post %s", upErr)
		}

		if reminder.Username != user.Username {
			if originalUser, uErr := p.API.GetUserByUsername(reminder.Username); uErr != nil {
				p.API.LogError(uErr.Error())
				writePostActionIntegrationResponseError(w, &model.PostActionIntegrationResponse{})
				return
			} else {
				if channel, cErr := p.API.GetDirectChannel(p.botUserId, originalUser.Id); cErr != nil {
					p.API.LogError("failed to create channel " + cErr.Error())
					writePostActionIntegrationResponseError(w, &model.PostActionIntegrationResponse{})
				} else {
					var postbackUpdateParameters = map[string]interface{}{
						"User":    "@" + user.Username,
						"Message": reminder.Message,
					}
					if _, pErr := p.API.CreatePost(&model.Post{
						ChannelId: channel.Id,
						UserId:    p.botUserId,
						Message:   T("action.complete.callback", postbackUpdateParameters),
					}); pErr != nil {
						p.API.LogError(pErr.Error())
						writePostActionIntegrationResponseError(w, &model.PostActionIntegrationResponse{})
					}
				}
			}
		}

		writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})
	}

}

func (p *Plugin) handleDelete(w http.ResponseWriter, r *http.Request) {

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var request *model.PostActionIntegrationRequest
	_ = json.Unmarshal(body, &request)

	reminder := p.GetReminder(request.Context["orig_user_id"].(string), request.Context["reminder_id"].(string))
	user, uErr := p.API.GetUser(request.UserId)
	if uErr != nil {
		p.API.LogError(uErr.Error())
		writePostActionIntegrationResponseError(w, &model.PostActionIntegrationResponse{})
		return
	}
	T, _ := p.translation(user)

	for _, occurrence := range reminder.Occurrences {
		p.ClearScheduledOccurrence(reminder, occurrence)
	}

	message := reminder.Message
	dErr := p.DeleteReminder(request.Context["orig_user_id"].(string), reminder)
	if dErr != nil {
		p.API.LogError("failed to delete reminder %s", dErr)
	}

	if post, pErr := p.API.GetPost(request.PostId); pErr != nil {
		p.API.LogError(pErr.Error())
		writePostActionIntegrationResponseError(w, &model.PostActionIntegrationResponse{})
	} else {
		var deleteParameters = map[string]interface{}{
			"Message": message,
		}
		post.Message = T("action.delete", deleteParameters)
		post.Props = model.StringInterface{}
		_, upErr := p.API.UpdatePost(post)
		if upErr != nil {
			p.API.LogError("failed to update post %s", upErr)
		}
		writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})
	}

}

func (p *Plugin) handleDeleteEphemeral(w http.ResponseWriter, r *http.Request) {

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var request *model.PostActionIntegrationRequest
	_ = json.Unmarshal(body, &request)

	reminder := p.GetReminder(request.UserId, request.Context["reminder_id"].(string))
	user, uErr := p.API.GetUser(request.UserId)
	if uErr != nil {
		p.API.LogError(uErr.Error())
		writePostActionIntegrationResponseError(w, &model.PostActionIntegrationResponse{})
		return
	}
	T, _ := p.translation(user)

	for _, occurrence := range reminder.Occurrences {
		p.ClearScheduledOccurrence(reminder, occurrence)
	}

	message := reminder.Message
	dErr := p.DeleteReminder(request.UserId, reminder)
	if dErr != nil {
		p.API.LogError("failed to delete reminder %s", dErr)
	}
	var deleteParameters = map[string]interface{}{
		"Message": message,
	}
	post := &model.Post{
		Id:        request.PostId,
		UserId:    p.botUserId,
		ChannelId: request.ChannelId,
		Message:   T("action.delete", deleteParameters),
	}
	p.API.UpdateEphemeralPost(request.UserId, post)
	writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})

}

func (p *Plugin) handleSnooze(w http.ResponseWriter, r *http.Request) {

	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var request *model.PostActionIntegrationRequest
	_ = json.Unmarshal(body, &request)

	reminder := p.GetReminder(request.Context["orig_user_id"].(string), request.Context["reminder_id"].(string))
	user, uErr := p.API.GetUser(request.UserId)
	if uErr != nil {
		p.API.LogError(uErr.Error())
		writePostActionIntegrationResponseError(w, &model.PostActionIntegrationResponse{})
		return
	}
	T, _ := p.translation(user)

	for _, occurrence := range reminder.Occurrences {
		if occurrence.Id == request.Context["occurrence_id"].(string) {
			p.ClearScheduledOccurrence(reminder, occurrence)
		}
	}

	if post, pErr := p.API.GetPost(request.PostId); pErr != nil {
		p.API.LogError("unable to get post " + pErr.Error())
		writePostActionIntegrationResponseError(w, &model.PostActionIntegrationResponse{})
	} else {
		var snoozeParameters = map[string]interface{}{
			"Message": reminder.Message,
		}

		switch request.Context["selected_option"].(string) {
		case "20min":
			for i, occurrence := range reminder.Occurrences {
				if occurrence.Id == request.Context["occurrence_id"].(string) {
					occurrence.Snoozed = time.Now().UTC().Round(time.Second).Add(time.Minute * time.Duration(20))
					reminder.Occurrences[i] = occurrence
					upErr := p.UpdateReminder(request.Context["orig_user_id"].(string), reminder)
					if upErr != nil {
						p.API.LogError("failed to update reminder %s", upErr)
					}
					p.upsertSnoozedOccurrence(&occurrence)
					post.Message = T("action.snooze.20min", snoozeParameters)
					break
				}
			}
		case "1hr":
			for i, occurrence := range reminder.Occurrences {
				if occurrence.Id == request.Context["occurrence_id"].(string) {
					occurrence.Snoozed = time.Now().UTC().Round(time.Second).Add(time.Hour * time.Duration(1))
					reminder.Occurrences[i] = occurrence
					upErr := p.UpdateReminder(request.Context["orig_user_id"].(string), reminder)
					if upErr != nil {
						p.API.LogError("failed to update reminder %s", upErr)
					}
					p.upsertSnoozedOccurrence(&occurrence)
					post.Message = T("action.snooze.1hr", snoozeParameters)
					break
				}
			}
		case "3hrs":
			for i, occurrence := range reminder.Occurrences {
				if occurrence.Id == request.Context["occurrence_id"].(string) {
					occurrence.Snoozed = time.Now().UTC().Round(time.Second).Add(time.Hour * time.Duration(3))
					reminder.Occurrences[i] = occurrence
					upErr := p.UpdateReminder(request.Context["orig_user_id"].(string), reminder)
					if upErr != nil {
						p.API.LogError("failed to update reminder %s", upErr)
					}
					p.upsertSnoozedOccurrence(&occurrence)
					post.Message = T("action.snooze.3hr", snoozeParameters)
					break
				}
			}
		case "tomorrow":
			for i, occurrence := range reminder.Occurrences {
				if occurrence.Id == request.Context["occurrence_id"].(string) {

					if user, uErr := p.API.GetUser(request.UserId); uErr != nil {
						p.API.LogError(uErr.Error())
						return
					} else {
						location := p.location(user)
						tt := time.Now().In(location).Add(time.Hour * time.Duration(24))
						occurrence.Snoozed = time.Date(tt.Year(), tt.Month(), tt.Day(), 9, 0, 0, 0, location).UTC()
						reminder.Occurrences[i] = occurrence
						upErr := p.UpdateReminder(request.Context["orig_user_id"].(string), reminder)
						if upErr != nil {
							p.API.LogError("failed to update reminder %s", upErr)
						}
						p.upsertSnoozedOccurrence(&occurrence)
						post.Message = T("action.snooze.tomorrow", snoozeParameters)
						break
					}
				}
			}
		case "nextweek":
			for i, occurrence := range reminder.Occurrences {
				if occurrence.Id == request.Context["occurrence_id"].(string) {

					if user, uErr := p.API.GetUser(request.UserId); uErr != nil {
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
						reminder.Occurrences[i] = occurrence
						upErr := p.UpdateReminder(request.Context["orig_user_id"].(string), reminder)
						if upErr != nil {
							p.API.LogError("failed to update reminder %s", upErr)
						}
						p.upsertSnoozedOccurrence(&occurrence)
						post.Message = T("action.snooze.nextweek", snoozeParameters)
						break
					}
				}
			}
		}

		post.Props = model.StringInterface{}
		_, upErr := p.API.UpdatePost(post)
		if upErr != nil {
			p.API.LogError("failed to update post %s", upErr)
		}
		writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})
	}
}

func (p *Plugin) handleNextReminders(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var request *model.PostActionIntegrationRequest
	_ = json.Unmarshal(body, &request)
	p.UpdateListReminders(request.UserId, request.PostId, request.ChannelId, int(request.Context["offset"].(float64)))
	writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})
}

func (p *Plugin) handleCompleteList(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var request *model.PostActionIntegrationRequest
	_ = json.Unmarshal(body, &request)
	reminder := p.GetReminder(request.UserId, request.Context["reminder_id"].(string))

	for _, occurrence := range reminder.Occurrences {
		p.ClearScheduledOccurrence(reminder, occurrence)
	}

	reminder.Completed = time.Now().UTC()
	upErr := p.UpdateReminder(request.UserId, reminder)
	if upErr != nil {
		p.API.LogError("failed to update reminder %s", upErr)
	}
	p.UpdateListReminders(request.UserId, request.PostId, request.ChannelId, int(request.Context["offset"].(float64)))
	writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})
}

func (p *Plugin) handleViewCompleteList(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var request *model.PostActionIntegrationRequest
	_ = json.Unmarshal(body, &request)
	p.ListCompletedReminders(request.UserId, request.PostId, request.ChannelId)
	writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})
}

func (p *Plugin) handleDeleteList(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var request *model.PostActionIntegrationRequest
	_ = json.Unmarshal(body, &request)
	reminder := p.GetReminder(request.UserId, request.Context["reminder_id"].(string))

	for _, occurrence := range reminder.Occurrences {
		p.ClearScheduledOccurrence(reminder, occurrence)
	}

	dErr := p.DeleteReminder(request.UserId, reminder)
	if dErr != nil {
		p.API.LogError("failed to update post %s", dErr)
	}
	p.UpdateListReminders(request.UserId, request.PostId, request.ChannelId, int(request.Context["offset"].(float64)))
	writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})
}

func (p *Plugin) handleDeleteCompleteList(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var request *model.PostActionIntegrationRequest
	_ = json.Unmarshal(body, &request)
	p.DeleteCompletedReminders(request.UserId)
	p.UpdateListReminders(request.UserId, request.PostId, request.ChannelId, int(request.Context["offset"].(float64)))
	writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})
}

func (p *Plugin) handleSnoozeList(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var request *model.PostActionIntegrationRequest
	_ = json.Unmarshal(body, &request)
	reminder := p.GetReminder(request.UserId, request.Context["reminder_id"].(string))

	for _, occurrence := range reminder.Occurrences {
		if occurrence.Id == request.Context["occurrence_id"].(string) {
			p.ClearScheduledOccurrence(reminder, occurrence)
		}
	}

	switch request.Context["selected_option"].(string) {
	case "20min":
		for i, occurrence := range reminder.Occurrences {
			if occurrence.Id == request.Context["occurrence_id"].(string) {
				occurrence.Snoozed = time.Now().UTC().Round(time.Second).Add(time.Minute * time.Duration(20))
				reminder.Occurrences[i] = occurrence
				upErr := p.UpdateReminder(request.UserId, reminder)
				if upErr != nil {
					p.API.LogError("failed to update reminder %s", upErr)
				}
				p.upsertSnoozedOccurrence(&occurrence)
				break
			}
		}
	case "1hr":
		for i, occurrence := range reminder.Occurrences {
			if occurrence.Id == request.Context["occurrence_id"].(string) {
				occurrence.Snoozed = time.Now().UTC().Round(time.Second).Add(time.Hour * time.Duration(1))
				reminder.Occurrences[i] = occurrence
				upErr := p.UpdateReminder(request.UserId, reminder)
				if upErr != nil {
					p.API.LogError("failed to update reminder %s", upErr)
				}
				p.upsertSnoozedOccurrence(&occurrence)
				break
			}
		}
	case "3hrs":
		for i, occurrence := range reminder.Occurrences {
			if occurrence.Id == request.Context["occurrence_id"].(string) {
				occurrence.Snoozed = time.Now().UTC().Round(time.Second).Add(time.Hour * time.Duration(3))
				reminder.Occurrences[i] = occurrence
				upErr := p.UpdateReminder(request.UserId, reminder)
				if upErr != nil {
					p.API.LogError("failed to update reminder %s", upErr)
				}
				p.upsertSnoozedOccurrence(&occurrence)
				break
			}
		}
	case "tomorrow":
		for i, occurrence := range reminder.Occurrences {
			if occurrence.Id == request.Context["occurrence_id"].(string) {

				if user, uErr := p.API.GetUser(request.UserId); uErr != nil {
					p.API.LogError(uErr.Error())
					return
				} else {
					location := p.location(user)
					tt := time.Now().In(location).Add(time.Hour * time.Duration(24))
					occurrence.Snoozed = time.Date(tt.Year(), tt.Month(), tt.Day(), 9, 0, 0, 0, location).UTC()
					reminder.Occurrences[i] = occurrence
					upErr := p.UpdateReminder(request.UserId, reminder)
					if upErr != nil {
						p.API.LogError("failed to update reminder %s", upErr)
					}
					p.upsertSnoozedOccurrence(&occurrence)
					break
				}
			}
		}
	case "nextweek":
		for i, occurrence := range reminder.Occurrences {
			if occurrence.Id == request.Context["occurrence_id"].(string) {

				if user, uErr := p.API.GetUser(request.UserId); uErr != nil {
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
					reminder.Occurrences[i] = occurrence
					upErr := p.UpdateReminder(request.UserId, reminder)
					if upErr != nil {
						p.API.LogError("failed to update reminder %s", upErr)
					}
					p.upsertSnoozedOccurrence(&occurrence)
					break
				}
			}
		}
	}

	p.UpdateListReminders(request.UserId, request.PostId, request.ChannelId, int(request.Context["offset"].(float64)))
	writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})
}

func (p *Plugin) handleCloseList(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	var request *model.PostActionIntegrationRequest
	_ = json.Unmarshal(body, &request)
	p.API.DeleteEphemeralPost(request.UserId, request.PostId)
	writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})
}

func writePostActionIntegrationResponseOk(w http.ResponseWriter, response *model.PostActionIntegrationResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	responseJSON, _ := json.Marshal(response)
	_, _ = w.Write(responseJSON)
}

func writePostActionIntegrationResponseError(w http.ResponseWriter, response *model.PostActionIntegrationResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	responseJSON, _ := json.Marshal(response)
	_, _ = w.Write(responseJSON)
}
