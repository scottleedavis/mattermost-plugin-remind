package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/mattermost/mattermost-server/v5/model"
	"github.com/mattermost/mattermost-server/v5/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTriggerReminders(t *testing.T) {
  testTime := time.Now().UTC().Round(time.Second)
  serializedTestTime := []byte(testTime.Format(time.RFC3339))

	t.Run("it triggers reminders scheduled for the current time", func(t *testing.T) {
		oneSecondAgo, _ := time.ParseDuration("-1s")
		lastTickAt := testTime.Add(oneSecondAgo)
	  serializedLastTickAt := []byte(lastTickAt.Format(time.RFC3339))

		api := &plugintest.API{}
		api.On("KVGet", string("LastTickAt")).Return(serializedLastTickAt, nil)
		api.On("KVSet", string("LastTickAt"), serializedTestTime).Return(nil)
		api.On("LogDebug", "Trigger reminders for " + fmt.Sprintf("%v", testTime))
		api.On("KVGet", string(fmt.Sprintf("%v", testTime))).Return(nil, nil)
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		p.TriggerReminders()
	})

	t.Run("when ticks have been missed, it triggers reminders for the missed ticks as well", func(t *testing.T) {
		oneSecondsAgo, _ := time.ParseDuration("-1s")
		twoSecondsAgo, _ := time.ParseDuration("-2s")
		threeSecondsAgo, _ := time.ParseDuration("-3s")
	  lastTickAt := testTime.Add(threeSecondsAgo)
	  serializedLastTickAt := []byte(lastTickAt.Format(time.RFC3339))

		api := &plugintest.API{}
		api.On("KVGet", string("LastTickAt")).Return(serializedLastTickAt, nil)
		api.On("KVSet", string("LastTickAt"), serializedTestTime).Return(nil)
		api.On("LogDebug", "Catching up on 2 reminder tick(s)...")
		api.On("LogDebug", "Trigger reminders for " + fmt.Sprintf("%v", testTime.Add(twoSecondsAgo)))
		api.On("KVGet", string(fmt.Sprintf("%v", testTime.Add(twoSecondsAgo)))).Return(nil, nil)
		api.On("LogDebug", "Trigger reminders for " + fmt.Sprintf("%v", testTime.Add(oneSecondsAgo)))
		api.On("KVGet", string(fmt.Sprintf("%v", testTime.Add(oneSecondsAgo)))).Return(nil, nil)
		api.On("LogDebug", "Caught up on missed reminder ticks.")
		api.On("LogDebug", "Trigger reminders for " + fmt.Sprintf("%v", testTime))
		api.On("KVGet", string(fmt.Sprintf("%v", testTime))).Return(nil, nil)
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		p.TriggerReminders()
	})
}

func TestTriggerRemindersForTick(t *testing.T) {

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
		Id: model.NewId(),
	}

	testTime := time.Now().UTC().Round(time.Second)
	hostname, _ := os.Hostname()
	reminderId := model.NewId()

	occurrences := []Occurrence{
		{
			Hostname:   hostname,
			Id:         model.NewId(),
			ReminderId: reminderId,
			Occurrence: testTime,
		},
	}

	reminders := []Reminder{
		{
			Id:          reminderId,
			TeamId:      model.NewId(),
			Username:    user.Username,
			Message:     "Hello",
			Target:      "me",
			When:        "in one minute",
			Occurrences: occurrences,
		},
	}

	stringOccurrences, _ := json.Marshal(occurrences)

	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("KVGet", string(fmt.Sprintf("%v", testTime))).Return(stringOccurrences, nil)
		api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
		api.On("CreatePost", mock.Anything).Return(post, nil)
		return api
	}

	t.Run("it doesn't trigger on different hostname", func(t *testing.T) {
		occurrences := []Occurrence{
			{
				Hostname:   model.NewId(),
				Id:         model.NewId(),
				ReminderId: reminderId,
				Occurrence: testTime,
			},
		}

		stringOccurrences, _ := json.Marshal(occurrences)
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("KVGet", string(fmt.Sprintf("%v", testTime))).Return(stringOccurrences, nil)
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		p.TriggerRemindersForTick(testTime)

	})

	t.Run("if triggers reminder for me", func(t *testing.T) {
		api := setupAPI()
		stringReminders, _ := json.Marshal(reminders)
		api.On("KVGet", user.Username).Return(stringReminders, nil)
		api.On("GetDirectChannel", mock.Anything, mock.Anything).Return(channel, nil)
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		p.TriggerRemindersForTick(testTime)

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

		p.TriggerRemindersForTick(testTime)

	})

	t.Run("if triggers reminder for channel", func(t *testing.T) {
		api := setupAPI()
		reminders[0].Target = "~off-topic"
		stringReminders, _ := json.Marshal(reminders)
		api.On("KVGet", user.Username).Return(stringReminders, nil)
		api.On("GetChannelByName", mock.Anything, mock.Anything, mock.Anything).Return(channel, nil)
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		p.TriggerRemindersForTick(testTime)

	})

	t.Run("if triggers user reminder recurring", func(t *testing.T) {
		occurrences := []Occurrence{
			{
				Hostname:   hostname,
				Id:         model.NewId(),
				ReminderId: reminderId,
				Occurrence: testTime,
				Repeat:     "every tuesday at 3pm",
			},
		}
		reminders := []Reminder{
			{
				Id:          reminderId,
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

		p.TriggerRemindersForTick(testTime)

	})

	t.Run("if triggers channel reminder recurring", func(t *testing.T) {
		occurrences := []Occurrence{
			{
				Hostname:   hostname,
				Id:         model.NewId(),
				ReminderId: reminderId,
				Occurrence: testTime,
				Repeat:     "every tuesday at 3pm",
			},
		}
		reminders := []Reminder{
			{
				Id:          reminderId,
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

		p.TriggerRemindersForTick(testTime)

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
