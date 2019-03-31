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
		AutoCompleteDesc: "Tests Translation",
		DisplayName:      "Remind Plugin Command",
		Description:      "Set a reminder",
	}); err != nil {
		return errors.Wrap(err, "failed to register command")
	}

	return nil
}

func (p *Plugin) unregisterCommand(teamId string) error {

	p.API.UnregisterCommand(teamId, CommandTrigger)
	return nil
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {

	user, _ := p.API.GetUser(args.UserId)
	T, _ := p.translation(user)

	// original //
	// return &model.CommandResponse{
	// 	ResponseType: model.COMMAND_RESPONSE_TYPE_IN_CHANNEL,
	// 	Text:         fmt.Sprintf(T("exception")),
	// }, nil

	if strings.HasSuffix(args.Command, T("help")) {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:         fmt.Sprintf(T("help-response")),
		}, nil
	}

	// if strings.HasSuffix(args.Command, T("list")) {
	// 	return &model.CommandResponse{
	// 		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
	// 		Text:         fmt.Sprintf(a.ListReminders(user.Id, args.ChannelId)),
	// 	}
	// }

	// if strings.HasSuffix(args.Command, T("clear")) {
	// 	return &model.CommandResponse{
	// 		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
	// 		Text:         fmt.Sprintf(a.DeleteReminders(user.Id)),
	// 	}
	// }

	// payload := strings.Trim(strings.Replace(args.Command, "/"+model.CMD_REMIND, "", -1), " ")

	// if strings.HasPrefix(payload, T("app.reminder.me")) ||
	// 	strings.HasPrefix(payload, "@") ||
	// 	strings.HasPrefix(payload, "~") {

	// 	request := model.ReminderRequest{
	// 		TeamId:      args.TeamId,
	// 		UserId:      args.UserId,
	// 		Payload:     payload,
	// 		Reminder:    model.Reminder{},
	// 		Occurrences: model.Occurrences{},
	// 	}
	// 	response, err := a.ScheduleReminder(&request)

	// 	if err != nil {
	// 		return &model.CommandResponse{
	// 			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
	// 			Text:         fmt.Sprintf(T(model.REMIND_EXCEPTION_TEXT)),
	// 		}
	// 	}

	// 	return &model.CommandResponse{
	// 		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
	// 		Text:         fmt.Sprintf("%s", response),
	// 	}
	// }

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text:         fmt.Sprintf(T("exception-response")),
	}, nil

}
