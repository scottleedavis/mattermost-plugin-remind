package main

import (
	"fmt"
	"github.com/blang/semver"
	"github.com/mattermost/mattermost-server/model"
	"github.com/pkg/errors"
	"path/filepath"
	"time"
)

const minimumServerVersion = "5.10.0"
const botName = "remindbot"
const botDisplayName = "Remindbot"

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

	teams, err := p.API.GetTeams()
	if err != nil {
		return errors.Wrap(err, "failed to query teams OnActivate")
	}

	p.router = p.InitAPI()

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

func (p *Plugin) OnConfigurationChange() error {
	p.ServerConfig = p.API.GetConfig()
	p.URL = "http://127.0.0.1" + fmt.Sprintf("%s", *p.ServerConfig.ServiceSettings.ListenAddress)
	return nil
}

func (p *Plugin) ensureBotExists() (string, *model.AppError) {
	p.API.LogDebug("Ensuring Remindbot exists")

	bot, createErr := p.API.CreateBot(&model.Bot{
		Username:    botName,
		DisplayName: botDisplayName,
		Description: "Sets and triggers reminders",
	})
	if createErr != nil {
		p.API.LogDebug("Failed to create "+botDisplayName+". Attempting to find existing one.", "err", createErr)

		// Unable to create the bot, so it should already exist
		user, err := p.API.GetUserByUsername(botName)
		if err != nil || user == nil {
			p.API.LogError("Failed to find "+botDisplayName+" user", "err", err)
			return "", err
		}

		bot, err = p.API.GetBot(user.Id, true)
		if err != nil {
			p.API.LogError("Failed to find "+botDisplayName, "err", err)
			return "", err
		}

		p.API.LogDebug("Found " + botDisplayName)
	} else {
		if err := p.setBotProfileImage(bot.UserId); err != nil {
			p.API.LogWarn("Failed to set profile image for bot", "err", err)
		}

		p.API.LogDebug(botDisplayName + " created")
	}

	p.remindUserId = bot.UserId

	return bot.UserId, nil
}

func (p *Plugin) setBotProfileImage(botUserId string) *model.AppError {
	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		return &model.AppError{Message: err.Error()}
	}

	profileImage, err := p.readFile(filepath.Join(bundlePath, "assets", "icon.png"))
	if err != nil {
		return &model.AppError{Message: err.Error()}
	}

	return p.API.SetProfileImage(botUserId, profileImage)
}
