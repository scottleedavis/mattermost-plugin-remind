package main

import (
	"fmt"
	"time"
	"encoding/json"

	"github.com/mattermost/mattermost-server/model"
)

type Reminder struct {

    Id string

	Username string

	Target string

	Message string

	When string

	Occurrences []ReminderOccurrence

	Completed time.Time
}

type ReminderRequest struct {
	Username string

	Payload string

	Reminder Reminder
}

func (p *Plugin) upsertReminder(request ReminderRequest) {

	user, u_err := p.API.GetUserByUsername(request.Username)
	
	if u_err != nil {
		p.API.LogError("failed to query user %s", request.Username)
		return
	}

	bytes, b_err := p.API.KVGet(user.Username)
	if b_err != nil {
		p.API.LogError("failed KVGet %s", b_err)
		return
	}

	var reminders []Reminder
	err := json.Unmarshal(bytes, &reminders)

	if err != nil {
		p.API.LogError("new reminder " + user.Username)
	} else {
		p.API.LogError("existing "+fmt.Sprintf("%v",reminders))
	}

	reminders = append(reminders, request.Reminder)
	ro,__ := json.Marshal(reminders)

	if __ != nil {
		p.API.LogError("failed to marshal reminders %s", user.Username)
		return
	}

	p.API.KVSet(user.Username,ro)
}

func (p *Plugin) triggerReminders() {

	bytes, err := p.API.KVGet(string(fmt.Sprintf("%v", time.Now().Round(time.Second))))

	if err != nil {
		p.API.LogError("failed KVGet %s", err)
	} else if string(bytes[:]) == "" {
	} else {

		var reminderOccurrences []ReminderOccurrence

		ro_err := json.Unmarshal(bytes, &reminderOccurrences)
		if ro_err != nil {
			p.API.LogError("Failed to unmarshal reminder occurrences "+ fmt.Sprintf("%v",ro_err))
			return
		}

		p.API.LogError("existing "+fmt.Sprintf("%v",reminderOccurrences))

		// TODO loop through array of occurrences, and trigger DM between remind user & user
		for _, ReminderOccurrence := range reminderOccurrences {

			user, err := p.API.GetUserByUsername(ReminderOccurrence.Username)
			
			if err != nil {
				p.API.LogError("failed to query user %s", user.Id)
			} else {
				channel, cerr := p.API.GetDirectChannel(p.remindUserId, user.Id)

				if cerr != nil {
					p.API.LogError("fail to get channel ", fmt.Sprintf("%v", cerr))
				} else {
					p.API.LogError("got channel "+ fmt.Sprintf("%v", channel))


						if _, err = p.API.CreatePost(&model.Post{
							UserId:    p.remindUserId,
							ChannelId: channel.Id,
							Message:   fmt.Sprintf(":wave: hello @%s", user.Username),
						}); err != nil {
							p.API.LogError(
								"failed to post DM message",
								"user_id", user.Id,
								"error", err.Error(),
							)
						}


				}
			}

		}

	}

}