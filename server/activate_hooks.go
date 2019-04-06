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

	// configuration := p.getConfiguration()10

	teams, err := p.API.GetTeams()
	if err != nil {
		return errors.Wrap(err, "failed to query teams OnActivate")
	}

	p.activateBotUser()

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

	p.router = p.InitAPI()

	p.Run()

	return nil
}

func (p *Plugin) OnDeactivate() error {

	teams, err := p.API.GetTeams()
	if err != nil {
		return errors.Wrap(err, "failed to query teams OnDeactivate")
	}

	p.Stop()
	p.deactivateBotUser()

	for _, team := range teams {
		if err := p.API.UnregisterCommand(team.Id, CommandTrigger); err != nil {
			return errors.Wrap(err, "failed to unregister command")
		}
	}

	return nil
}

func (p *Plugin) activateBotUser() (*model.Bot, error) {

	bot, err := p.API.GetBot("remind", true)
	if err != nil {
		p.API.LogError("failed to get user remind")

		bot := model.Bot{
			UserId:      model.NewId(),
			Username:    "remind",
			DisplayName: "Remind",
			Description: "Sets and triggers reminders",
		}

		newBot, bErr := p.API.CreateBot(&bot)
		if err != nil {
			p.API.LogError("failed to create remind user " + fmt.Sprintf("%v", err))
			return newBot, bErr
		}

		p.remindUserId = newBot.UserId

		return newBot, nil
	}

	p.remindUserId = bot.UserId

	return bot, nil

}

func (p *Plugin) deactivateBotUser() {

	botUser, err := p.API.GetBot(p.remindUserId, true)
	if err != nil {
		return
	}
	derr := p.API.PermanentDeleteBot(botUser.UserId)
	if derr != nil {
		p.API.LogError("Failed to delete remind bot " + fmt.Sprintf("%v", derr))
	}
}
