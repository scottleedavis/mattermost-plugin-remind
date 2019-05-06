package main

import (
	"testing"
	//"encoding/json"
	//"time"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	//"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestScheduleReminders(t *testing.T) {

	//user := &model.User{
	//	Email:    "-@-.-",
	//	Nickname: "TestUser",
	//	Password: model.NewId(),
	//	Username: "testuser",
	//	Roles:    model.SYSTEM_USER_ROLE_ID,
	//	Locale:   "en",
	//}
	//
	//occurrences := []Occurrence{
	//	{
	//		Id: model.NewId(),
	//		ReminderId: model.NewId(),
	//		Occurrence:  time.Now(),
	//	},
	//}
	//
	//stringOccurrences, _ := json.Marshal(occurrences)
	//channel := &model.Channel{
	//	Id: model.NewId(),
	//}
	//
	//setupAPI := func() *plugintest.API {
	//	api := &plugintest.API{}
	//	api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
	//	api.On("LogInfo", mock.Anything).Maybe()
	//	api.On("GetUserByUsername", mock.AnythingOfType("string")).Return(user, nil)
	//	api.On("KVGet", mock.AnythingOfType("string")).Return(stringOccurrences, nil)
	//	api.On("KVSet", mock.AnythingOfType("string"), stringOccurrences).Maybe()
	//	return api
	//}
	//
	//t.Run("if scheduled reminder" , func(t *testing.T) {
	//
	//	api := setupAPI()
	//	defer api.AssertExpectations(t)
	//
	//	p := &Plugin{}
	//	p.API = api
	//
	//	request := &ReminderRequest{
	//		Username:   user.Username,
	//		Payload:    "me foo in 1 seconds",
	//	}
	//	post, err := p.ScheduleReminder(request, channel.Id)
	//	assert.Nil(t, err)
	//	assert.Equal(t, len(post.Attachments()), 2)
	//
	//})

}


func TestInteractiveSchedule(t *testing.T) {

	user := &model.User{
		Email:    "-@-.-",
		Nickname: "TestUser",
		Password: model.NewId(),
		Username: "testuser",
		Roles:    model.SYSTEM_USER_ROLE_ID,
		Locale:   "en",
	}

	//occurrences := []Occurrence{
	//	{
	//		Id: model.NewId(),
	//		ReminderId: model.NewId(),
	//		Occurrence:  time.Now(),
	//	},
	//}
	//
	//stringOccurrences, _ := json.Marshal(occurrences)
	//channel := &model.Channel{
	//	Id: model.NewId(),
	//}

	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		//api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		//api.On("LogInfo", mock.Anything).Maybe()
		api.On("OpenInteractiveDialog", mock.AnythingOfType("model.OpenDialogRequest")).Return(nil)
		//api.On("KVGet", mock.AnythingOfType("string")).Return(stringOccurrences, nil)
		//api.On("KVSet", mock.AnythingOfType("string"), stringOccurrences).Maybe()
		return api
	}

	t.Run("if interactive schedule" , func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		p.InteractiveSchedule(model.NewId(), user)

	})

}
