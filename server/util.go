package main

import (
	"time"

	"github.com/mattermost/mattermost-server/model"
	"github.com/nicksnyder/go-i18n/i18n"
	"github.com/scottleedavis/go-timezone"
)

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
	tz_code := user.GetPreferredTimezone()
	if tz_code == "" {
		tz_code, _ = time.Now().Zone()
	}
	if tzLoc, err := timezone.GetTimezones(tz_code); err != nil {
		return time.Now().Location()
	} else {
		if l, lErr := time.LoadLocation(tzLoc[0]); lErr != nil {
			return time.Now().Location()
		} else {
			return l
		}

	}
}
