package main

import (
	"encoding/json"
	"testing"
	"time"
	"fmt"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"

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
	testTime :=  time.Now().UTC().Round(time.Second)

	occurrences := []Occurrence{
		{
			Id: model.NewId(),
			ReminderId: model.NewId(),
			Occurrence:  testTime,
		},
	}

	reminders := []Reminder{
		{
			Id:        model.NewId(),
			TeamId:    model.NewId(),
			Username:  user.Username,
			Message:   "Hello",
			Target:    "me",
			When:      "in one minute",
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
		return api
	}

	t.Run("if triggers reminder for user", func(t *testing.T) {
		api := setupAPI()
		stringReminders, _ := json.Marshal(reminders)
		api.On("KVGet", user.Username).Return(stringReminders, nil)
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
	testTime :=  time.Now().UTC().Round(time.Second)

	occurrences := []Occurrence{
		{
			Id: model.NewId(),
			ReminderId: model.NewId(),
			Occurrence:  testTime,
		},
	}

	reminders := []Reminder{
		{
			Id:        model.NewId(),
			TeamId:    model.NewId(),
			Username:  user.Username,
			Message:   "Hello",
			Target:    "me",
			When:      "in one minute",
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
	testTime :=  time.Now().UTC().Round(time.Second)

	occurrences := []Occurrence{
		{
			Id: model.NewId(),
			ReminderId: model.NewId(),
			Occurrence:  testTime,
		},
	}

	reminders := []Reminder{
		{
			Id:        model.NewId(),
			TeamId:    model.NewId(),
			Username:  user.Username,
			Message:   "Hello",
			Target:    "me",
			When:      "in one minute",
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

	//user := &model.User{
	//	Email:    "-@-.-",
	//	Nickname: "TestUser",
	//	Password: model.NewId(),
	//	Username: "testuser",
	//	Roles:    model.SYSTEM_USER_ROLE_ID,
	//	Locale:   "en",
	//}
	//testTime :=  time.Now().UTC().Round(time.Second)
	//
	//occurrences := []Occurrence{
	//	{
	//		Id: model.NewId(),
	//		ReminderId: model.NewId(),
	//		Occurrence:  testTime,
	//	},
	//}
	//
	//reminders := []Reminder{
	//	{
	//		Id:        model.NewId(),
	//		TeamId:    model.NewId(),
	//		Username:  user.Username,
	//		Message:   "Hello",
	//		Target:    "me",
	//		When:      "in one minute",
	//		Occurrences: occurrences,
	//	},
	//}
	//
	//stringReminders, _ := json.Marshal(reminders)
	//setupAPI := func() *plugintest.API {
	//	api := &plugintest.API{}
	//	api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
	//	api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
	//	api.On("LogInfo", mock.Anything).Maybe()
	//	api.On("KVGet", user.Username).Return(stringReminders, nil)
	//	api.On("KVSet", user.Username).Maybe()
	//	api.On("GetUser", mock.AnythingOfType("string")).Return(user, nil)
	//	//api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
	//	return api
	//}
	//
	//t.Run("if updates reminders", func(t *testing.T) {
	//	api := setupAPI()
	//	defer api.AssertExpectations(t)
	//
	//	p := &Plugin{}
	//	p.API = api
	//
	//	assert.Nil(t, p.UpdateReminder(user.Id, reminders[0]))
	//
	//})
	//

}

func TestUpsertReminder(t *testing.T) {

	//user := &model.User{
	//	Email:    "-@-.-",
	//	Nickname: "TestUser",
	//	Password: model.NewId(),
	//	Username: "testuser",
	//	Roles:    model.SYSTEM_USER_ROLE_ID,
	//	Locale:   "en",
	//}
	//testTime :=  time.Now().UTC().Round(time.Second)
	//
	//occurrences := []Occurrence{
	//	{
	//		Id: model.NewId(),
	//		ReminderId: model.NewId(),
	//		Occurrence:  testTime,
	//	},
	//}
	//
	//reminders := []Reminder{
	//	{
	//		Id:        model.NewId(),
	//		TeamId:    model.NewId(),
	//		Username:  user.Username,
	//		Message:   "Hello",
	//		Target:    "me",
	//		When:      "in one minute",
	//		Occurrences: occurrences,
	//	},
	//}
	//
	//
	//request := &ReminderRequest{
	//	TeamId:   model.NewId(),
	//	Username: user.Username,
	//	Payload:  "Hello in one minute",
	//	Reminder: Reminder{
	//		Id:        model.NewId(),
	//		TeamId:    model.NewId(),
	//		Username:  user.Username,
	//		Message:   "Hello",
	//		Target:    "me",
	//		When:      "in one minute",
	//	},
	//}
	//
	//stringReminders, _ := json.Marshal(reminders)
	//setupAPI := func() *plugintest.API {
	//	api := &plugintest.API{}
	//	api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
	//	api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
	//	api.On("LogInfo", mock.Anything).Maybe()
	//	api.On("KVGet", user.Username).Return(stringReminders, nil)
	//	api.On("KVSet", user.Username).Maybe()
	//	//api.On("GetUser", mock.AnythingOfType("string")).Return(user, nil)
	//	api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
	//	return api
	//}
	//
	//t.Run("if updates reminders", func(t *testing.T) {
	//	api := setupAPI()
	//	defer api.AssertExpectations(t)
	//
	//	p := &Plugin{}
	//	p.API = api
	//
	//	assert.Nil(t, p.UpsertReminder(request))
	//
	//})

}

func TestDeleteReminder(t *testing.T) {


}

func TestDeleteReminders(t *testing.T) {

}