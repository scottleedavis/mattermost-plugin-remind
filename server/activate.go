package main

import (
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/pkg/errors"
)

const (
	botUserName    = "remindbot"
	botDisplayName = "Remindbot"
)

func (p *Plugin) OnActivate() error {
	p.ServerConfig = p.API.GetConfig()
	p.router = p.InitAPI()
	p.emptyTime = time.Time{}.AddDate(1, 1, 1)

	botUserID, err := p.ensureBotExists()
	if err != nil {
		return errors.Wrap(err, "failed to ensure bot user")
	}
	p.remindUserId = botUserID

	err = p.registerCommand()
	if err != nil {
		return errors.Wrap(err, "failed to register command")
	}

	if err := p.TranslationsPreInit(); err != nil {
		return errors.Wrap(err, "failed to initialize translations")
	}
	p.Run()

	return nil
}

func (p *Plugin) ensureBotExists() (string, error) {
	bot := &model.Bot{
		Username:    botUserName,
		DisplayName: botDisplayName,
	}
	options := []plugin.EnsureBotOption{
		plugin.ProfileImagePath("assets/icon.png"),
	}

	return p.Helpers.EnsureBot(bot, options...)
}

func (p *Plugin) OnDeactivate() error {
	p.Stop()
	return nil
}
