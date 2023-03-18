package main

import (
	"io/ioutil"
	"time"

	"github.com/gorilla/mux"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
)

type Plugin struct {
	plugin.MattermostPlugin
	router      *mux.Router
	botUserId   string
	running     bool
	emptyTime   time.Time
	defaultTime time.Time

	ServerConfig *model.Config

	readFile func(path string) ([]byte, error)
	locales  map[string]string
}

func NewPlugin() *Plugin {
	return &Plugin{
		readFile: ioutil.ReadFile,
		locales:  make(map[string]string),
	}
}
