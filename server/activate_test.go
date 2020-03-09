package main

import (
	"testing"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEnsureBotExists(t *testing.T) {
	bot := &model.Bot{
		Username:    "remindbot",
		DisplayName: "Remindbot",
	}

	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("LogError", mock.Anything).Maybe()
		return api
	}

	t.Run("EnsureBotOption passes", func(t *testing.T) {
		expectedBotId := model.NewId()

		api := setupAPI()
		defer api.AssertExpectations(t)

		helpers := &plugintest.Helpers{}
		helpers.On("EnsureBot", bot, mock.AnythingOfType("plugin.EnsureBotOption")).Return(expectedBotId, nil)
		defer helpers.AssertExpectations(t)

		p := &Plugin{}
		p.SetAPI(api)
		p.SetHelpers(helpers)

		botId, err := p.ensureBotExists()
		assert.Equal(t, expectedBotId, botId)
		assert.NoError(t, err)
	})

	t.Run("EnsureBotOption fails", func(t *testing.T) {
		api := setupAPI()
		defer api.AssertExpectations(t)

		helpers := &plugintest.Helpers{}
		helpers.On("EnsureBot", bot, mock.AnythingOfType("plugin.EnsureBotOption")).Return("", errors.New(""))
		defer helpers.AssertExpectations(t)

		p := &Plugin{}
		p.SetAPI(api)
		p.SetHelpers(helpers)

		botId, err := p.ensureBotExists()
		assert.Equal(t, "", botId)
		assert.Error(t, err)
	})

}
