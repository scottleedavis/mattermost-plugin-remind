package main

import (
	"github.com/mattermost/mattermost-server/plugin"
)

type Plugin struct {
	plugin.MattermostPlugin

	Username string

	ChannelName string

	remindUserId string

	running bool
}
