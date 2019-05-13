package main

import (
	"testing"
	"fmt"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleCommand(t *testing.T) {

	trigger := "remind"

	user := &model.User{
		Id:       model.NewId(),
		Username: model.NewRandomString(10),
	}

	t.Run("/remind", func(t *testing.T) {

		setupAPI := func() *plugintest.API {
			api := &plugintest.API{}
			api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
			api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
			api.On("LogInfo", mock.Anything).Maybe()
			api.On("GetUser", mock.Anything).Return(user, nil)
			api.On("OpenInteractiveDialog", mock.Anything).Return(nil)

			return api
		}

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.URL = fmt.Sprintf("http://localhost/plugins/%s", manifest.Id)
		p.router = p.InitAPI()
		p.API = api

		r, err := p.ExecuteCommand(nil, &model.CommandArgs{
			Command: fmt.Sprintf("/%s", trigger),
			UserId:  "userID1",
		})

		assert.NotNil(t,r)
		assert.Nil(t, err)
	})

	t.Run("/remind help", func(t *testing.T) {

		setupAPI := func() *plugintest.API {
			api := &plugintest.API{}
			api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
			api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
			api.On("LogInfo", mock.Anything).Maybe()
			api.On("GetUser", mock.Anything).Return(user, nil)
			api.On("SendEphemeralPost", mock.Anything, mock.Anything).Return(nil)
			return api
		}

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.URL = fmt.Sprintf("http://localhost/plugins/%s", manifest.Id)
		p.router = p.InitAPI()
		p.API = api

		r, err := p.ExecuteCommand(nil, &model.CommandArgs{
			Command: fmt.Sprintf("/%s help", trigger),
			UserId:  "userID1",
		})

		assert.NotNil(t,r)
		assert.Nil(t, err)
	})

	t.Run("/remind list", func(t *testing.T) {

		//channel := &model.Channel{
		//	Id: 	model.NewId(),
		//	Name:   model.NewRandomString(10),
		//}
		//post := &model.Post{
		//	Id:        model.NewId(),
		//	ChannelId: channel.Id,
		//}
		//
		//testTime := time.Now().UTC().Round(time.Second)
		//
		//occurrences := []Occurrence{
		//	{
		//		Id:         model.NewId(),
		//		ReminderId: model.NewId(),
		//		Occurrence: testTime,
		//	},
		//}
		//reminders := []Reminder{
		//	{
		//		Id:          model.NewId(),
		//		TeamId:      model.NewId(),
		//		Username:    user.Username,
		//		Message:     "Hello",
		//		Target:      "me",
		//		When:        "in one minute",
		//		Occurrences: occurrences,
		//	},
		//}
		//stringReminders, _ := json.Marshal(reminders)
		//
		//setupAPI := func() *plugintest.API {
		//	api := &plugintest.API{}
		//	api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		//	api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		//	api.On("LogInfo", mock.Anything).Maybe()
		//	api.On("GetUser", mock.Anything).Return(user, nil)
		//	api.On("GetUserByUsername", mock.Anything).Return(user, nil)
		//	api.On("KVGet", user.Username).Return(stringReminders, nil)
		//	api.On("GetDirectChannel", mock.Anything, mock.Anything).Return(channel, nil)
		//	api.On("CreatePost", post).Return(post, nil)
		//
		//	return api
		//}
		//
		//api := setupAPI()
		//defer api.AssertExpectations(t)
		//
		//p := &Plugin{}
		//p.URL = fmt.Sprintf("http://localhost/plugins/%s", manifest.Id)
		//p.router = p.InitAPI()
		//p.API = api
		//
		//r, err := p.ExecuteCommand(nil, &model.CommandArgs{
		//	Command: fmt.Sprintf("/%s list", trigger),
		//	UserId:  "userID1",
		//})
		//
		//assert.NotNil(t,r)
		//assert.Nil(t, err)
	})


	t.Run("/remind me foo in 2 seconds", func(t *testing.T) {

		//channel := &model.Channel{
		//	Id: 	model.NewId(),
		//	Name:   model.NewRandomString(10),
		//}
		//
		//testTime := time.Now().UTC().Round(time.Second)
		//
		//occurrences := []Occurrence{
		//	{
		//		Id:         model.NewId(),
		//		ReminderId: model.NewId(),
		//		Occurrence: testTime,
		//	},
		//}
		//reminders := []Reminder{
		//	{
		//		Id:          model.NewId(),
		//		TeamId:      model.NewId(),
		//		Username:    user.Username,
		//		Message:     "Hello",
		//		Target:      "me",
		//		When:        "in 2 seconds",
		//		Occurrences: occurrences,
		//	},
		//}
		//stringReminders, _ := json.Marshal(reminders)
		//
		//setupAPI := func() *plugintest.API {
		//	api := &plugintest.API{}
		//	api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		//	api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		//	api.On("LogInfo", mock.Anything).Maybe()
		//	api.On("GetUser", mock.Anything).Return(user, nil)
		//	api.On("GetUserByUsername", mock.Anything).Return(user, nil)
		//	api.On("KVGet", user.Username).Return(stringReminders, nil)
		//	api.On("GetDirectChannel", mock.Anything, mock.Anything).Return(channel, nil)
		//	api.On("SendEphemeralPost", mock.Anything, mock.Anything).Return(nil)
		//
		//	return api
		//}
		//
		//api := setupAPI()
		//defer api.AssertExpectations(t)
		//
		//p := &Plugin{}
		//p.URL = fmt.Sprintf("http://localhost/plugins/%s", manifest.Id)
		//p.router = p.InitAPI()
		//p.API = api
		//
		//r, err := p.ExecuteCommand(nil, &model.CommandArgs{
		//	Command: fmt.Sprintf("/%s me Hello in 2 seconds", trigger),
		//	UserId:  "userID1",
		//})
		//
		//assert.NotNil(t,r)
		//assert.Nil(t, err)
	})

	t.Run("/remind __clear", func(t *testing.T) {

		//setupAPI := func() *plugintest.API {
		//	api := &plugintest.API{}
		//	api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		//	api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		//	api.On("LogInfo", mock.Anything).Maybe()
		//	api.On("GetUser", mock.Anything).Return(user, nil)
		//	api.On("SendEphemeralPost", mock.Anything, mock.Anything).Return(nil)
		//	return api
		//}
		//
		//api := setupAPI()
		//defer api.AssertExpectations(t)
		//
		//p := &Plugin{}
		//p.URL = fmt.Sprintf("http://localhost/plugins/%s", manifest.Id)
		//p.router = p.InitAPI()
		//p.API = api
		//
		//r, err := p.ExecuteCommand(nil, &model.CommandArgs{
		//	Command: fmt.Sprintf("/%s __clear", trigger),
		//	UserId:  "userID1",
		//})
		//
		//assert.NotNil(t,r)
		//assert.Nil(t, err)
	})

	t.Run("/remind __version", func(t *testing.T) {

		setupAPI := func() *plugintest.API {
			api := &plugintest.API{}
			api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
			api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
			api.On("LogInfo", mock.Anything).Maybe()
			api.On("GetUser", mock.Anything).Return(user, nil)
			api.On("SendEphemeralPost", mock.Anything, mock.Anything).Return(nil)
			return api
		}

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.URL = fmt.Sprintf("http://localhost/plugins/%s", manifest.Id)
		p.router = p.InitAPI()
		p.API = api

		r, err := p.ExecuteCommand(nil, &model.CommandArgs{
			Command: fmt.Sprintf("/%s __version", trigger),
			UserId:  "userID1",
		})

		assert.NotNil(t,r)
		assert.Nil(t, err)
	})

	t.Run("/remind __user", func(t *testing.T) {

		setupAPI := func() *plugintest.API {
			api := &plugintest.API{}
			api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
			api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
			api.On("LogInfo", mock.Anything).Maybe()
			api.On("GetUser", mock.Anything).Return(user, nil)
			api.On("SendEphemeralPost", mock.Anything, mock.Anything).Return(nil)
			return api
		}

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.URL = fmt.Sprintf("http://localhost/plugins/%s", manifest.Id)
		p.router = p.InitAPI()
		p.API = api

		r, err := p.ExecuteCommand(nil, &model.CommandArgs{
			Command: fmt.Sprintf("/%s __user", trigger),
			UserId:  "userID1",
		})

		assert.NotNil(t,r)
		assert.Nil(t, err)
	})
}