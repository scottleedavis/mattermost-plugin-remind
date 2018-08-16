package main

import (
	"github.com/mattermost/mattermost-server/plugin"
)

type Plugin struct {
	plugin.MattermostPlugin

	// The user to use as part of the demo plugin, created automatically if it does not exist.
	Username string

	// The channel to use as part of the demo plugin, created for each team automatically if it does not exist.
	ChannelName string

	// disabled tracks whether or not the plugin has been disabled after activation. It always starts enabled.
	disabled bool

	// demoUserId is the id of the user specified above.
	demoUserId string

	// demoChannelIds maps team ids to the channels created for each using the channel name above.
	demoChannelIds map[string]string
}
