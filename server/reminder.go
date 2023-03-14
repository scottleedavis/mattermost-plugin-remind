package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
)

type Reminder struct {
	TeamId string

	Id string

	Username string

	Target string

	Message string

	When string

	Occurrences []Occurrence

	Completed time.Time
}

type ReminderRequest struct {
	TeamId string

	Username string

	Payload string

	Reminder Reminder
}

func (p *Plugin) TriggerReminders() {
	tickAt := time.Now().UTC().Round(time.Second)
	lastTickAt := p.getLastTickTimeWithDefault(tickAt)

	// Before handling more operations, save the updated LastTickAt time
	p.setLastTickTime(tickAt)

	// Catch up on missed ticks (if any)
	tickDelta := tickAt.Sub(lastTickAt)
	ticksMissed := tickDelta.Seconds() - 1
	if ticksMissed > 0 {
		oneSecond             := time.Second
		maxCatchupDuration, _ := time.ParseDuration("-10m")
		catchupStart          := lastTickAt.Add(oneSecond)
		earliestCatchupStart  := tickAt.Add(maxCatchupDuration)

		if (catchupStart.Before(earliestCatchupStart)) {
			catchupStart = earliestCatchupStart
			p.API.LogInfo(fmt.Sprintf("Too many reminder ticks were missed: occurrences between %v and %v will be dropped.", lastTickAt, catchupStart))
		}

		p.API.LogDebug(fmt.Sprintf("Catching up on %v reminder tick(s)...", tickAt.Sub(catchupStart).Seconds()))
		for tick := catchupStart; tick.Before(tickAt); tick = tick.Add(oneSecond) {
			p.TriggerRemindersForTick(tick)
		}
		p.API.LogDebug("Caught up on missed reminder ticks.")
	}

	// Trigger the actual tick
	p.TriggerRemindersForTick(tickAt)
}

func (p *Plugin) TriggerRemindersForTick(tickAt time.Time) {
	p.API.LogDebug("Trigger reminders for " + fmt.Sprintf("%v", tickAt))

	// Look up reminders to be triggered for the tick time
	bytes, err := p.API.KVGet(string(fmt.Sprintf("%v", tickAt)))
	if err != nil {
		p.API.LogError("failed KVGet %s", err)
	}

	if string(bytes[:]) != "" {

		var occurrences []Occurrence
		oErr := json.Unmarshal(bytes, &occurrences)
		if oErr != nil {
			p.API.LogError("Failed to unmarshal reminder occurrences " + fmt.Sprintf("%v", oErr))
			return
		}
		hostname, _ := os.Hostname()

		for _, occurrence := range occurrences {

			if hostname != occurrence.Hostname {
				continue
			}
			user, uErr := p.API.GetUserByUsername(occurrence.Username)
			if uErr != nil {
				p.API.LogError("failed to query user %s", user.Id)
				continue
			}
			bytes, bErr := p.API.KVGet(user.Username)
			if bErr != nil {
				p.API.LogError("failed KVGet %s", bErr)
				continue
			}
			var reminders []Reminder
			rsErr := json.Unmarshal(bytes, &reminders)
			if rsErr != nil {
				p.API.LogError("failed json Unmarshal %s", rsErr)
				continue
			}

			T, _ := p.translation(user)
			reminder := p.findReminder(reminders, occurrence)

			if reminder.Target == "" {
				p.API.LogError("No target found for reminder")
				continue
			}

			if strings.HasPrefix(reminder.Target, "@") || strings.HasPrefix(reminder.Target, T("me")) { //@user
				var targetId string
				if strings.HasPrefix(reminder.Target, "@") {
					target := strings.Trim(reminder.Target, "@")
					targetUser, tErr := p.API.GetUserByUsername(target)
					if tErr != nil {
						p.API.LogError("failed to find target user " + reminder.Target)
						continue
					}
					targetId = targetUser.Id
				} else {
					targetId = user.Id
				}

				channel, cErr := p.API.GetDirectChannel(p.remindUserId, targetId)
				if cErr != nil {
					p.API.LogError("failed to create channel " + cErr.Error())
					continue
				}

				finalTarget := reminder.Target
				if finalTarget == T("me") {
					finalTarget = T("you")
				} else {
					finalTarget = "@" + user.Username
				}

				messageParameters := map[string]interface{}{
					"FinalTarget": finalTarget,
					"Message":     reminder.Message,
				}

				interactivePost := model.Post{}

				if occurrence.Repeat == "" {

					interactivePost = model.Post{
						ChannelId:     channel.Id,
						PendingPostId: model.NewId() + ":" + fmt.Sprint(model.GetMillis()),
						UserId:        p.remindUserId,
						Props: model.StringInterface{
							"attachments": []*model.SlackAttachment{
								{
									Text: T("reminder.message", messageParameters),
									Actions: []*model.PostAction{
										{
											Id: model.NewId(),
											Integration: &model.PostActionIntegration{
												Context: model.StringInterface{
													"orig_user_id":  user.Id,
													"reminder_id":   reminder.Id,
													"occurrence_id": occurrence.Id,
													"action":        "complete",
												},
												URL: fmt.Sprintf("/plugins/%s/complete", manifest.ID),
											},
											Type: model.POST_ACTION_TYPE_BUTTON,
											Name: T("button.complete"),
										},
										{
											Integration: &model.PostActionIntegration{
												Context: model.StringInterface{
													"orig_user_id":  user.Id,
													"reminder_id":   reminder.Id,
													"occurrence_id": occurrence.Id,
													"action":        "delete",
												},
												URL: fmt.Sprintf("/plugins/%s/delete", manifest.ID),
											},
											Name: T("button.delete"),
											Type: "action",
										},
										{
											Integration: &model.PostActionIntegration{
												Context: model.StringInterface{
													"orig_user_id":  user.Id,
													"reminder_id":   reminder.Id,
													"occurrence_id": occurrence.Id,
													"action":        "snooze",
												},
												URL: fmt.Sprintf("/plugins/%s/snooze", manifest.ID),
											},
											Name: T("button.snooze"),
											Type: "select",
											Options: []*model.PostActionOptions{
												{
													Text:  T("button.snooze.20min"),
													Value: "20min",
												},
												{
													Text:  T("button.snooze.1hr"),
													Value: "1hr",
												},
												{
													Text:  T("button.snooze.3hr"),
													Value: "3hrs",
												},
												{
													Text:  T("button.snooze.tomorrow"),
													Value: "tomorrow",
												},
												{
													Text:  T("button.snooze.nextweek"),
													Value: "nextweek",
												},
											},
										},
									},
								},
							},
						},
					}

				} else {

					interactivePost = model.Post{
						ChannelId:     channel.Id,
						PendingPostId: model.NewId() + ":" + fmt.Sprint(model.GetMillis()),
						UserId:        p.remindUserId,
						Props: model.StringInterface{
							"attachments": []*model.SlackAttachment{
								{
									Text: T("reminder.message", messageParameters),
									Actions: []*model.PostAction{
										{
											Integration: &model.PostActionIntegration{
												Context: model.StringInterface{
													"reminder_id":   reminder.Id,
													"occurrence_id": occurrence.Id,
													"action":        "snooze",
												},
												URL: fmt.Sprintf("/plugins/%s/snooze", manifest.ID),
											},
											Name: T("button.snooze"),
											Type: "select",
											Options: []*model.PostActionOptions{
												{
													Text:  T("button.snooze.20min"),
													Value: "20min",
												},
												{
													Text:  T("button.snooze.1hr"),
													Value: "1hr",
												},
												{
													Text:  T("button.snooze.3hr"),
													Value: "3hrs",
												},
												{
													Text:  T("button.snooze.tomorrow"),
													Value: "tomorrow",
												},
												{
													Text:  T("button.snooze.nextweek"),
													Value: "nextweek",
												},
											},
										},
									},
								},
							},
						},
					}
				}

				if _, pErr := p.API.CreatePost(&interactivePost); pErr != nil {
					p.API.LogError(fmt.Sprintf("%v", pErr))
					continue
				}

				if occurrence.Repeat != "" {
					defer p.rescheduleOccurrence(&occurrence)
				}

			} else if strings.HasPrefix(reminder.Target, "~") { //~ channel

				channel, cErr := p.API.GetChannelByName(
					reminder.TeamId,
					strings.Replace(reminder.Target, "~", "", -1),
					false,
				)

				if cErr != nil {
					p.API.LogError("fail to get channel " + fmt.Sprintf("%v", cErr))
				} else {

					var messageParameters = map[string]interface{}{
						"FinalTarget": "@" + user.Username,
						"Message":     reminder.Message,
					}

					interactivePost := model.Post{
						ChannelId:     channel.Id,
						PendingPostId: model.NewId() + ":" + fmt.Sprint(model.GetMillis()),
						UserId:        p.remindUserId,
						Type:          model.POST_CUSTOM_TYPE_PREFIX + "reminder",
						Message:       T("reminder.message", messageParameters),
						Props:         model.StringInterface{},
					}

					if _, pErr := p.API.CreatePost(&interactivePost); pErr != nil {
						p.API.LogError(fmt.Sprintf("%v", pErr))
					}

					if occurrence.Repeat != "" {
						defer p.rescheduleOccurrence(&occurrence)
					}

				}
			}

		}

	}

}

func (p *Plugin) GetReminder(userId string, reminderId string) Reminder {

	user, uErr := p.API.GetUser(userId)
	if uErr != nil {
		return Reminder{}
	}

	reminders := p.GetReminders(user.Username)
	for _, reminder := range reminders {
		if reminder.Id == reminderId {
			return reminder
		}
	}

	return Reminder{}
}

func (p *Plugin) GetReminders(username string) []Reminder {

	user, uErr := p.API.GetUserByUsername(username)

	if uErr != nil {
		p.API.LogError("failed to query user " + username)
		return []Reminder{}
	}

	bytes, bErr := p.API.KVGet(user.Username)
	if bErr != nil {
		p.API.LogError("failed KVGet " + bErr.Error())
		return []Reminder{}
	}

	var reminders []Reminder
	err := json.Unmarshal(bytes, &reminders)

	if err != nil {
		p.API.LogError("no reminders for " + user.Username)
	}
	return reminders
}

func (p *Plugin) UpdateReminder(userId string, reminder Reminder) error {

	user, uErr := p.API.GetUser(userId)

	if uErr != nil {
		p.API.LogError("failed to query user %s", user.Username)
		return uErr
	}

	bytes, bErr := p.API.KVGet(user.Username)
	if bErr != nil {
		p.API.LogError("failed KVGet %s", bErr)
		return bErr
	}

	var reminders []Reminder
	if err := json.Unmarshal(bytes, &reminders); err != nil {
		return err
	}

	updatedReminders := []Reminder{}
	for _, r := range reminders {
		if r.Id == reminder.Id {
			updatedReminders = append(updatedReminders, reminder)
		} else {
			updatedReminders = append(updatedReminders, r)
		}
	}

	ro, rErr := json.Marshal(updatedReminders)

	if rErr != nil {
		p.API.LogError("failed to marshal reminders %s", user.Username)
		return rErr
	}

	p.API.KVSet(user.Username, ro)

	return nil
}

func (p *Plugin) UpsertReminder(request *ReminderRequest) error {

	user, uErr := p.API.GetUserByUsername(request.Username)

	if uErr != nil {
		p.API.LogError("failed to query user %s", request.Username)
		return uErr
	}

	bytes, bErr := p.API.KVGet(user.Username)
	if bErr != nil {
		p.API.LogError("failed KVGet %s", bErr)
		return bErr
	}

	var reminders []Reminder
	err := json.Unmarshal(bytes, &reminders)
	if err != nil {
		p.API.LogDebug("new reminder " + user.Username)
	}

	duplicateReminder := false
	// a feature requested by PM.   seems to create undesired results
	//for _, r := range reminders {
	//	r.Id = request.Reminder.Id
	//	if r.Username == request.Reminder.Username &&
	//		r.Target == request.Reminder.Target &&
	//		r.Message == request.Reminder.Message &&
	//		r.When == request.Reminder.When {
	//		duplicateReminder = true
	//		break
	//	}
	//}
	if !duplicateReminder {
		reminders = append(reminders, request.Reminder)
	}
	ro, rErr := json.Marshal(reminders)
	if rErr != nil {
		p.API.LogError("failed to marshal reminders %s", user.Username)
		return rErr
	}

	p.API.KVSet(user.Username, ro)

	return nil
}

func (p *Plugin) DeleteReminder(userId string, reminder Reminder) error {

	user, uErr := p.API.GetUser(userId)

	if uErr != nil {
		p.API.LogError("failed to query user %s", user.Username)
		return uErr
	}

	bytes, bErr := p.API.KVGet(user.Username)
	if bErr != nil {
		p.API.LogError("failed KVGet %s", bErr)
		return bErr
	}

	var reminders []Reminder
	if err := json.Unmarshal(bytes, &reminders); err != nil {
		return err
	}

	var prunedReminders []Reminder
	for _, r := range reminders {
		if r.Id != reminder.Id {
			prunedReminders = append(prunedReminders, r)
		}
	}

	ro, rErr := json.Marshal(prunedReminders)

	if rErr != nil {
		p.API.LogError("failed to marshal reminders %s", user.Username)
		return rErr
	}

	p.API.KVSet(user.Username, ro)

	return nil
}

func (p *Plugin) DeleteReminders(user *model.User) string {
	T, _ := p.translation(user)

	bytes, bErr := p.API.KVGet(user.Username)
	if bErr != nil {
		p.API.LogError("failed KVGet %s", bErr)
		return T("exception.response")
	}

	var reminders []Reminder
	if err := json.Unmarshal(bytes, &reminders); err != nil {
		p.API.LogError("failed json unmarshal %s", bErr)
		return T("exception.reponse")
	}

	for _, r := range reminders {
		for _, o := range r.Occurrences {
			if o.Snoozed != p.emptyTime && r.Completed != p.emptyTime {
				p.deleteSnoozedOccurrence(o)
			} else if o.Occurrence != p.emptyTime && r.Completed != p.emptyTime {
				p.deleteOccurrence(o)
			}
		}
	}

	dErr := p.API.KVDelete(user.Username)
	if dErr != nil {
		p.API.LogError("failed KVDelete %s", dErr)
		return T("exception.response")
	}
	return T("clear.response")
}

func (p *Plugin) rescheduleOccurrence(occurrence *Occurrence) {

	time.Sleep(1000 * time.Millisecond)

	user, _ := p.API.GetUserByUsername(occurrence.Username)
	_, locale := p.translation(user)

	var times []time.Time

	switch locale {
	case "en":
		times, _ = p.rescheduleOccurrenceEN(occurrence)
	default:
		times, _ = p.rescheduleOccurrenceEN(occurrence)
	}

	if len(times) > 1 {
		for _, ts := range times {
			if ts.Weekday() == occurrence.Occurrence.Weekday() {
				occurrence.Occurrence = ts
				p.upsertOccurrence(occurrence)
				return
			}
		}
	} else {
		occurrence.Occurrence = times[0]
		p.upsertOccurrence(occurrence)
	}

	bytes, bErr := p.API.KVGet(user.Username)
	if bErr != nil {
		p.API.LogError("failed KVGet %s", bErr)
	}
	var reminders []Reminder
	rsErr := json.Unmarshal(bytes, &reminders)
	if rsErr != nil {
		p.API.LogError("failed json Unmarshal %s", rsErr)
	}

	reminder := p.findReminder(reminders, *occurrence)

	updatedOccurrences := []Occurrence{}
	for _, o := range reminder.Occurrences {
		if o.Id == occurrence.Id {
			updatedOccurrences = append(updatedOccurrences, *occurrence)
		} else {
			updatedOccurrences = append(updatedOccurrences, o)
		}
	}
	reminder.Occurrences = updatedOccurrences
	p.UpdateReminder(user.Id, reminder)

}

func (p *Plugin) rescheduleOccurrenceEN(occurrence *Occurrence) ([]time.Time, error) {

	user, _ := p.API.GetUserByUsername(occurrence.Username)
	T, _ := p.translation(user)

	if strings.HasPrefix(occurrence.Repeat, T("in")) {
		return p.in(occurrence.Repeat, user)
	} else if strings.HasPrefix(occurrence.Repeat, T("at")) {
		return p.at(occurrence.Repeat, user)
	} else if strings.HasPrefix(occurrence.Repeat, T("on")) {
		return p.on(occurrence.Repeat, user)
	} else if strings.HasPrefix(occurrence.Repeat, T("every")) {
		return p.every(occurrence.Repeat, user)
	} else {
		return p.freeForm(occurrence.Repeat, user)
	}
}

func (p *Plugin) findReminder(reminders []Reminder, occurrence Occurrence) Reminder {
	for _, reminder := range reminders {
		if reminder.Id == occurrence.ReminderId {
			return reminder
		}
	}
	return Reminder{}
}

func (p *Plugin) getLastTickTimeWithDefault(defaultValue time.Time) time.Time {
	bytes, err := p.API.KVGet("LastTickAt")
	if err != nil {
		p.API.LogInfo(fmt.Sprintf("Failed to read LastTickAt (%v); returning the default value", err))
		return defaultValue
	}
	if bytes == nil {
		p.API.LogDebug("LastTickAt is not set; returning the default value")
		return defaultValue
	}

	lastTickAt, parseErr := time.Parse(time.RFC3339, string(bytes[:]))
	if parseErr != nil {
		p.API.LogInfo(fmt.Sprintf("Failed to parse LastTickAt value (%v); returning the default value", parseErr))
		return defaultValue
	}

	return lastTickAt
}

func (p *Plugin) setLastTickTime(lastTickAt time.Time) {
	serializedTime := lastTickAt.Format(time.RFC3339)
	p.API.KVSet("LastTickAt", []byte(serializedTime))
}
