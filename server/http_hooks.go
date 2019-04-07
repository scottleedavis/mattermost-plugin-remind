package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

// ActionContext passed from action buttons
type ActionContext struct {
	ReminderID   string `json:"reminder_id"`
	OccurrenceId string `json:"occurrence_id"`
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

	p.API.LogInfo("UserId: " + fmt.Sprintf("%v", action.UserID))
	p.API.LogInfo("PostID: " + fmt.Sprintf("%v", action.PostID))
	p.API.LogInfo("Context: " + fmt.Sprintf("%v", action.Context))

	// complete reminder

	//strike through original reminder

	// response := &model.PostActionIntegrationResponse{}
	// response.EphemeralText = "TODO: ~~strikethrough original reminder~~"
	// writePostActionIntegrationResponse(w, response)
	// writePostActionIntegrationResponseOk(w, r)

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
