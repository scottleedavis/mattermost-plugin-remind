package main

import (
	"io/ioutil"
	"time"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

type Plugin struct {
	plugin.MattermostPlugin

	ServerConfig *model.Config

	remindUserId string

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
