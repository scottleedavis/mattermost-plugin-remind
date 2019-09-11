package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTriggerReminders(t *testing.T) {

	user := &model.User{
		Email:    "-@-.-",
		Nickname: "TestUser",
		Password: model.NewId(),
		Username: "testuser",
		Roles:    model.SYSTEM_USER_ROLE_ID,
		Locale:   "en",
	}

	channel := &model.Channel{
		Id: model.NewId(),
	}

	post := &model.Post{
		Id:        model.NewId(),
		ChannelId: channel.Id,
	}

	testTime := time.Now().UTC().Round(time.Second)
	hostname, _ := os.Hostname()
	reminder1 := model.NewId()
	reminder2 := model.NewId()
	reminder3 := model.NewId()

	occurrences := []Occurrence{
		{
			Hostname:   hostname,
			Username:   user.Username,
			Id:         model.NewId(),
			ReminderId: reminder1,
			Occurrence: testTime,
		},
		{
			Hostname:   hostname,
			Username:   user.Username,
			Id:         model.NewId(),
			ReminderId: reminder2,
			Occurrence: testTime,
		},
		{
			Hostname:   hostname,
			Username:   user.Username,
			Id:         model.NewId(),
			ReminderId: reminder3,
			Occurrence: testTime,
		},
	}

	reminders := []Reminder{
		{
			Id:          reminder1,
			TeamId:      model.NewId(),
			Username:    user.Username,
			Message:     "Hello",
			Target:      "me",
			When:        "in one minute",
			Occurrences: []Occurrence{occurrences[0]},
		},
		{
			Id:          reminder2,
			TeamId:      model.NewId(),
			Username:    user.Username,
			Message:     "Hello 2",
			Target:      "me",
			When:        "in one minute",
			Occurrences: []Occurrence{occurrences[1]},
		},
		{
			Id:          reminder3,
			PostId:      model.NewId(),
			TeamId:      model.NewId(),
			Username:    user.Username,
			Message:     "Hello 2",
			Target:      "me",
			When:        "in one minute",
			Occurrences: []Occurrence{occurrences[2]},
		},
	}

	stringOccurrences, _ := json.Marshal(occurrences)

	team := &model.Team{
		Id: model.NewId(),
	}

	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("KVGet", string(fmt.Sprintf("%v", testTime))).Return(stringOccurrences, nil)
		api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
		api.On("GetDirectChannel", mock.Anything, mock.Anything).Return(channel, nil)
		api.On("CreatePost", mock.Anything).Return(post, nil)
		api.On("GetPost", mock.Anything).Return(post, nil)
		api.On("GetTeam", mock.Anything).Return(team, nil)
		api.On("GetUser", mock.Anything).Return(user, nil)

		return api
	}

	t.Run("it doesn't trigger on different hostname", func(t *testing.T) {
		occurrences := []Occurrence{
			{
				Hostname:   model.NewId(),
				Id:         model.NewId(),
				ReminderId: reminder1,
				Occurrence: testTime,
			},
		}

		stringOccurrences, _ := json.Marshal(occurrences)
		api := &plugintest.API{}
		api.On("KVGet", string(fmt.Sprintf("%v", testTime))).Return(stringOccurrences, nil)
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		p.TriggerReminders()

	})

	t.Run("if triggers reminder for me", func(t *testing.T) {
		api := setupAPI()
		stringReminders, _ := json.Marshal(reminders)
		api.On("KVGet", user.Username).Return(stringReminders, nil)
		api.On("GetDirectChannel", mock.Anything, mock.Anything).Return(channel, nil)
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		p.TriggerReminders()

	})

	t.Run("if triggers reminder for user", func(t *testing.T) {
		api := setupAPI()
		reminders[0].Target = "@testuser"
		stringReminders, _ := json.Marshal(reminders)
		api.On("KVGet", user.Username).Return(stringReminders, nil)
		api.On("GetDirectChannel", mock.Anything, mock.Anything).Return(channel, nil)
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		p.TriggerReminders()

	})

	t.Run("if triggers reminder for channel", func(t *testing.T) {
		api := setupAPI()
		reminders[0].Target = "~off-topic"
		stringReminders, _ := json.Marshal(reminders)
		api.On("KVGet", user.Username).Return(stringReminders, nil)
		api.On("GetChannelByName", mock.Anything, mock.Anything, mock.Anything).Return(channel, nil)
		api.On("CreatePost", mock.Anything).Return(nil, nil)
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		p.TriggerReminders()

	})

	t.Run("if triggers user reminder recurring", func(t *testing.T) {
		occurrences := []Occurrence{
			{
				Hostname:   hostname,
				Id:         model.NewId(),
				ReminderId: reminder1,
				Occurrence: testTime,
				Repeat:     "every tuesday at 3pm",
			},
		}
		reminders := []Reminder{
			{
				Id:          reminder1,
				TeamId:      model.NewId(),
				Username:    user.Username,
				Message:     "Hello",
				Target:      "me",
				When:        "every tuesday at 3pm",
				Occurrences: occurrences,
			},
		}

		stringReminders, _ := json.Marshal(reminders)
		stringOccurrences, _ := json.Marshal(occurrences)
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
		api.On("CreatePost", mock.Anything).Return(post, nil)
		api.On("GetDirectChannel", mock.Anything, mock.Anything).Return(channel, nil)

		api.On("KVGet", string(fmt.Sprintf("%v", testTime))).Return(stringOccurrences, nil)
		api.On("KVGet", user.Username).Return(stringReminders, nil)
		api.On("KVGet", mock.Anything).Return(stringOccurrences, nil)
		api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
		api.On("GetUser", mock.Anything).Return(user, nil)
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		p.TriggerReminders()

	})

	t.Run("if triggers channel reminder recurring", func(t *testing.T) {
		occurrences := []Occurrence{
			{
				Hostname:   hostname,
				Id:         model.NewId(),
				ReminderId: reminder1,
				Occurrence: testTime,
				Repeat:     "every tuesday at 3pm",
			},
		}
		reminders := []Reminder{
			{
				Id:          reminder1,
				TeamId:      model.NewId(),
				Username:    user.Username,
				Message:     "Hello",
				Target:      "~town-square",
				When:        "every tuesday at 3pm",
				Occurrences: occurrences,
			},
		}

		stringReminders, _ := json.Marshal(reminders)
		stringOccurrences, _ := json.Marshal(occurrences)
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
		api.On("CreatePost", mock.Anything).Return(post, nil)
		api.On("GetChannelByName", mock.Anything, mock.Anything, mock.Anything).Return(channel, nil)
		api.On("KVGet", user.Username).Return(stringReminders, nil)
		api.On("KVGet", mock.Anything).Return(stringOccurrences, nil)
		api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
		api.On("GetUser", mock.Anything).Return(user, nil)
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		p.TriggerReminders()

	})

}

func TestGetReminder(t *testing.T) {

	user := &model.User{
		Email:    "-@-.-",
		Nickname: "TestUser",
		Password: model.NewId(),
		Username: "testuser",
		Roles:    model.SYSTEM_USER_ROLE_ID,
		Locale:   "en",
	}
	testTime := time.Now().UTC().Round(time.Second)

	occurrences := []Occurrence{
		{
			Id:         model.NewId(),
			ReminderId: model.NewId(),
			Occurrence: testTime,
		},
	}

	reminders := []Reminder{
		{
			Id:          model.NewId(),
			TeamId:      model.NewId(),
			Username:    user.Username,
			Message:     "Hello",
			Target:      "me",
			When:        "in one minute",
			Occurrences: occurrences,
		},
	}

	stringReminders, _ := json.Marshal(reminders)
	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("KVGet", user.Username).Return(stringReminders, nil)
		api.On("GetUser", mock.AnythingOfType("string")).Return(user, nil)
		api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
		return api
	}

	t.Run("if gets reminder", func(t *testing.T) {
		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		reminder := p.GetReminder(user.Id, reminders[0].Id)
		assert.Equal(t, reminders[0].Id, reminder.Id)

	})

}

func TestGetReminders(t *testing.T) {

	user := &model.User{
		Email:    "-@-.-",
		Nickname: "TestUser",
		Password: model.NewId(),
		Username: "testuser",
		Roles:    model.SYSTEM_USER_ROLE_ID,
		Locale:   "en",
	}
	testTime := time.Now().UTC().Round(time.Second)

	occurrences := []Occurrence{
		{
			Id:         model.NewId(),
			ReminderId: model.NewId(),
			Occurrence: testTime,
		},
	}

	reminders := []Reminder{
		{
			Id:          model.NewId(),
			TeamId:      model.NewId(),
			Username:    user.Username,
			Message:     "Hello",
			Target:      "me",
			When:        "in one minute",
			Occurrences: occurrences,
		},
	}

	stringReminders, _ := json.Marshal(reminders)
	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("KVGet", user.Username).Return(stringReminders, nil)
		api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
		return api
	}

	t.Run("if gets reminders", func(t *testing.T) {
		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		reminders := p.GetReminders(user.Username)
		assert.True(t, len(reminders) == 1)

	})

}

func TestUpdateReminder(t *testing.T) {

	user := &model.User{
		Email:    "-@-.-",
		Nickname: "TestUser",
		Password: model.NewId(),
		Username: "testuser",
		Roles:    model.SYSTEM_USER_ROLE_ID,
		Locale:   "en",
	}
	testTime := time.Now().UTC().Round(time.Second)

	occurrences := []Occurrence{
		{
			Id:         model.NewId(),
			ReminderId: model.NewId(),
			Occurrence: testTime,
		},
	}

	reminders := []Reminder{
		{
			Id:          model.NewId(),
			TeamId:      model.NewId(),
			Username:    user.Username,
			Message:     "Hello",
			Target:      "me",
			When:        "in one minute",
			Occurrences: occurrences,
		},
	}

	stringReminders, _ := json.Marshal(reminders)
	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("KVGet", user.Username).Return(stringReminders, nil)
		api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
		api.On("GetUser", mock.AnythingOfType("string")).Return(user, nil)
		return api
	}

	t.Run("if updates reminders", func(t *testing.T) {
		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		assert.Nil(t, p.UpdateReminder(user.Id, reminders[0]))

	})

}

func TestUpsertReminder(t *testing.T) {

	user := &model.User{
		Email:    "-@-.-",
		Nickname: "TestUser",
		Password: model.NewId(),
		Username: "testuser",
		Roles:    model.SYSTEM_USER_ROLE_ID,
		Locale:   "en",
	}
	testTime := time.Now().UTC().Round(time.Second)

	occurrences := []Occurrence{
		{
			Id:         model.NewId(),
			ReminderId: model.NewId(),
			Occurrence: testTime,
		},
	}

	reminders := []Reminder{
		{
			Id:          model.NewId(),
			TeamId:      model.NewId(),
			Username:    user.Username,
			Message:     "Hello",
			Target:      "me",
			When:        "in one minute",
			Occurrences: occurrences,
		},
	}

	request := &ReminderRequest{
		TeamId:   model.NewId(),
		Username: user.Username,
		Payload:  "Hello in one minute",
		Reminder: Reminder{
			Id:          model.NewId(),
			TeamId:      model.NewId(),
			Username:    user.Username,
			Occurrences: occurrences,
			Message:     "Hello",
			Target:      "me",
			When:        "in one minute",
		},
	}

	stringReminders, _ := json.Marshal(reminders)
	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("KVGet", user.Username).Return(stringReminders, nil)
		api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
		api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
		return api
	}

	t.Run("if updates reminders", func(t *testing.T) {
		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		assert.Nil(t, p.UpsertReminder(request))

	})

}

func TestDeleteReminder(t *testing.T) {

	user := &model.User{
		Email:    "-@-.-",
		Nickname: "TestUser",
		Password: model.NewId(),
		Username: "testuser",
		Roles:    model.SYSTEM_USER_ROLE_ID,
		Locale:   "en",
	}
	testTime := time.Now().UTC().Round(time.Second)

	occurrences := []Occurrence{
		{
			Id:         model.NewId(),
			ReminderId: model.NewId(),
			Occurrence: testTime,
		},
	}

	reminders := []Reminder{
		{
			Id:          model.NewId(),
			TeamId:      model.NewId(),
			Username:    user.Username,
			Message:     "Hello",
			Target:      "me",
			When:        "in one minute",
			Occurrences: occurrences,
		},
	}

	stringReminders, _ := json.Marshal(reminders)
	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("KVGet", user.Username).Return(stringReminders, nil)
		api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
		api.On("GetUser", mock.AnythingOfType("string")).Return(user, nil)
		return api
	}

	t.Run("if deletes reminder", func(t *testing.T) {
		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		assert.Nil(t, p.DeleteReminder(user.Id, reminders[0]))

	})

}

func TestDeleteReminders(t *testing.T) {
	user := &model.User{
		Email:    "-@-.-",
		Nickname: "TestUser",
		Password: model.NewId(),
		Username: "testuser",
		Roles:    model.SYSTEM_USER_ROLE_ID,
		Locale:   "en",
	}
	testTime := time.Now().UTC().Round(time.Second)

	occurrences := []Occurrence{
		{
			Id:         model.NewId(),
			ReminderId: model.NewId(),
			Occurrence: testTime,
		},
	}

	reminders := []Reminder{
		{
			Id:          model.NewId(),
			TeamId:      model.NewId(),
			Username:    user.Username,
			Message:     "Hello",
			Target:      "me",
			When:        "in one minute",
			Occurrences: occurrences,
		},
	}

	stringReminders, _ := json.Marshal(reminders)
	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("KVGet", user.Username).Return(stringReminders, nil)
		api.On("KVDelete", user.Username).Return(nil)
		return api
	}

	t.Run("if deletes reminders", func(t *testing.T) {
		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		assert.Equal(t, p.DeleteReminders(user), "clear.response")

	})

}
