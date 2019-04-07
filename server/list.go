package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/model"
)

// TODO
func (p *Plugin) ListReminders(user *model.User, channelId string) string {

	reminders := p.GetReminders(user.Username)

	var output string
	output = ""
	for _, reminder := range reminders {

		if len(reminder.Occurrences) > 0 {

			location, _ := time.LoadLocation(user.Timezone["automaticTimezone"])
			for _, occurrence := range reminder.Occurrences {
				if !occurrence.Occurrence.Equal(time.Time{}) {
					output = strings.Join([]string{output, "* \"" + reminder.Message + "\" " + fmt.Sprintf("%v", occurrence.Occurrence.In(location).Format(time.UnixDate))}, "\n")
				}
			}

		}
	}

	return output + "test"
}
