package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

// ActionContext passed from action buttons
type ActionContext struct {
	ReminderID   string `json:"reminder_id"`
	OccurrenceID string `json:"occurrence_id"`
	Action       string `json:"action"`
}

// Action type for decoding action buttons
type Action struct {
	UserID  string         `json:"user_id"`
	PostID  string         `json:"post_id"`
	Context *ActionContext `json:"context"`
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {

	var action *Action
	json.NewDecoder(r.Body).Decode(&action)

	switch action.Context.Action {
	case "complete":
		p.handleComplete(w, r, action)
	case "delete":
		p.handleDelete(w, r, action)
	case "snooze":
		p.handleSnooze(w, r, action)
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
		var updateParameters = map[string]interface{}{
			"Message": reminder.Message,
		}
		post.Message = "~~" + post.Message + "~~\n" + T("action.complete", updateParameters)
		post.Props = model.StringInterface{}
		p.API.UpdatePost(post)
		writePostActionIntegrationResponseOk(w, &model.PostActionIntegrationResponse{})
	}

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

func (p *Plugin) handleSnooze(w http.ResponseWriter, r *http.Request, action *Action) {

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
