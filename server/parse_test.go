package main

import (
	"testing"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestParseRequest(t *testing.T) {

	user := &model.User{
		Email:    "-@-.-",
		Nickname: "TestUser",
		Password: model.NewId(),
		Username: "testuser",
		Roles:    model.SYSTEM_USER_ROLE_ID,
		Locale:   "en",
	}

	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
		return api
	}

	t.Run("if no quotes", func(t *testing.T) {
		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		request := &ReminderRequest{
			TeamId:   model.NewId(),
			Username: user.Username,
			Payload:  "Hello in one minute",
			Reminder: Reminder{
				Id:        model.NewId(),
				TeamId:    model.NewId(),
				Username:  user.Username,
				Message:   "Hello",
				Completed: p.emptyTime,
				Target:    "me",
				When:      "in one minute",
			},
		}

		assert.NotNil(t, p.ParseRequest(request))
	})

	t.Run("if with quotes", func(t *testing.T) {
		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		request := &ReminderRequest{
			TeamId:   model.NewId(),
			Username: user.Username,
			Payload:  "\"Hello\" in one minute",
			Reminder: Reminder{
				Id:        model.NewId(),
				TeamId:    model.NewId(),
				Username:  user.Username,
				Message:   "Hello",
				Completed: p.emptyTime,
				Target:    "me",
				When:      "in one minute",
			},
		}

		assert.NotNil(t, p.ParseRequest(request))
	})
}

func TestFindWhen(t *testing.T) {

	user := &model.User{
		Email:    "-@-.-",
		Nickname: "TestUser",
		Password: model.NewId(),
		Username: "testuser",
		Roles:    model.SYSTEM_USER_ROLE_ID,
		Locale:   "en",
	}

	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
		return api
	}

	t.Run("if findWhen", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		request := &ReminderRequest{
			TeamId:   model.NewId(),
			Username: user.Username,
			Payload:  "Hello in one minute",
			Reminder: Reminder{
				Id:        model.NewId(),
				TeamId:    model.NewId(),
				Username:  user.Username,
				Message:   "Hello",
				Completed: p.emptyTime,
				Target:    "me",
				When:      "in one minute",
			},
		}

		err := p.findWhen(request)
		assert.True(t, err == nil)

		request = &ReminderRequest{
			TeamId:   model.NewId(),
			Username: user.Username,
			Payload:  "Hello every tuesday at 10am",
			Reminder: Reminder{
				Id:        model.NewId(),
				TeamId:    model.NewId(),
				Username:  user.Username,
				Message:   "Hello",
				Completed: p.emptyTime,
				Target:    "me",
				When:      "every tuesday at 10am",
			},
		}

		err = p.findWhen(request)
		assert.True(t, err == nil)

		request = &ReminderRequest{
			TeamId:   model.NewId(),
			Username: user.Username,
			Payload:  "Hello today at noon",
			Reminder: Reminder{
				Id:        model.NewId(),
				TeamId:    model.NewId(),
				Username:  user.Username,
				Message:   "Hello",
				Completed: p.emptyTime,
				Target:    "me",
				When:      "today at noon",
			},
		}

		err = p.findWhen(request)
		assert.True(t, err == nil)

		request = &ReminderRequest{
			TeamId:   model.NewId(),
			Username: user.Username,
			Payload:  "Hello tomorrow at noon",
			Reminder: Reminder{
				Id:        model.NewId(),
				TeamId:    model.NewId(),
				Username:  user.Username,
				Message:   "Hello",
				Completed: p.emptyTime,
				Target:    "me",
				When:      "tomorrow at noon",
			},
		}

		err = p.findWhen(request)
		assert.True(t, err == nil)

		request = &ReminderRequest{
			TeamId:   model.NewId(),
			Username: user.Username,
			Payload:  "Hello monday at 11:11am",
			Reminder: Reminder{
				Id:        model.NewId(),
				TeamId:    model.NewId(),
				Username:  user.Username,
				Message:   "Hello",
				Completed: p.emptyTime,
				Target:    "me",
				When:      "monday at 11:11am",
			},
		}

		err = p.findWhen(request)
		assert.True(t, err == nil)

		request = &ReminderRequest{
			TeamId:   model.NewId(),
			Username: user.Username,
			Payload:  "Hello monday",
			Reminder: Reminder{
				Id:        model.NewId(),
				TeamId:    model.NewId(),
				Username:  user.Username,
				Message:   "Hello",
				Completed: p.emptyTime,
				Target:    "me",
				When:      "monday",
			},
		}

		err = p.findWhen(request)
		assert.True(t, err == nil)

		request = &ReminderRequest{
			TeamId:   model.NewId(),
			Username: user.Username,
			Payload:  "Hello at 2:04 pm",
			Reminder: Reminder{
				Id:        model.NewId(),
				TeamId:    model.NewId(),
				Username:  user.Username,
				Message:   "Hello",
				Completed: p.emptyTime,
				Target:    "me",
				When:      "at 2:04 pm",
			},
		}

		err = p.findWhen(request)
		assert.True(t, err == nil)

		request = &ReminderRequest{
			TeamId:   model.NewId(),
			Username: user.Username,
			Payload:  "Hello at noon every monday",
			Reminder: Reminder{
				Id:        model.NewId(),
				TeamId:    model.NewId(),
				Username:  user.Username,
				Message:   "Hello",
				Completed: p.emptyTime,
				Target:    "me",
				When:      "at noon every monday",
			},
		}

		err = p.findWhen(request)
		assert.True(t, err == nil)

		request = &ReminderRequest{
			TeamId:   model.NewId(),
			Username: user.Username,
			Payload:  "Hello tomorrow",
			Reminder: Reminder{
				Id:        model.NewId(),
				TeamId:    model.NewId(),
				Username:  user.Username,
				Message:   "Hello",
				Completed: p.emptyTime,
				Target:    "me",
				When:      "tomorrow",
			},
		}

		err = p.findWhen(request)
		assert.True(t, err == nil)

	})

}
