package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

const CommandTrigger = "remind"

func runSchedule() {
   for i := 3; i < 1000; i+=3 {
        fmt.Printf("i loop: %v", i)
    }
}

func (p *Plugin) registerCommand(teamId string) error {
	if err := p.API.RegisterCommand(&model.Command{
		TeamId:           teamId,
		Trigger:          CommandTrigger,
		AutoComplete:     true,
		AutoCompleteHint: "[@someone or ~channel] [what] [when]",
		AutoCompleteDesc: "Enables or disables the demo plugin hooks.",
		DisplayName:      "Remind Plugin Command",
		Description:      "Set a reminder",
	}); err != nil {
		p.API.LogError(
			"failed to register command",
			"error", err.Error(),
		)
	}

	// runSchedule()

	return nil
}

func (p *Plugin) emitStatusChange() {
	// p.API.PublishWebSocketEvent("status_change", map[string]interface{}{
	// 	"enabled": true,
	// }, &model.WebsocketBroadcast{})
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {

	user, err := p.API.GetUser(args.UserId)
	
	if err != nil {
		p.API.LogError("failed to query user %s", args.UserId)
	}

	if strings.HasSuffix(args.Command, "help") {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text: fmt.Sprintf("Ok! %s %s", user.Username, user.Timezone["automaticTimezone"]),
		}, nil
	}

	if strings.HasSuffix(args.Command, "list") {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text: fmt.Sprintf("Ok! %s %s", user.Username, user.Timezone["automaticTimezone"]),
		}, nil
	}

	if strings.HasSuffix(args.Command, "version") {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text: fmt.Sprintf("Ok! %s %s", user.Username, user.Timezone["automaticTimezone"]),
		}, nil
	}

	payload := strings.Trim(strings.Replace(args.Command, "/"+CommandTrigger, "", -1),"")
	commandSplit := strings.Split(payload," ")

	if len(commandSplit) == 0 {	
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("No Command Found"),
		}, nil
	}

	if commandSplit[1] == "me" {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text: fmt.Sprintf("Ok! %s %s", user.Username, user.Timezone["automaticTimezone"]),
		}, nil
	}

	if strings.HasPrefix(commandSplit[1][:1], "@"){
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text: fmt.Sprintf("Ok! %s %s", user.Username, user.Timezone["automaticTimezone"]),
		}, nil
	}

	if strings.HasPrefix(commandSplit[1][:1], "~") {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text: fmt.Sprintf("Ok! %s %s", user.Username, user.Timezone["automaticTimezone"]),
		}, nil
	}

	p.emitStatusChange()

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf("Unknown command action: " + args.Command),
	}, nil
}

