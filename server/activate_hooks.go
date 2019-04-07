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

	// if dErr := p.deactivateBotUser(); dErr != nil {
	// return dErr
	// }

	for _, team := range teams {
		if cErr := p.API.UnregisterCommand(team.Id, CommandTrigger); cErr != nil {
			return errors.Wrap(cErr, "failed to unregister command")
		}
	}

	return nil
}

// func (p *Plugin) activateBotUser() (*model.Bot, error) {
func (p *Plugin) activateBotUser() (*model.User, error) {

	// bot, err := p.API.GetBot(CommandTrigger, true)
	bot, err := p.API.GetUserByUsername(CommandTrigger)
	if err != nil {
		p.API.LogError(fmt.Sprintf("failed to get user %s: %v", CommandTrigger, err))

		user := model.User{
			Email:    "-@-.-",
			Nickname: "Remind",
			Password: model.NewId(),
			Username: CommandTrigger,
			Roles:    model.SYSTEM_USER_ROLE_ID,
		}

		cuser, cerr := p.API.CreateUser(&user)
		if cerr != nil {
			p.API.LogError("failed to create user: " + fmt.Sprintf("%v", cerr))
			return nil, cerr
		}

		// b := model.Bot{
		// 	UserId:      cuser.Id,
		// 	Username:    CommandTrigger + "_bot",
		// 	OwnerId:     manifest.Id,
		// 	DisplayName: "Remind",
		// 	Description: "Sets and triggers reminders",
		// }

		// newBot, bErr := p.API.CreateBot(&b)
		// if bErr != nil {
		// 	p.API.LogError(fmt.Sprintf("failed to create %s user: %v", CommandTrigger, bErr))
		// 	return nil, bErr
		// }

		// p.remindUserId = newBot.UserId
		p.remindUserId = cuser.Id

		// return newBot, nil
		return cuser, nil
	}

	// p.remindUserId = bot.UserId
	p.remindUserId = bot.Id

	return bot, nil

}

func (p *Plugin) deactivateBotUser() error {

	// botUser, err := p.API.GetBot(p.remindUserId, true)
	botUser, err := p.API.GetUser(p.remindUserId)
	if err != nil {
		return err
	}
	// dErr := p.API.PermanentDeleteBot(botUser.UserId)
	dErr := p.API.DeleteUser(botUser.Id)
	if dErr != nil {
		p.API.LogError("Failed to delete " + CommandTrigger + " bot " + fmt.Sprintf("%v", dErr))
		return dErr
	}
	return nil
}
