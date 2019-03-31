package main

import (
	"fmt"
	"github.com/mattermost/mattermost-server/mlog"

	"github.com/blang/semver"
	"github.com/google/uuid"
	"github.com/mattermost/mattermost-server/model"
	"github.com/pkg/errors"
)

const minimumServerVersion = "5.4.0"

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

	p.createBotUser()

	for _, team := range teams {
		if err := p.registerCommand(team.Id); err != nil {
			return errors.Wrap(err, "failed to register command")
		}
	}

	if err := TranslationsPreInit(); err != nil {
		mlog.Error(err.Error())
	}

	return nil
}

func (p *Plugin) OnDeactivate() error {

	teams, err := p.API.GetTeams()
	if err != nil {
		return errors.Wrap(err, "failed to query teams OnDeactivate")
	}

	p.deleteBotUser()

	for _, team := range teams {
		if err := p.API.UnregisterCommand(team.Id, CommandTrigger); err != nil {
			return errors.Wrap(err, "failed to unregister command")
		}
	}

	return nil
}

func (p *Plugin) createBotUser() (*model.User, error) {

	user, err := p.API.GetUserByUsername("remind")
	if err != nil {
		p.API.LogError("failed to get user remind")

		guid, gerr := uuid.NewRandom()
		if gerr != nil {
			p.API.LogError("failed to generate guid")
			return nil, gerr
		}

		user := model.User{Email: "-@-.-", Nickname: "Remind", Password: guid.String(), Username: "remind", Roles: model.SYSTEM_USER_ROLE_ID}

		cuser, err := p.API.CreateUser(&user)
		if err != nil {
			p.API.LogError("failed to create remind user " + fmt.Sprintf("%v", err))
			return cuser, err
		}

		p.remindUserId = cuser.Id

		return cuser, nil
	}

	p.API.LogDebug("user.Id " + user.Id)

	p.remindUserId = user.Id

	return user, nil

}

func (p *Plugin) deleteBotUser() {

	user, err := p.API.GetUserByUsername("remind")
	if err != nil {
		return
	}
	derr := p.API.DeleteUser(user.Id)
	if derr != nil {
		p.API.LogError("Failed to delete remind user " + fmt.Sprintf("%v", derr))
	}
}
