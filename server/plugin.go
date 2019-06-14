package main

import (
	"github.com/gorilla/mux"
	"io/ioutil"
	"time"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

type Plugin struct {
	plugin.MattermostPlugin

	router *mux.Router

	ServerConfig *model.Config

	URL string

	remindUserId string

	activated bool

	running bool

	emptyTime time.Time

	defaultTime time.Time

	supportedLocales []string

	readFile func(path string) ([]byte, error)
}

func NewPlugin() *Plugin {
	return &Plugin{
		readFile: ioutil.ReadFile,
	}
}
