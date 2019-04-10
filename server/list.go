package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/model"
)

const RemindersPerPage = 2 //3

func (p *Plugin) ListReminders(user *model.User, channelId string) string {

	offset := 1
	reminders := p.GetReminders(user.Username)

	if len(reminders) == 0 {
		return T("no.reminders")
	}

	upcomingOccurrences,
		recurringOccurrences,
		pastOccurrences,
		channelOccurrences := p.categorizeOccurrences(reminders)

	attachments, endOffset := p.pagedOccurrences(
		user,
		reminders,
		upcomingOccurrences,
		recurringOccurrences,
		pastOccurrences,
		channelOccurrences,
		offset)

	completedCount, attachments := p.completedReminders(reminders, attachments)
	attachments = p.pageControl(len(reminders), completedCount, offset, endOffset, attachments)

	channel, cErr := p.API.GetDirectChannel(p.remindUserId, user.Id)
	if cErr != nil {
		p.API.LogError("failed to create channel " + cErr.Error())
		return ""
	}

	interactivePost := model.Post{
		ChannelId:     channel.Id,
		PendingPostId: model.NewId() + ":" + fmt.Sprint(model.GetMillis()),
		UserId:        p.remindUserId,
		Props: model.StringInterface{
			"attachments": attachments,
		},
	}

	defer p.API.CreatePost(&interactivePost)

	if channel.Id == channelId {
		return ""
	} else {
		messageParameters := map[string]interface{}{
			"RemindUser": CommandTrigger,
		}
		return T("list.reminders", messageParameters)
	}
}

func (p *Plugin) UpdateListReminders(userId string, postId string, offset int) {

	user, _ := p.API.GetUser(userId)
	reminders := p.GetReminders(user.Username)

	upcomingOccurrences,
		recurringOccurrences,
		pastOccurrences,
		channelOccurrences := p.categorizeOccurrences(reminders)

	attachments, endOffset := p.pagedOccurrences(
		user,
		reminders,
		upcomingOccurrences,
		recurringOccurrences,
		pastOccurrences,
		channelOccurrences,
		offset)

	completedCount, attachments := p.completedReminders(reminders, attachments)
	attachments = p.pageControl(len(reminders), completedCount, offset, endOffset, attachments)

	if post, pErr := p.API.GetPost(postId); pErr != nil {
		p.API.LogError("unable to get list reminders post " + pErr.Error())
	} else {
		post.Props = model.StringInterface{
			"attachments": attachments,
		}
		defer p.API.UpdatePost(post)
	}
}

func (p *Plugin) categorizeOccurrences(reminders []Reminder) (
	upcomingOccurrences []Occurrence,
	recurringOccurrences []Occurrence,
	pastOccurrences []Occurrence,
	channelOccurrences []Occurrence) {

	for _, reminder := range reminders {
		occurrences := reminder.Occurrences
		if len(occurrences) > 0 {
			for _, occurrence := range occurrences {
				t := occurrence.Occurrence
				s := occurrence.Snoozed

				if !strings.HasPrefix(reminder.Target, "~") &&
					reminder.Completed == p.emptyTime &&
					((occurrence.Repeat == "" && t.After(time.Now())) ||
						(s != p.emptyTime && s.After(time.Now()))) {
					upcomingOccurrences = append(upcomingOccurrences, occurrence)
				}

				if !strings.HasPrefix(reminder.Target, "~") &&
					occurrence.Repeat != "" && t.After(time.Now()) {
					recurringOccurrences = append(recurringOccurrences, occurrence)
				}

				if !strings.HasPrefix(reminder.Target, "~") &&
					reminder.Completed == p.emptyTime &&
					t.Before(time.Now()) &&
					s == p.emptyTime {
					pastOccurrences = append(pastOccurrences, occurrence)
				}

				if strings.HasPrefix(reminder.Target, "~") &&
					reminder.Completed == p.emptyTime &&
					t.After(time.Now()) {
					channelOccurrences = append(channelOccurrences, occurrence)
				}

			}
		}
	}

	return upcomingOccurrences, recurringOccurrences, pastOccurrences, channelOccurrences

}

func (p *Plugin) pagedOccurrences(
	user *model.User,
	reminders []Reminder,
	upcomingOccurrences []Occurrence,
	recurringOccurrences []Occurrence,
	pastOccurrences []Occurrence,
	channelOccurrences []Occurrence,
	offset int) (attachments []*model.SlackAttachment, endOffset int) {

	offsetCount := 1
	endOffset = offset + RemindersPerPage

	if len(upcomingOccurrences) > 0 {
		for _, o := range upcomingOccurrences {
			offsetCount += 1
			if offsetCount > offset && offsetCount <= endOffset {
				attachments = append(attachments, p.addAttachment(user, o, reminders, "upcoming"))
			} else {
				break
			}
		}
	}

	if len(recurringOccurrences) > 0 {
		for _, o := range recurringOccurrences {
			offsetCount += 1
			if offsetCount > offset && offsetCount <= endOffset {
				attachments = append(attachments, p.addAttachment(user, o, reminders, "recurring"))
			} else {
				break
			}
		}
	}

	if len(pastOccurrences) > 0 {
		for _, o := range pastOccurrences {
			offsetCount += 1
			if offsetCount > offset && offsetCount <= endOffset {
				attachments = append(attachments, p.addAttachment(user, o, reminders, "past"))
			} else {
				break
			}
		}
	}
	if len(channelOccurrences) > 0 {
		for _, o := range channelOccurrences {
			offsetCount += 1
			if offsetCount > offset && offsetCount <= endOffset {
				attachments = append(attachments, p.addAttachment(user, o, reminders, "channel"))
			} else {
				break
			}
		}
	}

	return attachments, endOffset

}

func (p *Plugin) completedReminders(reminders []Reminder, attachments []*model.SlackAttachment) (
	completedCount int,
	outputAttachments []*model.SlackAttachment) {

	completedCount = 0
	for _, reminder := range reminders {
		if reminder.Completed != p.emptyTime {
			completedCount += 1
		}
	}
	if completedCount > 0 {
		siteURL := fmt.Sprintf("%s", *p.ServerConfig.ServiceSettings.SiteURL)
		attachments = append(attachments, &model.SlackAttachment{
			Actions: []*model.PostAction{
				{
					Integration: &model.PostActionIntegration{
						Context: model.StringInterface{
							"action": "view/complete/list",
						},
						URL: fmt.Sprintf("%s/plugins/%s", siteURL, manifest.Id),
					},
					Type: model.POST_ACTION_TYPE_BUTTON,
					Name: T("button.view.complete"),
				},
				{
					Integration: &model.PostActionIntegration{
						Context: model.StringInterface{
							"action": "delete/complete/list",
						},
						URL: fmt.Sprintf("%s/plugins/%s", siteURL, manifest.Id),
					},
					Name: T("button.delete.complete"),
					Type: "action",
				},
				{
					Integration: &model.PostActionIntegration{
						Context: model.StringInterface{
							"action": "close/list",
						},
						URL: fmt.Sprintf("%s/plugins/%s", siteURL, manifest.Id),
					},
					Name: T("button.close.list"),
					Type: "action",
				},
			},
		})
	}

	outputAttachments = attachments
	return completedCount, outputAttachments
}

// attachments = p.pageControl(len(reminders), completedCount, offset, endOffset, attachments)
func (p *Plugin) pageControl(reminderCount int, completedCount int, offset int, endOffset int, attachments []*model.SlackAttachment) (
	outputAttachments []*model.SlackAttachment) {

	activeReminderCount := reminderCount - completedCount

	if activeReminderCount > RemindersPerPage {
		siteURL := fmt.Sprintf("%s", *p.ServerConfig.ServiceSettings.SiteURL)
		reminderCount := map[string]interface{}{
			"ReminderCount": RemindersPerPage,
		}
		remindersPageCount := map[string]interface{}{
			"Start": offset,
			"Stop":  RemindersPerPage,
			"Total": activeReminderCount,
		}

		if offset <= 1 {
			attachments = append(attachments, &model.SlackAttachment{
				Text: T("reminders.page.numbers", remindersPageCount),
				Actions: []*model.PostAction{
					{
						Integration: &model.PostActionIntegration{
							Context: model.StringInterface{
								"action": "next/reminders",
								"offset": endOffset,
							},
							URL: fmt.Sprintf("%s/plugins/%s", siteURL, manifest.Id),
						},
						Type: model.POST_ACTION_TYPE_BUTTON,
						Name: T("button.next.reminders", reminderCount),
					},
					{
						Integration: &model.PostActionIntegration{
							Context: model.StringInterface{
								"action": "close/list",
							},
							URL: fmt.Sprintf("%s/plugins/%s", siteURL, manifest.Id),
						},
						Name: T("button.close.list"),
						Type: "action",
					},
				},
			})
		} else {
			attachments = append(attachments, &model.SlackAttachment{
				Text: T("reminders.page.numbers", remindersPageCount),
				Actions: []*model.PostAction{
					{
						Integration: &model.PostActionIntegration{
							Context: model.StringInterface{
								"action": "previous/reminders",
								"offset": endOffset - RemindersPerPage,
							},
							URL: fmt.Sprintf("%s/plugins/%s", siteURL, manifest.Id),
						},
						Type: model.POST_ACTION_TYPE_BUTTON,
						Name: T("button.previous.reminders", reminderCount),
					},
					{
						Integration: &model.PostActionIntegration{
							Context: model.StringInterface{
								"action": "next/reminders",
								"offset": endOffset,
							},
							URL: fmt.Sprintf("%s/plugins/%s", siteURL, manifest.Id),
						},
						Type: model.POST_ACTION_TYPE_BUTTON,
						Name: T("button.next.reminders", reminderCount),
					},
					{
						Integration: &model.PostActionIntegration{
							Context: model.StringInterface{
								"action": "close/list",
							},
							URL: fmt.Sprintf("%s/plugins/%s", siteURL, manifest.Id),
						},
						Name: T("button.close.list"),
						Type: "action",
					},
				},
			})
		}

	}

	outputAttachments = attachments
	return outputAttachments
}

func (p *Plugin) addAttachment(user *model.User, occurrence Occurrence, reminders []Reminder, gType string) *model.SlackAttachment {

	location := p.location(user)
	T, _ := p.translation(user)

	siteURL := fmt.Sprintf("%s", *p.ServerConfig.ServiceSettings.SiteURL)
	reminder := p.findReminder(reminders, occurrence)

	t := occurrence.Occurrence
	s := occurrence.Snoozed

	formattedOccurrence := ""
	formattedOccurrence = p.formatWhen(user.Username, reminder.When, t.In(location).Format(time.RFC3339), false)

	formattedSnooze := ""
	if s != p.emptyTime {
		formattedSnooze = p.formatWhen(user.Username, reminder.When, s.In(location).Format(time.RFC3339), true)
	}
	var messageParameters = map[string]interface{}{
		"Message":    reminder.Message,
		"Occurrence": formattedOccurrence,
		"Snoozed":    formattedSnooze,
	}
	if !t.Equal(p.emptyTime) {
		switch gType {
		case "upcoming":
			output := ""
			if formattedSnooze == "" {
				output = T("list.upcoming") + " " + T("list.element.upcoming", messageParameters)
			} else {
				output = T("list.upcoming") + " " + T("list.element.upcoming.snoozed", messageParameters)
			}
			return &model.SlackAttachment{
				Text: output,
				Actions: []*model.PostAction{
					{
						Integration: &model.PostActionIntegration{
							Context: model.StringInterface{
								"reminder_id":   reminder.Id,
								"occurrence_id": occurrence.Id,
								"action":        "complete/list",
							},
							URL: fmt.Sprintf("%s/plugins/%s", siteURL, manifest.Id),
						},
						Type: model.POST_ACTION_TYPE_BUTTON,
						Name: T("button.complete"),
					},
					{
						Integration: &model.PostActionIntegration{
							Context: model.StringInterface{
								"reminder_id":   reminder.Id,
								"occurrence_id": occurrence.Id,
								"action":        "delete/list",
							},
							URL: fmt.Sprintf("%s/plugins/%s", siteURL, manifest.Id),
						},
						Name: T("button.delete"),
						Type: "action",
					},
				},
			}
		case "recurring":
			output := ""
			if formattedSnooze == "" {
				output = T("list.recurring") + " " + T("list.element.recurring", messageParameters)
			} else {
				output = T("list.recurring") + " " + T("list.element.recurring.snoozed", messageParameters)
			}
			return &model.SlackAttachment{
				Text: output,
				Actions: []*model.PostAction{
					{
						Integration: &model.PostActionIntegration{
							Context: model.StringInterface{
								"reminder_id":   reminder.Id,
								"occurrence_id": occurrence.Id,
								"action":        "delete/list",
							},
							URL: fmt.Sprintf("%s/plugins/%s", siteURL, manifest.Id),
						},
						Name: T("button.delete"),
						Type: "action",
					},
				},
			}
		case "past":
			output := ""
			if formattedSnooze == "" {
				output = T("list.past.and.incomplete") + " " + T("list.element.past", messageParameters)
			} else {
				output = T("list.past.and.incomplete") + " " + T("list.element.past.snoozed", messageParameters)
			}
			return &model.SlackAttachment{
				Text: output,
				Actions: []*model.PostAction{
					{
						Integration: &model.PostActionIntegration{
							Context: model.StringInterface{
								"reminder_id":   reminder.Id,
								"occurrence_id": occurrence.Id,
								"action":        "complete/list",
							},
							URL: fmt.Sprintf("%s/plugins/%s", siteURL, manifest.Id),
						},
						Type: model.POST_ACTION_TYPE_BUTTON,
						Name: T("button.complete"),
					},
					{
						Integration: &model.PostActionIntegration{
							Context: model.StringInterface{
								"reminder_id":   reminder.Id,
								"occurrence_id": occurrence.Id,
								"action":        "delete/list",
							},
							URL: fmt.Sprintf("%s/plugins/%s", siteURL, manifest.Id),
						},
						Name: T("button.delete"),
						Type: "action",
					},
					{
						Integration: &model.PostActionIntegration{
							Context: model.StringInterface{
								"reminder_id":   reminder.Id,
								"occurrence_id": occurrence.Id,
								"action":        "snooze/list",
							},
							URL: fmt.Sprintf("%s/plugins/%s", siteURL, manifest.Id),
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
			}
		case "channel":
			output := ""
			if formattedSnooze == "" {
				output = T("list.channel") + " " + T("list.element.channel", messageParameters)
			} else {
				output = T("list.channel") + " " + T("list.element.channel.snoozed", messageParameters)
			}
			return &model.SlackAttachment{
				Text: output,
				Actions: []*model.PostAction{
					{
						Integration: &model.PostActionIntegration{
							Context: model.StringInterface{
								"reminder_id":   reminder.Id,
								"occurrence_id": occurrence.Id,
								"action":        "delete/list",
							},
							URL: fmt.Sprintf("%s/plugins/%s/api/v1/delete", siteURL, manifest.Id),
						},
						Name: T("button.delete"),
						Type: "action",
					},
				},
			}
		}
	}

	return &model.SlackAttachment{}
}

func (p *Plugin) ListCompletedReminders(userId string, postId string) {

	user, uErr := p.API.GetUser(userId)
	if uErr != nil {
		p.API.LogError(uErr.Error())
	}
	reminders := p.GetReminders(user.Username)

	var output string
	output = ""
	for _, reminder := range reminders {
		if reminder.Completed != p.emptyTime {
			location := p.location(user)
			output = strings.Join([]string{output, "* \"" + reminder.Message + "\" " + fmt.Sprintf("%v", reminder.Completed.In(location).Format(time.RFC3339))}, "\n")
		}
	}

	if post, pErr := p.API.GetPost(postId); pErr != nil {
		p.API.LogError(pErr.Error())
	} else {
		post.Message = output
		post.Props = model.StringInterface{}
		p.API.UpdatePost(post)
	}
}

func (p *Plugin) DeleteCompletedReminders(userId string) {

	user, uErr := p.API.GetUser(userId)
	if uErr != nil {
		p.API.LogError(uErr.Error())
	}
	reminders := p.GetReminders(user.Username)
	for _, reminder := range reminders {
		if reminder.Completed != p.emptyTime {
			p.DeleteReminder(userId, reminder)
		}
	}

}

/*
func (p *Plugin) ListReminders(user *model.User, channelId string) string {

	T, _ := p.translation(user)
	// location := p.location(user)

	var upcomingOccurrences []Occurrence
	var recurringOccurrences []Occurrence
	var pastOccurrences []Occurrence
	var channelOccurrences []Occurrence

	reminders := p.GetReminders(user.Username)

	output := ""

	for _, reminder := range reminders {
		occurrences := reminder.Occurrences
		if len(occurrences) > 0 {
			for _, occurrence := range occurrences {
				t := occurrence.Occurrence
				s := occurrence.Snoozed

				if !strings.HasPrefix(reminder.Target, "~") &&
					reminder.Completed == p.emptyTime &&
					(occurrence.Repeat == "" && t.After(time.Now())) ||
					(s != p.emptyTime && s.After(time.Now())) {
					upcomingOccurrences = append(upcomingOccurrences, occurrence)
				}

				if !strings.HasPrefix(reminder.Target, "~") &&
					occurrence.Repeat != "" && (t.After(time.Now()) ||
					(s != p.emptyTime && s.After(time.Now()))) {
					recurringOccurrences = append(recurringOccurrences, occurrence)
				}

				if !strings.HasPrefix(reminder.Target, "~") &&
					reminder.Completed == p.emptyTime &&
					t.Before(time.Now()) &&
					s == p.emptyTime {
					pastOccurrences = append(pastOccurrences, occurrence)
				}

				if strings.HasPrefix(reminder.Target, "~") &&
					reminder.Completed == p.emptyTime &&
					t.After(time.Now()) {
					channelOccurrences = append(channelOccurrences, occurrence)
				}

			}
		}
	}

	// channel, cErr := a.GetChannel(channelId)
	// if cErr == nil {
	// 	schan := a.Srv.Store.Remind().GetByChannel("~" + channel.Name)
	// 	result := <-schan

	// 	if result.Err != nil {
	// 		mlog.Error(result.Err.Error())
	// 	} else {
	// 		inChannel := result.Data.(model.ChannelReminders)

	// 		if len(inChannel.Occurrences) > 0 {
	// 			output = strings.Join([]string{
	// 				output,
	// 				T("list_inchannel"),
	// 				a.listReminderGroup(userId, &inChannel.Occurrences, &reminders, "inchannel"),
	// 				"\n",
	// 			}, "\n")
	// 		}
	// 	}
	// }

	if len(upcomingOccurrences) > 0 {
		output = strings.Join([]string{
			output,
			T("list_upcoming"),
			p.listReminderGroup(user, &upcomingOccurrences, &reminders, "upcoming"),
			"\n",
		}, "\n")
	}

	if len(recurringOccurrences) > 0 {
		output = strings.Join([]string{
			output,
			T("list_recurring"),
			p.listReminderGroup(user, &recurringOccurrences, &reminders, "recurring"),
			"\n",
		}, "\n")
	}

	if len(pastOccurrences) > 0 {
		output = strings.Join([]string{
			output,
			T("list_past_and_incomplete"),
			p.listReminderGroup(user, &pastOccurrences, &reminders, "past"),
			"\n",
		}, "\n")
	}

	if len(channelOccurrences) > 0 {
		output = strings.Join([]string{
			output,
			T("list_channel"),
			p.listReminderGroup(user, &channelOccurrences, &reminders, "channel"),
			"\n",
		}, "\n")
	}

	return output + T("list_footer")
}

func (p *Plugin) listReminderGroup(user *model.User, occurrences *[]Occurrence, reminders *[]Reminder, gType string) string {

	location := p.location(user)
	T, _ := p.translation(user)

	output := ""

	for _, occurrence := range *occurrences {

		reminder := p.findReminder(*reminders, occurrence)
		t := occurrence.Occurrence
		s := occurrence.Snoozed

		formattedOccurrence := ""
		formattedOccurrence = p.formatWhen(user.Username, reminder.When, t.In(location).Format(time.RFC3339), false)

		formattedSnooze := ""
		if s != p.emptyTime {
			formattedSnooze = p.formatWhen(user.Username, reminder.When, s.In(location).Format(time.RFC3339), true)
		}

		var messageParameters = map[string]interface{}{
			"Message":    reminder.Message,
			"Occurrence": formattedOccurrence,
			"Snoozed":    formattedSnooze,
		}
		if !t.Equal(p.emptyTime) {
			switch gType {
			case "upcoming":
				if formattedSnooze == "" {
					output = strings.Join([]string{output, T("list.element.upcoming", messageParameters)}, "\n")
				} else {
					output = strings.Join([]string{output, T("list.element.upcoming.snoozed", messageParameters)}, "\n")
				}
			case "recurring":
				if formattedSnooze == "" {
					output = strings.Join([]string{output, T("list.element.recurring", messageParameters)}, "\n")
				} else {
					output = strings.Join([]string{output, T("list.element.recurring.snoozed", messageParameters)}, "\n")
				}
			case "past":
				output = strings.Join([]string{output, T("list.element.past", messageParameters)}, "\n")
			case "channel":
				output = strings.Join([]string{output, T("list.element.channel", messageParameters)}, "\n")
			case "inchannel":
				output = strings.Join([]string{output, T("list.element.inchannel", messageParameters)}, "\n")
			}
		}
	}
	return output
}
*/
