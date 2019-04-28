package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/pkg/errors"
)

const CommandTrigger = "remind"

func (p *Plugin) registerCommand(teamId string) error {
	if err := p.API.RegisterCommand(&model.Command{
		TeamId:           teamId,
		Trigger:          CommandTrigger,
		Username:         botName,
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

	user, uErr := p.API.GetUser(args.UserId)
	if uErr != nil {
		return &model.CommandResponse{}, uErr
	}

	T, locale := p.translation(user)
	location := p.location(user)

	if strings.Trim(args.Command, " ") == "/"+CommandTrigger {
		p.InteractiveSchedule(args.TriggerId, user)
		return &model.CommandResponse{}, nil
	}

	if strings.HasSuffix(args.Command, T("help")) {
		post := model.Post{
			ChannelId: args.ChannelId,
			UserId:    p.remindUserId,
			Message:   T("help.response"),
		}
		p.API.SendEphemeralPost(user.Id, &post)
		return &model.CommandResponse{}, nil
	}

	if strings.HasSuffix(args.Command, T("list")) {
		listMessage := p.ListReminders(user, args.ChannelId)
		if listMessage != "" {
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text:         fmt.Sprintf(listMessage),
				Username:     botName,
			}, nil
		}
		// p.API.SendEphemeralPost(user.Id, p.ListReminders(user, args.ChannelId))
		return &model.CommandResponse{}, nil
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
		reminder, err := p.ScheduleReminder(&request, args.ChannelId)

		if err != nil {
			post := model.Post{
				ChannelId: args.ChannelId,
				UserId:    p.remindUserId,
				Message:   T("exception.response"),
			}
			p.API.SendEphemeralPost(user.Id, &post)
			return &model.CommandResponse{}, nil
		}

		p.API.SendEphemeralPost(user.Id, reminder)
		return &model.CommandResponse{}, nil

	}

	// debug & troubleshooting commands

	// clear all reminders for current user
	if strings.HasSuffix(args.Command, "__clear") {
		post := model.Post{
			ChannelId: args.ChannelId,
			UserId:    p.remindUserId,
			Message:   p.DeleteReminders(user),
		}
		p.API.SendEphemeralPost(user.Id, &post)
		return &model.CommandResponse{}, nil
	}

	// display the plugin version
	if strings.HasSuffix(args.Command, "__version") {
		post := model.Post{
			ChannelId: args.ChannelId,
			UserId:    p.remindUserId,
			Message:   manifest.Version,
		}
		p.API.SendEphemeralPost(user.Id, &post)
		return &model.CommandResponse{}, nil
	}

	// display the locale & location of user
	if strings.HasSuffix(args.Command, "__user") {
		post := model.Post{
			ChannelId: args.ChannelId,
			UserId:    p.remindUserId,
			Message:   "locale: " + locale + "\nlocation: " + location.String(),
		}
		p.API.SendEphemeralPost(user.Id, &post)
		return &model.CommandResponse{}, nil
	}

	post := model.Post{
		ChannelId: args.ChannelId,
		UserId:    p.remindUserId,
		Message:   T("exception.response"),
	}
	p.API.SendEphemeralPost(user.Id, &post)
	return &model.CommandResponse{}, nil

}
