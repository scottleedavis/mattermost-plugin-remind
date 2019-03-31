package main

import (
	"fmt"
	"github.com/pkg/errors"

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

	u, err := p.API.GetUser(args.UserId)

	if err != nil {
		p.API.LogError(err.Error())
	}
	T := GetUserTranslations(u.Locale)

	p.API.LogError(T("exception"))

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_IN_CHANNEL,
		Text:         fmt.Sprintf(T("exception")),
	}, nil
}
