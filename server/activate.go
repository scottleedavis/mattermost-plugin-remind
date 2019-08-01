package main

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/mattermost/mattermost-server/model"
	"github.com/pkg/errors"
)

const botName = "remindbot"
const botDisplayName = "Remindbot"

func (p *Plugin) OnActivate() error {
	p.ServerConfig = p.API.GetConfig()
	if p.ServerConfig.ServiceSettings.SiteURL == nil {
		return errors.New("siteURL is not set. Please set a siteURL and restart the plugin")
	}

	teams, err := p.API.GetTeams()
	if err != nil {
		return errors.Wrap(err, "failed to query teams OnActivate")
	}

	p.router = p.InitAPI()
	p.ensureBotExists()
	p.emptyTime = time.Time{}.AddDate(1, 1, 1)

	for _, team := range teams {
		if err := p.registerCommand(team.Id); err != nil {
			return errors.Wrap(err, "failed to register command")
		}
	}

	p.URL = fmt.Sprintf("%s", *p.ServerConfig.ServiceSettings.SiteURL)
	if err := p.TranslationsPreInit(); err != nil {
		return errors.Wrap(err, "failed to initialize translations")
	}
	p.Run()

	return nil
}

func (p *Plugin) OnDeactivate() error {
	p.Stop()
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
