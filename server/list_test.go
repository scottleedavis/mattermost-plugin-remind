package main

import (
	"testing"
	"time"

	"encoding/json"
	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListReminders(t *testing.T) {

	user := &model.User{
		Email:    "-@-.-",
		Nickname: "TestUser",
		Password: model.NewId(),
		Username: "testuser",
		Roles:    model.SYSTEM_USER_ROLE_ID,
		Locale:   "en",
	}

	occurrences := []Occurrence{
		{
			Id:         model.NewId(),
			ReminderId: model.NewId(),
			Occurrence: time.Now(),
		},
	}

	reminders := []Reminder{
		{
			Id:          model.NewId(),
			Username:    user.Username,
			Message:     "Message",
			When:        "in 1 second",
			Occurrences: occurrences,
			Completed:   time.Time{}.AddDate(1, 1, 1),
		},
	}

	stringReminders, _ := json.Marshal(reminders)

	channel := &model.Channel{
		Id: model.NewId(),
	}

	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
		api.On("KVGet", mock.Anything).Return(stringReminders, nil)
		return api
	}

	t.Run("if list happens", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		assert.NotNil(t, p.ListReminders(user, channel.Id))

	})

}

func TestUpdateListReminders(t *testing.T) {

	user := &model.User{
		Email:    "-@-.-",
		Nickname: "TestUser",
		Password: model.NewId(),
		Username: "testuser",
		Roles:    model.SYSTEM_USER_ROLE_ID,
		Locale:   "en",
	}

	occurrences := []Occurrence{
		{
			Id:         model.NewId(),
			ReminderId: "ididididid",
			Occurrence: time.Now(),
		},
	}

	reminders := []Reminder{
		{
			Id:          model.NewId(),
			Username:    user.Username,
			Message:     "Message",
			When:        "in 1 second",
			Occurrences: occurrences,
			Completed:   time.Time{}.AddDate(1, 1, 1),
		},
	}

	stringReminders, _ := json.Marshal(reminders)

	channel := &model.Channel{
		Id: model.NewId(),
	}

	post := &model.Post{
		ChannelId:     channel.Id,
		PendingPostId: model.NewId(),
		UserId:        user.Id,
		Props:         model.StringInterface{},
	}

	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
		api.On("GetUser", mock.AnythingOfType("string")).Return(user, nil)
		api.On("UpdateEphemeralPost", mock.Anything, mock.Anything).Return(post)
		api.On("KVGet", mock.Anything).Return(stringReminders, nil)
		return api
	}

	t.Run("if update list happens", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		p.UpdateListReminders(user.Id, model.NewId(), channel.Id, 0)

	})

}

func TestListCompletedReminders(t *testing.T) {

	user := &model.User{
		Email:    "-@-.-",
		Nickname: "TestUser",
		Password: model.NewId(),
		Username: "testuser",
		Roles:    model.SYSTEM_USER_ROLE_ID,
		Locale:   "en",
	}

	occurrences := []Occurrence{
		{
			Id:         model.NewId(),
			ReminderId: "ididididid",
			Occurrence: time.Now(),
		},
	}

	reminders := []Reminder{
		{
			Id:          model.NewId(),
			Username:    user.Username,
			Message:     "Message",
			When:        "in 1 second",
			Occurrences: occurrences,
			Completed:   time.Time{}.AddDate(1, 1, 1),
		},
	}

	stringReminders, _ := json.Marshal(reminders)

	channel := &model.Channel{
		Id: model.NewId(),
	}

	post := &model.Post{
		ChannelId:     channel.Id,
		PendingPostId: model.NewId(),
		UserId:        user.Id,
		Props:         model.StringInterface{},
	}

	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
		api.On("GetUser", mock.AnythingOfType("string")).Return(user, nil)
		api.On("UpdateEphemeralPost", mock.AnythingOfType("string"), mock.Anything).Return(post)
		api.On("KVGet", mock.Anything).Return(stringReminders, nil)
		return api
	}

	t.Run("if list completed happens", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		p.ListCompletedReminders(user.Id, post.Id, channel.Id)

	})

}

func TestDeleteCompletedReminders(t *testing.T) {

	user := &model.User{
		Email:    "-@-.-",
		Nickname: "TestUser",
		Password: model.NewId(),
		Username: "testuser",
		Roles:    model.SYSTEM_USER_ROLE_ID,
		Locale:   "en",
	}

	occurrences := []Occurrence{
		{
			Id:         model.NewId(),
			ReminderId: model.NewId(),
			Occurrence: time.Now(),
		},
	}

	reminders := []Reminder{
		{
			Id:          model.NewId(),
			Username:    user.Username,
			Message:     "Message",
			When:        "in 1 second",
			Occurrences: occurrences,
			Completed:   time.Time{}.AddDate(1, 1, 1),
		},
	}

	stringReminders, _ := json.Marshal(reminders)
	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
		api.On("GetUser", mock.AnythingOfType("string")).Return(user, nil)
		api.On("KVGet", mock.Anything).Return(stringReminders, nil)
		api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
		return api
	}

	t.Run("if delete completed happens", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		p.DeleteCompletedReminders(user.Id)

	})

}
