package main

import (
	"fmt"
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
	p.API.LogInfo("1) TIMEZONE TIMEZONE TIMEZONE : " + timezone)
	if timezone == "" {
		timezone, _ = time.Now().Zone()
	}
	location, _ := time.LoadLocation("America/New_York") // ("timezone")
	p.API.LogInfo("2) TIMEZONE TIMEZONE TIMEZONE : " + timezone)
	p.API.LogInfo("location: " + location.String())
	p.API.LogInfo(fmt.Sprintf("%v", time.Now().In(location)))

	return location
}
