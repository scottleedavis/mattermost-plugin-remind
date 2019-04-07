package main

import (
	"encoding/json"
	// "fmt"
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

	// p.API.LogInfo("UserId: " + fmt.Sprintf("%v", action.UserID))
	// p.API.LogInfo("PostID: " + fmt.Sprintf("%v", action.PostID))
	// p.API.LogInfo("Context: " + fmt.Sprintf("%v", action.Context))

	/*
		if result := <-a.Srv.Store.Remind().GetReminder(reminderId); result.Err != nil {
			return result.Err
		} else {
			reminder := result.Data.(model.Reminder)
			reminder.Completed = time.Now().Format(time.RFC3339)
			if result := <-a.Srv.Store.Remind().SaveReminder(&reminder); result.Err != nil {
				return result.Err
			}
			if result := <-a.Srv.Store.Remind().DeleteForReminder(reminderId); result.Err != nil {
				return result.Err
			}
			var updateParameters = map[string]interface{}{
				"Message": reminder.Message,
			}
			update.Message = "~~" + post.Message + "~~\n" + T("app.reminder.update.complete", updateParameters)
		}
	*/

	//get reminder
	reminder := p.GetReminder(action.UserID, action.Context.ReminderID)
	//remove occurrences
	for _, occurrence := range reminder.Occurrences {
		p.ClearScheduledOccurrence(reminder, occurrence)
	}
	// complete reminder
	reminder.Completed = time.Now().UTC()
	p.UpdateReminder(action.UserID, reminder)
	//strike through original reminder
	if post, pErr := p.API.GetPost(action.PostID); pErr != nil {
		p.API.LogError("unable to get post " + pErr.Error())
		response := &model.PostActionIntegrationResponse{}
		writePostActionIntegrationResponseError(w, response)
	} else {
		var updateParameters = map[string]interface{}{
			"Message": reminder.Message,
		}
		post.Message = "~~" + post.Message + "~~\n" + T("app.reminder.update.complete", updateParameters)
		p.API.UpdatePost(post)
		response := &model.PostActionIntegrationResponse{}
		response.EphemeralText = post.Message
		writePostActionIntegrationResponseOk(w, response)

	}

}

func (p *Plugin) handleDelete(w http.ResponseWriter, r *http.Request, action *Action) {

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
