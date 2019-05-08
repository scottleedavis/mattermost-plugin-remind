package main

import (
	"testing"
	"encoding/json"
	"time"

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
			Id: model.NewId(),
			ReminderId: "ididididid",
			Occurrence:  time.Now(),
		},
	}

	reminders := []Reminder{
		{
			Id: model.NewId(),
			Username: user.Username,
			Message: "Message",
			When: "in 1 second",
			Occurrences: occurrences,
			Completed:  time.Time{}.AddDate(1, 1, 1),
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
		api.On("GetDirectChannel", mock.Anything, mock.Anything).Return(channel, nil)
		api.On("CreatePost", mock.AnythingOfType("*model.Post")).Return(nil, nil)
		api.On("KVGet", mock.Anything).Return(stringReminders, nil)
		return api
	}

	t.Run("if list happens in remind channel", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		output := p.ListReminders(user, channel.Id)

		assert.Equal(t, output, "")

	})

	t.Run("if list happens in other channel", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		output := p.ListReminders(user, model.NewId())

		assert.Equal(t, output, "list.reminders")

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
			Id: model.NewId(),
			ReminderId: "ididididid",
			Occurrence:  time.Now(),
		},
	}

	reminders := []Reminder{
		{
			Id: model.NewId(),
			Username: user.Username,
			Message: "Message",
			When: "in 1 second",
			Occurrences: occurrences,
			Completed:  time.Time{}.AddDate(1, 1, 1),
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
		Props: model.StringInterface{},
	}

	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
		api.On("GetUser", mock.AnythingOfType("string")).Return(user, nil)
		api.On("UpdatePost", mock.AnythingOfType("*model.Post")).Return(nil, nil)
		api.On("GetPost", mock.AnythingOfType("string")).Return(post, nil)
		api.On("KVGet", mock.Anything).Return(stringReminders, nil)
		return api
	}

	t.Run("if update list happens", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		p.UpdateListReminders(user.Id, model.NewId(), 0)

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
			Id: model.NewId(),
			ReminderId: "ididididid",
			Occurrence:  time.Now(),
		},
	}

	reminders := []Reminder{
		{
			Id: model.NewId(),
			Username: user.Username,
			Message: "Message",
			When: "in 1 second",
			Occurrences: occurrences,
			Completed:  time.Time{}.AddDate(1, 1, 1),
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
		Props: model.StringInterface{},
	}

	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
		api.On("GetUser", mock.AnythingOfType("string")).Return(user, nil)
		api.On("UpdatePost", mock.AnythingOfType("*model.Post")).Return(nil, nil)
		api.On("GetPost", mock.AnythingOfType("string")).Return(post, nil)
		api.On("KVGet", mock.Anything).Return(stringReminders, nil)
		return api
	}

	t.Run("if list completed happens", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		p.ListCompletedReminders(user.Id, post.Id)

	})

}


//func TestDeleteCompletedReminders(t *testing.T) {
//
//	user := &model.User{
//		Email:    "-@-.-",
//		Nickname: "TestUser",
//		Password: model.NewId(),
//		Username: "testuser",
//		Roles:    model.SYSTEM_USER_ROLE_ID,
//		Locale:   "en",
//	}
//
//	occurrences := []Occurrence{
//		{
//			Id: model.NewId(),
//			ReminderId: "ididididid",
//			Occurrence:  time.Now(),
//		},
//	}
//
//	reminders := []Reminder{
//		{
//			Id: model.NewId(),
//			Username: user.Username,
//			Message: "Message",
//			When: "in 1 second",
//			Occurrences: occurrences,
//			Completed:  time.Time{}.AddDate(1, 1, 1),
//		},
//	}
//
//	stringReminders, _ := json.Marshal(reminders)
//
//	channel := &model.Channel{
//		Id: model.NewId(),
//	}
//
//	post := &model.Post{
//		ChannelId:     channel.Id,
//		PendingPostId: model.NewId(),
//		UserId:        user.Id,
//		Props: model.StringInterface{},
//	}
//
//	setupAPI := func() *plugintest.API {
//		api := &plugintest.API{}
//		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
//		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
//		api.On("LogInfo", mock.Anything).Maybe()
//		api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
//		api.On("GetUser", mock.AnythingOfType("string")).Return(user, nil)
//		api.On("UpdatePost", mock.AnythingOfType("*model.Post")).Return(nil, nil)
//		api.On("GetPost", mock.AnythingOfType("string")).Return(post, nil)
//		api.On("KVGet", mock.Anything).Return(stringReminders, nil)
//		api.On("KVSet", user.Username, stringReminders ).Return(nil)
//		return api
//	}
//
//	t.Run("if delete completed happens", func(t *testing.T) {
//
//		api := setupAPI()
//		defer api.AssertExpectations(t)
//
//		p := &Plugin{}
//		p.API = api
//
//		p.DeleteCompletedReminders(user.Id)
//
//	})
//
//}
