package main

import (
	"strings"
	"fmt"
	"time"
)

func (p *Plugin) ListReminders(username string) (string) {

	reminders := p.GetReminders(username)
	//fmt.Println("Same, in UTC:", t.UTC().Format(time.UnixDate))

	var output string
	output = ""
	for _, reminder := range reminders {

		if len(reminder.Occurrences) > 0 {

			for _, occurrence := range reminder.Occurrences {
				if !occurrence.Occurrence.Equal(time.Time{}) {
					output = strings.Join([]string{output, "* \"" + reminder.Message + "\" " + fmt.Sprintf("%v", occurrence.Occurrence.UTC().Format(time.UnixDate))}, "\n")
				}
			}

		}
	}

	// TODO categorize the set and group output.  Following same pattern at mattermost-remind
	//"*Upcoming*:\n"
	//"*Recurring*:
	//"*Past and incomplete*:"

	return output + "\n*Note*:  To interact with these reminders use `/remind list` in a direct message with the remind user";
	return "foo"
}
