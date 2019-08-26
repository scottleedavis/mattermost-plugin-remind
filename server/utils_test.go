package main

import (
	"testing"

	"github.com/mattermost/mattermost-server/model"
	"github.com/stretchr/testify/assert"
)

func TestTranslation(t *testing.T) {

	t.Run("if can translate with locale", func(t *testing.T) {

		p := &Plugin{}

		user := &model.User{
			Email:    "-@-.-",
			Nickname: "TestUser",
			Password: model.NewId(),
			Username: "testuser",
			Roles:    model.SYSTEM_USER_ROLE_ID,
			Locale:   "en",
		}

		p.locales = make(map[string]string)
		p.locales["en"] = "/test/en.json"
		T, locale := p.translation(user)
		assert.Equal(t, T("me"), "me")
		assert.Equal(t, locale, "en")
	})
}

func TestLocation(t *testing.T) {

	t.Run("if can locate", func(t *testing.T) {

		p := &Plugin{}
		timezone := make(map[string]string)
		timezone["manualTimezone"] = "America/Los_Angeles"
		user := &model.User{
			Email:    "-@-.-",
			Nickname: "TestUser",
			Password: model.NewId(),
			Username: "testuser",
			Roles:    model.SYSTEM_USER_ROLE_ID,
			Locale:   "en",
			Timezone: timezone,
		}

		location := p.location(user)
		assert.NotNil(t, location)

	})
}
