package main

import (
	// "encoding/json"
	// "fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

// InitAPI initializes the REST API
func (p *Plugin) InitAPI() *mux.Router {
	r := mux.NewRouter()

	apiV1 := r.PathPrefix("/api/v1").Subrouter()

	apiV1.HandleFunc("/complete", p.handleComplete).Methods("POST")

	return r
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}

func (p *Plugin) handleComplete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	p.API.LogInfo(vars["id"])
	p.API.LogInfo(vars["name"])
	p.API.LogInfo(vars["type"])
	p.API.LogInfo(vars["datasource"])
	p.API.LogInfo(vars["integration"])

	response := &model.PostActionIntegrationResponse{}
	response.EphemeralText = "bobs your uncle"
	writePostActionIntegrationResponse(w, response)
}

func writePostActionIntegrationResponse(w http.ResponseWriter, response *model.PostActionIntegrationResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(response.ToJson())
}
