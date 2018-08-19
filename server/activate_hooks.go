package main

import (
	"fmt"

	"github.com/mattermost/mattermost-server/model"
	"github.com/google/uuid"
)

func (p *Plugin) OnActivate() error {
	teams, err := p.API.GetTeams()
	if err != nil {
		p.API.LogError(
			"failed to query teams OnActivate",
			"error", err.Error(),
		)
	}

	p.API.LogError("OnActivate")

	p.createUser()

	for _, team := range teams {

		if err := p.registerCommand(team.Id); err != nil {
			p.API.LogError(
				"failed to register command",
				"error", err.Error(),
			)
		}
	}

	p.Run()

	return nil
}

func (p *Plugin) OnDeactivate() error {
	teams, err := p.API.GetTeams()
	if err != nil {
		p.API.LogError(
			"failed to query teams OnDeactivate",
			"error", err.Error(),
		)
	}

	p.API.LogError("OnDeactivate")

	// p.deleteUser()

	for _, team := range teams {

		if err := p.unregisterCommand(team.Id); err != nil {
			p.API.LogError(
				"failed to register command",
				"error", err.Error(),
			)
		}
	}

	p.Stop()

	return nil
}

func (p *Plugin) createUser() (*model.User, error) {

	p.API.LogError("create user")

	user, err := p.API.GetUserByUsername("remind")
	if err != nil {
		p.API.LogError("Failed to get user remind")

		guid, gerr := uuid.NewRandom()
		if gerr != nil {
			p.API.LogError("Failed to generate guid")
			return nil, gerr
		}

		user := model.User{Email: "scottleedavis@gmail.com", Nickname: "Remind", Password: guid.String(), Username: "remind", Roles: model.SYSTEM_USER_ROLE_ID}

		cuser, err := p.API.CreateUser(&user)
		if err != nil {
			p.API.LogError("Failed to create remind user " + fmt.Sprintf("%v", err))
			return cuser, err
		}

		p.remindUserId = cuser.Id

		return cuser, nil
	}

	p.API.LogDebug("id " + user.Id)

	p.remindUserId = user.Id

	return user, nil
}

// func (p *Plugin) deleteUser() {

//     p.API.LogError("delete user")
// 	user, err := p.API.GetUserByUsername("remind")
// 	if err != nil {
// 		return
// 	}
// 	derr := p.API.DeleteUser(user.Id)
// 	if derr != nil {
// 		p.API.LogError("Failed to delete remind user "+fmt.Sprintf("%v", derr))
// 	}
// }
