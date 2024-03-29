package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleDialog(t *testing.T) {
	user := &model.User{
		Id:       model.NewId(),
		Username: model.NewRandomString(10),
	}
	testTime := time.Now().UTC().Round(time.Second).Add(20 * time.Duration(time.Minute))
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
	stringReminders, _ := json.Marshal(reminders)

	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("KVGet", string(fmt.Sprintf("%v", testTime))).Return(stringOccurrences, nil)
		api.On("KVGet", user.Username).Return(stringReminders, nil)
		api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
		api.On("GetUser", mock.Anything).Return(user, nil)
		api.On("GetUserByUsername", mock.Anything).Return(user, nil)
		api.On("SendEphemeralPost", mock.Anything, mock.Anything).Return(nil)

		return api
	}

	t.Run("view dialog", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.router = p.InitAPI()
		p.API = api

		request := &model.SubmitDialogRequest{
			UserId: "userID1",
			Submission: model.StringInterface{
				"message": "hello",
				"target":  "me",
				"time":    "unit.test",
			},
		}

		w := httptest.NewRecorder()
		requestJSON, _ := json.Marshal(request)
		r := httptest.NewRequest("POST", "/dialog", bytes.NewReader(requestJSON))
		p.ServeHTTP(nil, w, r)

		result := w.Result()
		assert.NotNil(t, result)

		bodyBytes, err := io.ReadAll(result.Body)
		assert.Nil(t, err)
		bodyString := string(bodyBytes)
		assert.Equal(t, bodyString, "")

	})
}

func TestHandleViewEphmeral(t *testing.T) {

	user := &model.User{
		Id:       model.NewId(),
		Username: model.NewRandomString(10),
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
		api.On("GetUser", mock.Anything).Return(user, nil)
		api.On("GetUserByUsername", mock.Anything).Return(user, nil)
		api.On("KVGet", user.Username).Return(stringReminders, nil)
		api.On("CreatePost", mock.Anything).Maybe()
		api.On("SendEphemeralPost", mock.Anything, mock.Anything).Return(nil)

		return api
	}

	t.Run("view ephemeral", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.router = p.InitAPI()
		p.API = api

		request := &model.PostActionIntegrationRequest{
			UserId: "userID1",
			PostId: "postID1",
			Context: model.StringInterface{
				"reminder_id":   model.NewId(),
				"occurrence_id": model.NewId(),
			},
		}

		w := httptest.NewRecorder()
		requestJSON, _ := json.Marshal(request)
		r := httptest.NewRequest("POST", "/view/ephemeral", bytes.NewReader(requestJSON))
		p.ServeHTTP(nil, w, r)

		result := w.Result()
		assert.NotNil(t, result)

		bodyBytes, err := io.ReadAll(result.Body)
		assert.Nil(t, err)
		bodyString := string(bodyBytes)
		assert.Equal(t, bodyString, "{\"update\":null,\"ephemeral_text\":\"\",\"skip_slack_parsing\":false}")

	})

}

func TestHandleComplete(t *testing.T) {
	user := &model.User{
		Id:       model.NewId(),
		Username: model.NewRandomString(10),
	}
	post := &model.Post{
		Id:        model.NewId(),
		ChannelId: model.NewId(),
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
	channel := &model.Channel{
		Id: model.NewId(),
	}

	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("GetPost", mock.Anything).Return(post, nil)
		api.On("CreatePost", mock.Anything).Return(post, nil)
		api.On("UpdatePost", mock.Anything).Return(post, nil)
		api.On("GetUser", mock.Anything).Return(user, nil)
		api.On("GetUserByUsername", mock.Anything).Return(user, nil)
		api.On("KVGet", user.Username).Return(stringReminders, nil)
		api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
		api.On("GetDirectChannel", mock.Anything, mock.Anything).Return(channel, nil)

		return api
	}

	t.Run("complete", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.router = p.InitAPI()
		p.API = api

		request := &model.PostActionIntegrationRequest{
			UserId: "userID1",
			PostId: "postID1",
			Context: model.StringInterface{
				"orig_user_id":  "foobar",
				"reminder_id":   model.NewId(),
				"occurrence_id": model.NewId(),
			},
		}

		w := httptest.NewRecorder()
		requestJSON, _ := json.Marshal(request)
		r := httptest.NewRequest("POST", "/complete", bytes.NewReader(requestJSON))
		p.ServeHTTP(nil, w, r)

		result := w.Result()
		assert.NotNil(t, result)

		bodyBytes, err := io.ReadAll(result.Body)
		assert.Nil(t, err)
		bodyString := string(bodyBytes)
		assert.Equal(t, bodyString, "{\"update\":null,\"ephemeral_text\":\"\",\"skip_slack_parsing\":false}")

	})

}

func TestHandleDelete(t *testing.T) {

	user := &model.User{
		Id:       model.NewId(),
		Username: model.NewRandomString(10),
	}
	post := &model.Post{
		Id:        model.NewId(),
		ChannelId: model.NewId(),
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
		api.On("GetPost", mock.Anything).Return(post, nil)
		api.On("UpdatePost", mock.Anything).Return(post, nil)
		api.On("GetUser", mock.Anything).Return(user, nil)
		api.On("GetUserByUsername", mock.Anything).Return(user, nil)
		api.On("KVGet", user.Username).Return(stringReminders, nil)
		api.On("KVSet", mock.Anything, mock.Anything).Return(nil)

		return api
	}

	t.Run("delete", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.router = p.InitAPI()
		p.API = api

		request := &model.PostActionIntegrationRequest{
			UserId: "userID1",
			PostId: "postID1",
			Context: model.StringInterface{
				"orig_user_id":  "foobar",
				"reminder_id":   model.NewId(),
				"occurrence_id": model.NewId(),
			},
		}

		w := httptest.NewRecorder()
		requestJSON, _ := json.Marshal(request)
		r := httptest.NewRequest("POST", "/delete", bytes.NewReader(requestJSON))
		p.ServeHTTP(nil, w, r)

		result := w.Result()
		assert.NotNil(t, result)

		bodyBytes, err := io.ReadAll(result.Body)
		assert.Nil(t, err)
		bodyString := string(bodyBytes)
		assert.Equal(t, bodyString, "{\"update\":null,\"ephemeral_text\":\"\",\"skip_slack_parsing\":false}")

	})

}

func TestHandleDeleteEphemeral(t *testing.T) {
	user := &model.User{
		Id:       model.NewId(),
		Username: model.NewRandomString(10),
	}
	post := &model.Post{
		Id:        model.NewId(),
		ChannelId: model.NewId(),
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
		api.On("GetUser", mock.Anything).Return(user, nil)
		api.On("GetUserByUsername", mock.Anything).Return(user, nil)
		api.On("KVGet", user.Username).Return(stringReminders, nil)
		api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
		api.On("UpdateEphemeralPost", mock.Anything, mock.Anything).Return(post)

		return api
	}

	t.Run("delete ephemeral", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.router = p.InitAPI()
		p.API = api

		request := &model.PostActionIntegrationRequest{
			UserId: "userID1",
			PostId: "postID1",
			Context: model.StringInterface{
				"reminder_id":   model.NewId(),
				"occurrence_id": model.NewId(),
			},
		}

		w := httptest.NewRecorder()
		requestJSON, _ := json.Marshal(request)
		r := httptest.NewRequest("POST", "/delete/ephemeral", bytes.NewReader(requestJSON))
		p.ServeHTTP(nil, w, r)

		result := w.Result()
		assert.NotNil(t, result)

		bodyBytes, err := io.ReadAll(result.Body)
		assert.Nil(t, err)
		bodyString := string(bodyBytes)
		assert.Equal(t, bodyString, "{\"update\":null,\"ephemeral_text\":\"\",\"skip_slack_parsing\":false}")

	})

}

func TestHandleSnooze(t *testing.T) {

	user := &model.User{
		Id:       model.NewId(),
		Username: model.NewRandomString(10),
	}
	post := &model.Post{
		Id:        model.NewId(),
		ChannelId: model.NewId(),
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
	stringOccurrences, _ := json.Marshal(occurrences)
	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("GetPost", mock.Anything).Return(post, nil)
		api.On("UpdatePost", mock.Anything).Return(post, nil)
		api.On("GetUser", mock.Anything).Return(user, nil)
		api.On("GetUserByUsername", mock.Anything).Return(user, nil)
		api.On("KVGet", user.Username).Return(stringReminders, nil)
		api.On("KVGet", mock.Anything).Return(stringOccurrences, nil)
		api.On("KVSet", mock.Anything, mock.Anything).Return(nil)

		return api
	}

	for name, test := range map[string]struct {
		SnoozeTime string
	}{
		"snoozes item 20min": {
			SnoozeTime: "20min",
		},
		"snoozes item 1hr": {
			SnoozeTime: "1hr",
		},
		"snoozes item 3hrs": {
			SnoozeTime: "3hrs",
		},
		"snoozes item tomorrow": {
			SnoozeTime: "tomorrow",
		},
		"snoozes item nextweek": {
			SnoozeTime: "nextweek",
		},
	} {

		t.Run(name, func(t *testing.T) {

			api := setupAPI()
			defer api.AssertExpectations(t)

			p := &Plugin{}
			p.router = p.InitAPI()
			p.API = api

			request := &model.PostActionIntegrationRequest{
				UserId: "userID1",
				PostId: "postID1",
				Context: model.StringInterface{
					"orig_user_id":    "foobar",
					"reminder_id":     reminders[0].Id,
					"occurrence_id":   occurrences[0].Id,
					"selected_option": test.SnoozeTime,
				},
			}

			w := httptest.NewRecorder()
			requestJSON, _ := json.Marshal(request)
			r := httptest.NewRequest("POST", "/snooze", bytes.NewReader(requestJSON))
			p.ServeHTTP(nil, w, r)

			result := w.Result()
			assert.NotNil(t, result)

			bodyBytes, err := io.ReadAll(result.Body)
			assert.Nil(t, err)
			bodyString := string(bodyBytes)
			assert.Equal(t, bodyString, "{\"update\":null,\"ephemeral_text\":\"\",\"skip_slack_parsing\":false}")

		})
	}
}

func TestHandleNextReminders(t *testing.T) {

	user := &model.User{
		Id:       model.NewId(),
		Username: model.NewRandomString(10),
	}
	post := &model.Post{
		Id:        model.NewId(),
		ChannelId: model.NewId(),
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
		api.On("GetUser", mock.Anything).Return(user, nil)
		api.On("GetUserByUsername", mock.Anything).Return(user, nil)
		api.On("KVGet", user.Username).Return(stringReminders, nil)
		api.On("UpdateEphemeralPost", mock.Anything, mock.Anything).Return(post)

		return api
	}

	t.Run("next reminders", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.router = p.InitAPI()
		p.API = api

		request := &model.PostActionIntegrationRequest{
			UserId: user.Id,
			PostId: post.Id,
			Context: model.StringInterface{
				"reminder_id":   model.NewId(),
				"occurrence_id": model.NewId(),
				"offset":        0,
			},
		}

		w := httptest.NewRecorder()
		requestJSON, _ := json.Marshal(request)
		r := httptest.NewRequest("POST", "/next/reminders", bytes.NewReader(requestJSON))
		p.ServeHTTP(nil, w, r)

		result := w.Result()
		assert.NotNil(t, result)

		bodyBytes, err := io.ReadAll(result.Body)
		assert.Nil(t, err)
		bodyString := string(bodyBytes)
		assert.Equal(t, bodyString, "{\"update\":null,\"ephemeral_text\":\"\",\"skip_slack_parsing\":false}")

	})
}

func TestHandleCompleteList(t *testing.T) {

	user := &model.User{
		Id:       model.NewId(),
		Username: model.NewRandomString(10),
	}

	post := &model.Post{
		Id:        model.NewId(),
		ChannelId: model.NewId(),
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
		api.On("GetUser", mock.Anything).Return(user, nil)
		api.On("GetUserByUsername", mock.Anything).Return(user, nil)
		api.On("KVGet", user.Username).Return(stringReminders, nil)
		api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
		api.On("UpdateEphemeralPost", mock.Anything, mock.Anything).Return(post)

		return api
	}

	t.Run("complete list", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.router = p.InitAPI()
		p.API = api

		request := &model.PostActionIntegrationRequest{
			UserId: "userID1",
			PostId: "postID1",
			Context: model.StringInterface{
				"reminder_id":   model.NewId(),
				"occurrence_id": model.NewId(),
				"offset":        0,
			},
		}

		w := httptest.NewRecorder()
		requestJSON, _ := json.Marshal(request)
		r := httptest.NewRequest("POST", "/complete/list", bytes.NewReader(requestJSON))
		p.ServeHTTP(nil, w, r)

		result := w.Result()
		assert.NotNil(t, result)

		bodyBytes, err := io.ReadAll(result.Body)
		assert.Nil(t, err)
		bodyString := string(bodyBytes)
		assert.Equal(t, bodyString, "{\"update\":null,\"ephemeral_text\":\"\",\"skip_slack_parsing\":false}")

	})
}

func TestHandleViewCompleteList(t *testing.T) {

	user := &model.User{
		Id:       model.NewId(),
		Username: model.NewRandomString(10),
	}
	post := &model.Post{
		Id:        model.NewId(),
		ChannelId: model.NewId(),
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
		api.On("GetUser", mock.Anything).Return(user, nil)
		api.On("GetUserByUsername", mock.Anything).Return(user, nil)
		api.On("KVGet", user.Username).Return(stringReminders, nil)
		api.On("UpdateEphemeralPost", mock.Anything, mock.Anything).Return(post)

		return api
	}

	t.Run("view complete list", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.router = p.InitAPI()
		p.API = api

		request := &model.PostActionIntegrationRequest{UserId: user.Id, PostId: post.Id}

		w := httptest.NewRecorder()
		requestJSON, _ := json.Marshal(request)
		r := httptest.NewRequest("POST", "/view/complete/list", bytes.NewReader(requestJSON))
		p.ServeHTTP(nil, w, r)

		result := w.Result()
		assert.NotNil(t, result)

		bodyBytes, err := io.ReadAll(result.Body)
		assert.Nil(t, err)
		bodyString := string(bodyBytes)
		assert.Equal(t, bodyString, "{\"update\":null,\"ephemeral_text\":\"\",\"skip_slack_parsing\":false}")

	})

}

func TestHandleDeleteList(t *testing.T) {
	user := &model.User{
		Id:       model.NewId(),
		Username: model.NewRandomString(10),
	}
	post := &model.Post{
		Id:        model.NewId(),
		ChannelId: model.NewId(),
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
		api.On("GetUser", mock.Anything).Return(user, nil)
		api.On("GetUserByUsername", mock.Anything).Return(user, nil)
		api.On("KVGet", user.Username).Return(stringReminders, nil)
		api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
		api.On("UpdateEphemeralPost", mock.Anything, mock.Anything).Return(post)

		return api
	}

	t.Run("delete list item", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.router = p.InitAPI()
		p.API = api

		request := &model.PostActionIntegrationRequest{
			UserId: "userID1",
			PostId: "postID1",
			Context: model.StringInterface{
				"reminder_id":   model.NewId(),
				"occurrence_id": model.NewId(),
				"offset":        0,
			},
		}

		w := httptest.NewRecorder()
		requestJSON, _ := json.Marshal(request)
		r := httptest.NewRequest("POST", "/delete/list", bytes.NewReader(requestJSON))
		p.ServeHTTP(nil, w, r)

		result := w.Result()
		assert.NotNil(t, result)

		bodyBytes, err := io.ReadAll(result.Body)
		assert.Nil(t, err)
		bodyString := string(bodyBytes)
		assert.Equal(t, bodyString, "{\"update\":null,\"ephemeral_text\":\"\",\"skip_slack_parsing\":false}")

	})

}

func TestHandleDeleteCompleteList(t *testing.T) {

	user := &model.User{
		Id:       model.NewId(),
		Username: model.NewRandomString(10),
	}
	post := &model.Post{
		Id:        model.NewId(),
		ChannelId: model.NewId(),
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
		api.On("GetUser", mock.Anything).Return(user, nil)
		api.On("GetUserByUsername", mock.Anything).Return(user, nil)
		api.On("UpdateEphemeralPost", mock.Anything, mock.Anything).Return(post)
		api.On("KVGet", user.Username).Return(stringReminders, nil)

		return api
	}

	t.Run("delete completed list", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.router = p.InitAPI()
		p.API = api

		request := &model.PostActionIntegrationRequest{
			UserId: "userID1",
			PostId: "postID1",
			Context: model.StringInterface{
				"reminder_id":   model.NewId(),
				"occurrence_id": model.NewId(),
				"offset":        0,
			},
		}

		w := httptest.NewRecorder()
		requestJSON, _ := json.Marshal(request)
		r := httptest.NewRequest("POST", "/delete/complete/list", bytes.NewReader(requestJSON))
		p.ServeHTTP(nil, w, r)

		result := w.Result()
		assert.NotNil(t, result)

		bodyBytes, err := io.ReadAll(result.Body)
		assert.Nil(t, err)
		bodyString := string(bodyBytes)
		assert.Equal(t, bodyString, "{\"update\":null,\"ephemeral_text\":\"\",\"skip_slack_parsing\":false}")

	})
}

func TestHandleSnoozeList(t *testing.T) {

	user := &model.User{
		Id:       model.NewId(),
		Username: model.NewRandomString(10),
	}
	post := &model.Post{
		Id:        model.NewId(),
		ChannelId: model.NewId(),
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
	stringOccurrences, _ := json.Marshal(occurrences)
	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("GetUser", mock.Anything).Return(user, nil)
		api.On("GetUserByUsername", mock.Anything).Return(user, nil)
		api.On("KVGet", user.Username).Return(stringReminders, nil)
		api.On("KVGet", mock.Anything).Return(stringOccurrences, nil)
		api.On("KVSet", mock.Anything, mock.Anything).Return(nil)
		api.On("UpdateEphemeralPost", mock.Anything, mock.Anything).Return(post)

		return api
	}

	for name, test := range map[string]struct {
		SnoozeTime string
	}{
		"snoozes list item 20min": {
			SnoozeTime: "20min",
		},
		"snoozes list item 1hr": {
			SnoozeTime: "1hr",
		},
		"snoozes list item 3hrs": {
			SnoozeTime: "3hrs",
		},
		"snoozes list item tomorrow": {
			SnoozeTime: "tomorrow",
		},
		"snoozes list item nextweek": {
			SnoozeTime: "nextweek",
		},
	} {

		t.Run(name, func(t *testing.T) {

			api := setupAPI()
			defer api.AssertExpectations(t)

			p := &Plugin{}
			p.router = p.InitAPI()
			p.API = api

			request := &model.PostActionIntegrationRequest{
				UserId: "userID1",
				PostId: "postID1",
				Context: model.StringInterface{
					"reminder_id":     reminders[0].Id,
					"occurrence_id":   occurrences[0].Id,
					"selected_option": test.SnoozeTime,
					"offset":          0,
				},
			}

			w := httptest.NewRecorder()
			requestJSON, _ := json.Marshal(request)
			r := httptest.NewRequest("POST", "/snooze/list", bytes.NewReader(requestJSON))
			p.ServeHTTP(nil, w, r)

			result := w.Result()
			assert.NotNil(t, result)

			bodyBytes, err := io.ReadAll(result.Body)
			assert.Nil(t, err)
			bodyString := string(bodyBytes)
			assert.Equal(t, bodyString, "{\"update\":null,\"ephemeral_text\":\"\",\"skip_slack_parsing\":false}")

		})
	}
}

func TestHandleCloseList(t *testing.T) {

	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogError", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		api.On("DeleteEphemeralPost", mock.Anything, mock.Anything).Return(nil)
		return api
	}

	t.Run("closes list", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.router = p.InitAPI()
		p.API = api

		request := &model.PostActionIntegrationRequest{UserId: "userID1", PostId: "postID1"}

		w := httptest.NewRecorder()
		requestJSON, _ := json.Marshal(request)
		r := httptest.NewRequest("POST", "/close/list", bytes.NewReader(requestJSON))
		p.ServeHTTP(nil, w, r)

		result := w.Result()
		assert.NotNil(t, result)

		bodyBytes, err := io.ReadAll(result.Body)
		assert.Nil(t, err)
		bodyString := string(bodyBytes)
		assert.Equal(t, bodyString, "{\"update\":null,\"ephemeral_text\":\"\",\"skip_slack_parsing\":false}")

	})
}
