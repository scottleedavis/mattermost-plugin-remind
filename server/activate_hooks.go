package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/model"
)

// OnActivate is invoked when the plugin is activated.
//
// This demo implementation logs a message to the demo channel whenever the plugin is activated.
func (p *Plugin) OnActivate() error {
	teams, err := p.API.GetTeams()
	if err != nil {
		p.API.LogError(
			"failed to query teams OnActivate",
			"error", err.Error(),
		)
	}

	for _, team := range teams {
		if _, err := p.API.CreatePost(&model.Post{
			UserId:    p.demoUserId,
			ChannelId: p.demoChannelIds[team.Id],
			Message: fmt.Sprintf(
				"OnActivate: %s", PluginId,
			),
			Type: "custom_demo_plugin",
			Props: map[string]interface{}{
				"username":     p.Username,
				"channel_name": p.ChannelName,
			},
		}); err != nil {
			p.API.LogError(
				"failed to post OnActivate message",
				"error", err.Error(),
			)
		}

		if err := p.registerCommand(team.Id); err != nil {
			p.API.LogError(
				"failed to register command",
				"error", err.Error(),
			)
		}
	}

	return nil
}

// OnDeactivate is invoked when the plugin is deactivated. This is the plugin's last chance to use
// the API, and the plugin will be terminated shortly after this invocation.
//
// This demo implementation logs a message to the demo channel whenever the plugin is deactivated.
func (p *Plugin) OnDeactivate() error {
	teams, err := p.API.GetTeams()
	if err != nil {
		p.API.LogError(
			"failed to query teams OnDeactivate",
			"error", err.Error(),
		)
	}

	for _, team := range teams {
		if _, err := p.API.CreatePost(&model.Post{
			UserId:    p.demoUserId,
			ChannelId: p.demoChannelIds[team.Id],
			Message: fmt.Sprintf(
				"OnDeactivate: %s", PluginId,
			),
			Type: "custom_demo_plugin",
			Props: map[string]interface{}{
				"username":     p.Username,
				"channel_name": p.ChannelName,
			},
		}); err != nil {
			p.API.LogError(
				"failed to post OnDeactivate message",
				"error", err.Error(),
			)
		}

		if err := p.registerCommand(team.Id); err != nil {
			p.API.LogError(
				"failed to register command",
				"error", err.Error(),
			)
		}
	}

	return nil
}
