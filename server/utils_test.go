package main

import (
	"testing"

	"github.com/mattermost/mattermost-server/model"
	"github.com/stretchr/testify/assert"
)

func TestTranslation(t *testing.T) {

	t.Run("if can translate", func(t *testing.T) {

		p := &Plugin{}

		user := &model.User{
			Email:    "-@-.-",
			Nickname: "TestUser",
			Password: model.NewId(),
			Username: "testuser",
			Roles:    model.SYSTEM_USER_ROLE_ID,
			Locale:   "en",
		}

		T, locale := p.translation(user)
		assert.Equal(t, T("me"), "me")
		assert.Equal(t, locale, "en")

	})
}

func TestLocation(t *testing.T) {

	t.Run("if can locate", func(t *testing.T) {

		p := &Plugin{}

		user := &model.User{
			Email:    "-@-.-",
			Nickname: "TestUser",
			Password: model.NewId(),
			Username: "testuser",
			Roles:    model.SYSTEM_USER_ROLE_ID,
			Locale:   "en",
		}

		location := p.location(user)
		assert.NotNil(t, location)

	})
}