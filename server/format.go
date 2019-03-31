package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/model"
)

func (p *Plugin) ListReminders(user *model.User) string {

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

	return output + "\n*Note*:  To interact with these reminders use `/remind list` in a direct message with the remind user"
}
