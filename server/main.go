package main

import (
	"time"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/nicksnyder/go-i18n/i18n"
)

func main() {
	plugin.ClientMain(&Plugin{})
}

func (p *Plugin) translation(user *model.User) (i18n.TranslateFunc, string) {
	locale := "en"
	for _, l := range p.supportedLocales {
		if user.Locale == l {
			locale = user.Locale
			break
		}
	}
	return GetUserTranslations(locale), locale
}

func (p *Plugin) location(user *model.User) *time.Location {
	timezone := user.GetPreferredTimezone()
	if timezone == "" {
		timezone, _ = time.Now().Zone()
	}
	location, _ := time.LoadLocation(timezone)
	return location
}
