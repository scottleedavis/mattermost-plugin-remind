package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

const CommandTrigger = "remind"

func (p *Plugin) registerCommand(teamId string) error {
	if err := p.API.RegisterCommand(&model.Command{
		TeamId:           teamId,
		Trigger:          CommandTrigger,
		AutoComplete:     true,
		AutoCompleteHint: "[@someone or ~channel] [what] [when]",
		AutoCompleteDesc: "Enables or disables the demo plugin hooks.",
		DisplayName:      "Demo Plugin Command",
		Description:      "Set a reminder",
	}); err != nil {
		p.API.LogError(
			"failed to register command",
			"error", err.Error(),
		)
	}

	return nil
}

func (p *Plugin) emitStatusChange() {
	p.API.PublishWebSocketEvent("status_change", map[string]interface{}{
		"enabled": true,
	}, &model.WebsocketBroadcast{})
}

// ExecuteCommand executes a command that has been previously registered via the RegisterCommand
// API.
//
// This demo implementation responds to a /demo_plugin command, allowing the user to enable
// or disable the demo plugin's hooks functionality (but leave the command and webapp enabled).
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	if !strings.HasPrefix(args.Command, "/"+CommandTrigger) {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("Unknown command: " + args.Command),
		}, nil
	}

	if strings.HasSuffix(args.Command, "me foo in 2 seconds") {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text: "Ok!  Will remind you foo in 2 seconds",
		}, nil
	}

	
	// if strings.HasSuffix(args.Command, "true") {
	// 	if !p.disabled {
	// 		return &model.CommandResponse{
	// 			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
	// 			Text:         "The demo plugin hooks are already enabled.",
	// 		}, nil
	// 	}

	// 	p.disabled = false
	// 	p.emitStatusChange()

	// 	return &model.CommandResponse{
	// 		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
	// 		Text:         "Enabled demo plugin hooks.",
	// 	}, nil

	// } else if strings.HasSuffix(args.Command, "false") {
	// 	if p.disabled {
	// 		return &model.CommandResponse{
	// 			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
	// 			Text:         "The demo plugin hooks are already disabled.",
	// 		}, nil
	// 	}

	// 	p.disabled = true
		p.emitStatusChange()

	// 	return &model.CommandResponse{
	// 		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
	// 		Text:         "Disabled demo plugin hooks.",
	// 	}, nil
	// }

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf("Unknown command action: " + args.Command),
	}, nil
}
