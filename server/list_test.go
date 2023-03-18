package main

import (
	"testing"
	"time"

	"encoding/json"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListReminders(t *testing.T) {
	p := &Plugin{}

	user := &model.User{
		Email:    "-@-.-",
		Nickname: "TestUser",
		Password: model.NewId(),
		Username: "testuser",
		Roles:    model.SystemUserRoleId,
		Locale:   "en",
	}

	originChannel := &model.Channel{
		Id: model.NewId(),
	}

	publicChannel := &model.Channel{
		Id:   model.NewId(),
		Name: "public-channel",
	}

	T, _ := p.translation(user)

	pastOccurrenceDate := time.Now().Add(-1 * time.Minute)
	futureOccurrenceDate := time.Now().Add(1 * time.Minute)

	pastReminder := Reminder{
		Id:       model.NewId(),
		Username: user.Username,
		Message:  "This reminder was triggered a single time",
		When:     "in 1 minute",
		Occurrences: []Occurrence{
			{
				Id:         model.NewId(),
				ReminderId: model.NewId(),
				Occurrence: pastOccurrenceDate,
			},
		},
		Completed: p.emptyTime,
	}

	upcomingReminder := Reminder{
		Id:       model.NewId(),
		Username: user.Username,
		Message:  "This reminder triggers a single time",
		When:     "in 1 minute",
		Occurrences: []Occurrence{
			{
				Id:         model.NewId(),
				ReminderId: model.NewId(),
				Occurrence: futureOccurrenceDate,
			},
		},
		Completed: p.emptyTime,
	}

	channelReminder := Reminder{
		Id:       model.NewId(),
		Username: user.Username,
		Target:   "~" + publicChannel.Name,
		Message:  "This reminder posts in a channel a single time",
		When:     "in 1 minute",
		Occurrences: []Occurrence{
			{
				Id:         model.NewId(),
				ReminderId: model.NewId(),
				Occurrence: futureOccurrenceDate,
			},
		},
		Completed: p.emptyTime,
	}

	recurringReminder := Reminder{
		Id:       model.NewId(),
		Username: user.Username,
		Message:  "This reminder triggers several times",
		When:     "every Monday",
		Occurrences: []Occurrence{
			{
				Id:         model.NewId(),
				ReminderId: model.NewId(),
				Repeat:     "every Monday at 9:00AM",
				Occurrence: futureOccurrenceDate,
			},
		},
		Completed: p.emptyTime,
	}

	setupAPI := func(reminders []Reminder) *plugintest.API {
		serializedReminders, _ := json.Marshal(reminders)
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
		api.On("KVGet", "testuser").Return(serializedReminders, nil)
		return api
	}

	t.Run("the list categorizes reminders by type", func(t *testing.T) {
		reminders := []Reminder{pastReminder, upcomingReminder, channelReminder, recurringReminder}
		p.API = setupAPI(reminders)

		post := p.ListReminders(user, originChannel.Id)

		attachments := post.Attachments()
		assert.NotNil(t, post, "A post must be returned")
		assert.Equal(t, post.ChannelId, originChannel.Id, "The list must be posted to the channel it was requested")
		assert.Equal(t, len(reminders)+1, len(attachments), "The list must have one attachment per active reminder, plus an attachment for control")
		assert.Contains(t, attachments[0].Text, T("list.upcoming"), "The first displayed reminders must be upcoming reminders")
		assert.Contains(t, attachments[1].Text, T("list.recurring"), "The next displayed reminders must be recurring reminders")
		assert.Contains(t, attachments[2].Text, T("list.past.and.incomplete"), "The next displayed reminders must be past and incomplete reminders")
		assert.Contains(t, attachments[3].Text, T("list.channel"), "The next displayed reminders must be channel reminders")
		assert.Contains(t, attachments[len(attachments)-1].Text, T("reminders.page.numbers"), "The last attachment must be list pagination and controls")
	})
}

func TestUpdateListReminders(t *testing.T) {

	user := &model.User{
		Email:    "-@-.-",
		Nickname: "TestUser",
		Password: model.NewId(),
		Username: "testuser",
		Roles:    model.SystemUserRoleId,
		Locale:   "en",
	}

	occurrences := []Occurrence{
		{
			Id:         model.NewId(),
			ReminderId: "ididididid",
			Occurrence: time.Now(),
		},
	}

	reminders := []Reminder{
		{
			Id:          model.NewId(),
			Username:    user.Username,
			Message:     "Message",
			When:        "in 1 second",
			Occurrences: occurrences,
			Completed:   time.Time{}.AddDate(1, 1, 1),
		},
	}

	stringReminders, _ := json.Marshal(reminders)

	channel := &model.Channel{
		Id: model.NewId(),
	}

	post := &model.Post{
		ChannelId:     channel.Id,
		PendingPostId: model.NewId(),
		UserId:        user.Id,
		Props:         model.StringInterface{},
	}

	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
		api.On("GetUser", mock.AnythingOfType("string")).Return(user, nil)
		api.On("UpdateEphemeralPost", mock.Anything, mock.Anything).Return(post)
		api.On("KVGet", mock.Anything).Return(stringReminders, nil)
		return api
	}

	t.Run("if update list happens", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		p.UpdateListReminders(user.Id, model.NewId(), channel.Id, 0)

	})

}

func TestListCompletedReminders(t *testing.T) {

	user := &model.User{
		Email:    "-@-.-",
		Nickname: "TestUser",
		Password: model.NewId(),
		Username: "testuser",
		Roles:    model.SystemUserRoleId,
		Locale:   "en",
	}

	occurrences := []Occurrence{
		{
			Id:         model.NewId(),
			ReminderId: "ididididid",
			Occurrence: time.Now(),
		},
	}

	reminders := []Reminder{
		{
			Id:          model.NewId(),
			Username:    user.Username,
			Message:     "Message",
			When:        "in 1 second",
			Occurrences: occurrences,
			Completed:   time.Time{}.AddDate(1, 1, 1),
		},
	}

	stringReminders, _ := json.Marshal(reminders)

	channel := &model.Channel{
		Id: model.NewId(),
	}

	post := &model.Post{
		ChannelId:     channel.Id,
		PendingPostId: model.NewId(),
		UserId:        user.Id,
		Props:         model.StringInterface{},
	}

	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
		api.On("GetUser", mock.AnythingOfType("string")).Return(user, nil)
		api.On("UpdateEphemeralPost", mock.AnythingOfType("string"), mock.Anything).Return(post)
		api.On("KVGet", mock.Anything).Return(stringReminders, nil)
		return api
	}

	t.Run("if list completed happens", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		p.ListCompletedReminders(user.Id, post.Id, channel.Id)

	})

}

func TestDeleteCompletedReminders(t *testing.T) {

	user := &model.User{
		Email:    "-@-.-",
		Nickname: "TestUser",
		Password: model.NewId(),
		Username: "testuser",
		Roles:    model.SystemUserRoleId,
		Locale:   "en",
	}

	occurrences := []Occurrence{
		{
			Id:         model.NewId(),
			ReminderId: model.NewId(),
			Occurrence: time.Now(),
		},
	}

	reminders := []Reminder{
		{
			Id:          model.NewId(),
			Username:    user.Username,
			Message:     "Message",
			When:        "in 1 second",
			Occurrences: occurrences,
			Completed:   time.Time{}.AddDate(1, 1, 1),
		},
	}

	stringReminders, _ := json.Marshal(reminders)
	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
		api.On("GetUser", mock.AnythingOfType("string")).Return(user, nil)
		api.On("KVGet", mock.Anything).Return(stringReminders, nil)
		api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
		return api
	}

	t.Run("if delete completed happens", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		p.DeleteCompletedReminders(user.Id)

	})

}
