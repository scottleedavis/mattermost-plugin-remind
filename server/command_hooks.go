package main

import (
	"fmt"
	"github.com/pkg/errors"
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
		AutoCompleteDesc: "Set a reminder",
		DisplayName:      "Remind Plugin Command",
		Description:      "A command used to set a reminder",
	}); err != nil {
		return errors.Wrap(err, "failed to register command")
	}

	return nil
}

func (p *Plugin) unregisterCommand(teamId string) error {
	return p.API.UnregisterCommand(teamId, CommandTrigger)
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {

	user, _ := p.API.GetUser(args.UserId)
	T, _ := p.translation(user)

	if strings.HasSuffix(args.Command, T("help")) {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf(T("help.response")),
		}, nil
	}

	if strings.HasSuffix(args.Command, T("list")) {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf(p.ListReminders(user, args.ChannelId)),
		}, nil
	}

	if strings.HasSuffix(args.Command, T("clear")) {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf(p.DeleteReminders(user)),
		}, nil
	}

	payload := strings.Trim(strings.Replace(args.Command, "/"+CommandTrigger, "", -1), " ")

	if strings.HasPrefix(payload, T("me")) ||
		strings.HasPrefix(payload, "@") ||
		strings.HasPrefix(payload, "~") {

		request := ReminderRequest{
			TeamId:   args.TeamId,
			Username: user.Username,
			Payload:  payload,
			Reminder: Reminder{},
		}
		response, err := p.ScheduleReminder(&request)

		if err != nil {
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         fmt.Sprintf(T("exception.response")),
			}, nil
		}

		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf("%s", response),
		}, nil
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf(T("exception.response")),
	}, nil

}
