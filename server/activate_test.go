package main

import (
	"testing"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEnsureBotExists(t *testing.T) {
	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("LogError", mock.Anything).Maybe()
		return api
	}

	t.Run("if remindbot already exists", func(t *testing.T) {
		t.Run("should find and return the existing bot ID", func(t *testing.T) {
			expectedBotId := model.NewId()

			api := setupAPI()
			api.On("CreateBot", mock.Anything).Return(nil, &model.AppError{})
			api.On("GetUserByUsername", "remindbot").Return(&model.User{
				Id: expectedBotId,
			}, nil)
			api.On("GetBot", expectedBotId, true).Return(&model.Bot{
				UserId: expectedBotId,
			}, nil)
			defer api.AssertExpectations(t)

			p := &Plugin{}
			p.API = api

			botId, err := p.ensureBotExists()

			assert.Equal(t, expectedBotId, botId)
			assert.Nil(t, err)
		})

		t.Run("should return an error if unable to get user", func(t *testing.T) {
			api := setupAPI()
			api.On("CreateBot", mock.Anything).Return(nil, &model.AppError{})
			api.On("GetUserByUsername", "remindbot").Return(nil, &model.AppError{})
			api.On("LogError", mock.Anything, "err", mock.Anything)
			defer api.AssertExpectations(t)

			p := &Plugin{}
			p.API = api

			botId, err := p.ensureBotExists()

			assert.Equal(t, "", botId)
			assert.NotNil(t, err)
		})

		t.Run("should return an error if unable to get bot", func(t *testing.T) {
			botUserId := model.NewId()

			api := setupAPI()
			api.On("CreateBot", mock.Anything).Return(nil, &model.AppError{})
			api.On("GetUserByUsername", "remindbot").Return(&model.User{
				Id: botUserId,
			}, nil)
			api.On("GetBot", botUserId, true).Return(nil, &model.AppError{})
			api.On("LogError", mock.Anything, "err", mock.Anything)
			defer api.AssertExpectations(t)

			p := &Plugin{}
			p.API = api

			botId, err := p.ensureBotExists()

			assert.Equal(t, "", botId)
			assert.NotNil(t, err)
		})
	})

	t.Run("if remindbot doesn't exist", func(t *testing.T) {
		t.Run("should create the bot and return the ID", func(t *testing.T) {
			expectedBotId := model.NewId()
			// profileImageBytes := []byte("profileImage")

			api := setupAPI()
			api.On("CreateBot", mock.Anything).Return(&model.Bot{
				UserId: expectedBotId,
			}, nil)
			// api.On("GetBundlePath").Return("", nil)
			// api.On("SetProfileImage", expectedBotId, profileImageBytes).Return(nil)
			defer api.AssertExpectations(t)

			p := &Plugin{
				// readFile: func(path string) ([]byte, error) {
				// 	return profileImageBytes, nil
				// },
			}
			p.API = api

			botId, err := p.ensureBotExists()

			assert.Equal(t, expectedBotId, botId)
			assert.Nil(t, err)
		})

		// t.Run("should log a warning if unable to set the profile picture, but still return the bot", func(t *testing.T) {
		// 	expectedBotId := model.NewId()

		// 	api := setupAPI()
		// 	api.On("CreateBot", mock.Anything).Return(&model.Bot{
		// 		UserId: expectedBotId,
		// 	}, nil)
		// 	api.On("GetBundlePath").Return("", &model.AppError{})
		// 	api.On("LogWarn", mock.Anything, "err", mock.Anything)
		// 	defer api.AssertExpectations(t)

		// 	p := &Plugin{}
		// 	p.API = api

		// 	botId, err := p.ensureBotExists()

		// 	assert.Equal(t, expectedBotId, botId)
		// 	assert.Nil(t, err)
		// })
	})
}

// func TestSetBotProfileImage(t *testing.T) {
// 	t.Run("should set profile image correctly", func(t *testing.T) {
// 		botUserId := model.NewId()
// 		profileImageBytes := []byte("profile image")

// 		api := &plugintest.API{}
// 		api.On("GetBundlePath").Return("/foo/bar", nil)
// 		api.On("SetProfileImage", botUserId, profileImageBytes).Return(nil)
// 		defer api.AssertExpectations(t)

// 		p := &Plugin{
// 			readFile: func(path string) ([]byte, error) {
// 				assert.Equal(t, "/foo/bar/assets/icon-happy-bot-square@1x.png", path)

// 				return profileImageBytes, nil
// 			},
// 		}
// 		p.API = api

// 		assert.Nil(t, p.setBotProfileImage(botUserId))
// 	})

// 	t.Run("should return an error when SetProfileImage fails", func(t *testing.T) {
// 		botUserId := model.NewId()
// 		profileImageBytes := []byte("profile image")

// 		api := &plugintest.API{}
// 		api.On("GetBundlePath").Return("/foo/bar", nil)
// 		api.On("SetProfileImage", botUserId, profileImageBytes).Return(&model.AppError{})
// 		defer api.AssertExpectations(t)

// 		p := &Plugin{
// 			readFile: func(path string) ([]byte, error) {
// 				assert.Equal(t, "/foo/bar/assets/icon-happy-bot-square@1x.png", path)

// 				return profileImageBytes, nil
// 			},
// 		}
// 		p.API = api

// 		assert.NotNil(t, p.setBotProfileImage(botUserId))
// 	})

// 	t.Run("should return an error when readFile fails", func(t *testing.T) {
// 		botUserId := model.NewId()

// 		api := &plugintest.API{}
// 		api.On("GetBundlePath").Return("/foo/bar", nil)
// 		defer api.AssertExpectations(t)

// 		p := &Plugin{
// 			readFile: func(path string) ([]byte, error) {
// 				return nil, &model.AppError{}
// 			},
// 		}
// 		p.API = api

// 		assert.NotNil(t, p.setBotProfileImage(botUserId))
// 	})

// 	t.Run("should return an error when GetBundlePath fails", func(t *testing.T) {
// 		botUserId := model.NewId()

// 		api := &plugintest.API{}
// 		api.On("GetBundlePath").Return("", &model.AppError{})
// 		defer api.AssertExpectations(t)

// 		p := &Plugin{}
// 		p.API = api

// 		assert.NotNil(t, p.setBotProfileImage(botUserId))
// 	})
// }
