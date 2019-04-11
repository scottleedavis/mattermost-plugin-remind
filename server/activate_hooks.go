package main

import (
	"fmt"
	"time"

	"github.com/blang/semver"
	"github.com/mattermost/mattermost-server/model"
	"github.com/pkg/errors"
)

const minimumServerVersion = "5.10.0"

func (p *Plugin) checkServerVersion() error {
	serverVersion, err := semver.Parse(p.API.GetServerVersion())
	if err != nil {
		return errors.Wrap(err, "failed to parse server version")
	}

	r := semver.MustParseRange(">=" + minimumServerVersion)
	if !r(serverVersion) {
		return fmt.Errorf("this plugin requires Mattermost v%s or later", minimumServerVersion)
	}

	return nil
}

func (p *Plugin) OnActivate() error {
	if err := p.checkServerVersion(); err != nil {
		return err
	}

	// configuration := p.getConfiguration()

	teams, err := p.API.GetTeams()
	if err != nil {
		return errors.Wrap(err, "failed to query teams OnActivate")
	}

	p.ensureBotExists()

	for _, team := range teams {
		if err := p.registerCommand(team.Id); err != nil {
			return errors.Wrap(err, "failed to register command")
		}
	}

	if err := TranslationsPreInit(); err != nil {
		return errors.Wrap(err, "failed to initialize translations OnActivate message")
	}

	p.emptyTime = time.Time{}.AddDate(1, 1, 1)
	p.supportedLocales = []string{"en"}
	p.ServerConfig = p.API.GetConfig()

	p.Run()

	return nil
}

func (p *Plugin) OnDeactivate() error {

	teams, err := p.API.GetTeams()
	if err != nil {
		return errors.Wrap(err, "failed to query teams OnDeactivate")
	}

	p.Stop()

	for _, team := range teams {
		if cErr := p.API.UnregisterCommand(team.Id, CommandTrigger); cErr != nil {
			return errors.Wrap(cErr, "failed to unregister command")
		}
	}

	return nil
}

func (p *Plugin) ensureBotExists() (string, *model.AppError) {
	p.API.LogInfo("Ensuring Remindbot exists")

	bot, createErr := p.API.CreateBot(&model.Bot{
		Username:    "remindbot",
		DisplayName: "Remindbot",
		Description: "Sets and triggers reminders",
	})
	if createErr != nil {
		p.API.LogDebug("Failed to create Remindbot. Attempting to find existing one.", "err", createErr)

		// Unable to create the bot, so it should already exist
		user, err := p.API.GetUserByUsername("remindbot")
		if err != nil || user == nil {
			p.API.LogError("Failed to find Remind user", "err", err)
			return "", err
		}

		bot, err = p.API.GetBot(user.Id, true)
		if err != nil {
			p.API.LogError("Failed to find Remindbot", "err", err)
			return "", err
		}

		p.API.LogDebug("Found Remindbot")
	} else {
		p.API.LogInfo("Remindbot created")
	}

	p.remindUserId = bot.UserId

	return bot.UserId, nil
}
