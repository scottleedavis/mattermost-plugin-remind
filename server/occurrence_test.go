package main

import (
	"testing"
	"time"
	// "fmt"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestIn(t *testing.T) {

	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		return api
	}

	t.Run("if inEN locale is used", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		user := &model.User{
			Email:    "-@-.-",
			Nickname: "TestUser",
			Password: model.NewId(),
			Username: "testuser",
			Roles:    model.SYSTEM_USER_ROLE_ID,
			Locale:   "en",
		}

		times, err := p.inEN("in one second", user)
		testTimes := []time.Time{
			time.Now().Round(time.Second).Add(time.Second * time.Duration(1)).UTC(),
		}
		assert.Nil(t, err)
		assert.Equal(t, times, testTimes)

		times, err = p.inEN("in 712 minutes", user)
		testTimes = []time.Time{
			time.Now().Round(time.Second).Add(time.Minute * time.Duration(712)).UTC(),
		}
		assert.Nil(t, err)
		assert.Equal(t, times, testTimes)

		times, err = p.inEN("in 3 hours", user)
		testTimes = []time.Time{
			time.Now().Round(time.Second).Add(time.Hour * time.Duration(3)).UTC(),
		}
		assert.Nil(t, err)
		assert.Equal(t, times, testTimes)

		times, err = p.inEN("in 24 hours", user)
		testTimes = []time.Time{
			time.Now().Round(time.Second).Add(time.Hour * time.Duration(24)).UTC(),
		}
		assert.Nil(t, err)
		assert.Equal(t, times, testTimes)

		times, err = p.inEN("in 2 days", user)
		testTimes = []time.Time{
			time.Now().Round(time.Second).Add(time.Hour * time.Duration(24) * time.Duration(2)).UTC(),
		}
		assert.Nil(t, err)
		assert.Equal(t, times, testTimes)

		times, err = p.inEN("in 90 weeks", user)
		testTimes = []time.Time{
			time.Now().Round(time.Second).Add(time.Hour * time.Duration(24) * time.Duration(7) * time.Duration(90)).UTC(),
		}
		assert.Nil(t, err)
		assert.Equal(t, times, testTimes)

		times, err = p.inEN("in 4 months", user)
		testTimes = []time.Time{
			time.Now().Round(time.Second).Add(time.Hour * time.Duration(24) * time.Duration(30) * time.Duration(4)).UTC(),
		}
		assert.Nil(t, err)
		assert.Equal(t, times, testTimes)

		times, err = p.inEN("in 1 year", user)
		testTimes = []time.Time{
			time.Now().Round(time.Second).Add(time.Hour * time.Duration(24) * time.Duration(365)).UTC(),
		}
		assert.Nil(t, err)
		assert.Equal(t, times, testTimes)

	})

}

func TestAt(t *testing.T) {
	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		return api
	}

	t.Run("if atEN locale is used", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		user := &model.User{
			Email:    "-@-.-",
			Nickname: "TestUser",
			Password: model.NewId(),
			Username: "testuser",
			Roles:    model.SYSTEM_USER_ROLE_ID,
			Locale:   "en",
		}
		location := p.location(user)

		times, err := p.atEN("at noon", user)
		assert.Nil(t, err)
		assert.Equal(t, times[0].In(location).Hour(), 12)

		times, err = p.atEN("at midnight", user)
		assert.Nil(t, err)
		assert.Equal(t, times[0].In(location).Hour(), 0)

		times, err = p.atEN("at two", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Hour() == 2 || times[0].In(location).Hour() == 14)

		times, err = p.atEN("at 7", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Hour() == 7 || times[0].In(location).Hour() == 19)

		times, err = p.atEN("at 12:30pm", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Hour() == 12 && times[0].In(location).Minute() == 30)

		times, err = p.atEN("at 7:12pm", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Hour() == 19 && times[0].In(location).Minute() == 12)

		times, err = p.atEN("at 8:05 PM", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Hour() == 20 && times[0].In(location).Minute() == 5)

		times, err = p.atEN("at 9:52 am", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Hour() == 9 && times[0].In(location).Minute() == 52)

		times, err = p.atEN("at 9:12", user)
		assert.Nil(t, err)
		assert.True(t, (times[0].In(location).Hour() == 9 || times[0].In(location).Hour() == 21) && times[0].In(location).Minute() == 12)

		times, err = p.atEN("at 17:15", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Hour() == 17 && times[0].In(location).Minute() == 15)

		times, err = p.atEN("at 930am", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Hour() == 9 && times[0].In(location).Minute() == 30)

		times, err = p.atEN("at 1230 am", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Hour() == 0 && times[0].In(location).Minute() == 30)

		times, err = p.atEN("at 5PM", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Hour() == 17 && times[0].In(location).Minute() == 0)

		times, err = p.atEN("at 4 am", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Hour() == 4 && times[0].In(location).Minute() == 0)

		times, err = p.atEN("at 1400", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Hour() == 14 && times[0].In(location).Minute() == 0)

		times, err = p.atEN("at 11:00 every Thursday", user)
		assert.Nil(t, err)
		if err == nil {
			assert.True(t, (times[0].In(location).Hour() == 11 || times[0].In(location).Hour() == 23) &&
				times[0].In(location).Weekday().String() == "Thursday")
		}

		times, err = p.atEN("at 3pm every day", user)
		assert.Nil(t, err)
		if err == nil {
			assert.True(t, times[0].In(location).Hour() == 15)
		}
	})
}

func TestOn(t *testing.T) {
	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		return api
	}

	t.Run("if onEN locale is used", func(t *testing.T) {

		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		user := &model.User{
			Email:    "-@-.-",
			Nickname: "TestUser",
			Password: model.NewId(),
			Username: "testuser",
			Roles:    model.SYSTEM_USER_ROLE_ID,
			Locale:   "en",
		}
		location := p.location(user)

		times, err := p.onEN("on Monday", user)
		assert.Nil(t, err)
		assert.Equal(t, times[0].In(location).Weekday().String(), "Monday")

		times, err = p.onEN("on tuesday", user)
		assert.Nil(t, err)
		assert.Equal(t, times[0].In(location).Weekday().String(), "Tuesday")

		times, err = p.onEN("on WedNesday", user)
		assert.Nil(t, err)
		assert.Equal(t, times[0].In(location).Weekday().String(), "Wednesday")

		times, err = p.onEN("on thursDAY", user)
		assert.Nil(t, err)
		assert.Equal(t, times[0].In(location).Weekday().String(), "Thursday")

		times, err = p.onEN("on FrIdAy", user)
		assert.Nil(t, err)
		assert.Equal(t, times[0].In(location).Weekday().String(), "Friday")

		times, err = p.onEN("on Saturday", user)
		assert.Nil(t, err)
		assert.Equal(t, times[0].In(location).Weekday().String(), "Saturday")

		times, err = p.onEN("on sunday", user)
		assert.Nil(t, err)
		assert.Equal(t, times[0].In(location).Weekday().String(), "Sunday")

		times, err = p.onEN("on Mondays", user)
		assert.Nil(t, err)
		if err == nil {
			assert.Equal(t, times[0].In(location).Weekday().String(), "Monday")
		}

		times, err = p.onEN("on Tuesdays", user)
		assert.Nil(t, err)
		if err == nil {
			assert.Equal(t, times[0].In(location).Weekday().String(), "Tuesday")
		}

		times, err = p.onEN("on Wednesdays", user)
		assert.Nil(t, err)
		if err == nil {
			assert.Equal(t, times[0].In(location).Weekday().String(), "Wednesday")
		}

		times, err = p.onEN("on Thursdays", user)
		assert.Nil(t, err)
		if err == nil {
			assert.Equal(t, times[0].In(location).Weekday().String(), "Thursday")
		}

		times, err = p.onEN("on Fridays", user)
		assert.Nil(t, err)
		if err == nil {
			assert.Equal(t, times[0].In(location).Weekday().String(), "Friday")
		}

		times, err = p.onEN("on Saturdays", user)
		assert.Nil(t, err)
		if err == nil {
			assert.Equal(t, times[0].In(location).Weekday().String(), "Saturday")
		}

		times, err = p.onEN("on Sundays", user)
		assert.Nil(t, err)
		if err == nil {
			assert.Equal(t, times[0].In(location).Weekday().String(), "Sunday")
		}

		times, err = p.onEN("on mon", user)
		assert.Nil(t, err)
		assert.Equal(t, times[0].In(location).Weekday().String(), "Monday")

		times, err = p.onEN("on wED", user)
		assert.Nil(t, err)
		assert.Equal(t, times[0].In(location).Weekday().String(), "Wednesday")

		times, err = p.onEN("on tuesday at noon", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Weekday().String() == "Tuesday" && times[0].In(location).Hour() == 12)

		times, err = p.onEN("on sunday at 3:42am", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Weekday().String() == "Sunday" && times[0].In(location).Hour() == 3 &&
			times[0].In(location).Minute() == 42)

		times, err = p.onEN("on December 15", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Month().String() == "December" && times[0].In(location).Day() == 15)

		times, err = p.onEN("on jan 12", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Month().String() == "January" && times[0].In(location).Day() == 12)

		times, err = p.onEN("on july 12th", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Month().String() == "July" && times[0].In(location).Day() == 12)

		times, err = p.onEN("on mArch 22nd", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Month().String() == "March" && times[0].In(location).Day() == 22)

		times, err = p.onEN("on march 17 at 5:41pm", user)
		assert.Nil(t, err)
		if err == nil {
			assert.True(t, times[0].In(location).Month().String() == "March" && times[0].In(location).Day() == 17 &&
				times[0].In(location).Hour() == 17 && times[0].In(location).Minute() == 41)
		}

		times, err = p.onEN("on september 7th 2020", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Month().String() == "September" && times[0].In(location).Day() == 7)

		times, err = p.onEN("on April 17 2020", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Month().String() == "April" && times[0].In(location).Day() == 17)

		times, err = p.onEN("on April 9 2020 at 11am", user)
		assert.Nil(t, err)
		if err == nil {
			assert.True(t, times[0].In(location).Month().String() == "April" && times[0].In(location).Day() == 9 &&
				times[0].In(location).Hour() == 11)
		}

		times, err = p.onEN("on auguSt tenth 2019", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Month().String() == "August" && times[0].In(location).Day() == 10)

		times, err = p.onEN("on 7", user)
		assert.Nil(t, err)
		assert.Equal(t, times[0].In(location).Day(), 7)

		times, err = p.onEN("on 7th", user)
		assert.Nil(t, err)
		assert.Equal(t, times[0].In(location).Day(), 7)

		times, err = p.onEN("on seven", user)
		assert.Nil(t, err)
		assert.Equal(t, times[0].In(location).Day(), 7)

		times, err = p.onEN("on 1/17/20", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Year() == 2020 && times[0].In(location).Month() == 1 &&
			times[0].In(location).Day() == 17)

		times, err = p.onEN("on 12/17/2020", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Year() == 2020 && times[0].In(location).Month() == 12 &&
			times[0].In(location).Day() == 17)

		times, err = p.onEN("on 17.1.20", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Year() == 2020 && times[0].In(location).Month() == 1 &&
			times[0].In(location).Day() == 17)

		times, err = p.onEN("on 17.12.20", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Year() == 2020 && times[0].In(location).Month() == 12 &&
			times[0].In(location).Day() == 17)

		times, err = p.onEN("on 12/1", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Month() == 12 && times[0].In(location).Day() == 1)

		times, err = p.onEN("on 5-17-20", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Month() == 5 && times[0].In(location).Day() == 17)

		times, err = p.onEN("on 12-5-2020", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Month() == 12 && times[0].In(location).Day() == 5)

		times, err = p.onEN("on 12-12", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Month() == 12 && times[0].In(location).Day() == 12)

		times, err = p.onEN("on 1-1 at midnight", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Month() == 1 && times[0].In(location).Day() == 1 &&
			times[0].In(location).Hour() == 0)

	})
}

func TestEvery(t *testing.T) {

	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		return api
	}

	t.Run("if everyEN locale is used", func(t *testing.T) {
		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		user := &model.User{
			Email:    "-@-.-",
			Nickname: "TestUser",
			Password: model.NewId(),
			Username: "testuser",
			Roles:    model.SYSTEM_USER_ROLE_ID,
			Locale:   "en",
		}
		location := p.location(user)

		times, err := p.everyEN("every Thursday", user)
		assert.Nil(t, err)
		assert.Equal(t, times[0].In(location).Weekday().String(), "Thursday")

		times, err = p.everyEN("every day", user)
		assert.Nil(t, err)
		assert.Equal(t, times[0].In(location).Weekday().String(), time.Now().In(location).AddDate(0, 0, 1).Weekday().String())

		times, err = p.everyEN("every 12/18/2022", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Month() == 12 && times[0].In(location).Year() == 2022)

		times, err = p.everyEN("every january 25", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Month() == 1 && times[0].In(location).Day() == 25)

		times, err = p.everyEN("every other wednesday", user)
		assert.Nil(t, err)
		assert.Equal(t, times[0].In(location).Weekday().String(), "Wednesday")

		times, err = p.everyEN("every day at 11:32am", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Hour() == 11 && times[0].In(location).Minute() == 32)

		times, err = p.everyEN("every 5/5 at 7", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Month() == 5 && times[0].In(location).Day() == 5 &&
			(times[0].In(location).Hour() == 7 || times[0].In(location).Hour() == 19))

		times, err = p.everyEN("every 7/20 at 1100", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Month() == 7 && times[0].In(location).Day() == 20 &&
			(times[0].In(location).Hour() == 11 || times[0].In(location).Hour() == 23))

		times, err = p.everyEN("every Monday at 7:32am", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Weekday().String() == "Monday" && (times[0].In(location).Hour() == 7 ||
			times[0].In(location).Hour() == 32))

		times, err = p.everyEN("every monday and wednesday", user)
		assert.Nil(t, err)
		if err == nil {
			assert.True(t, times[0].In(location).Weekday().String() == "Monday" && times[1].In(location).Weekday().String() == "Wednesday")
		}

		times, err = p.everyEN("every wednesday, thursday", user)
		assert.Nil(t, err)
		if err == nil {
			assert.True(t, times[0].In(location).Weekday().String() == "Monday" && times[1].In(location).Weekday().String() == "Thursday")
		}

		times, err = p.everyEN("every other friday and saturday", user)
		assert.Nil(t, err)
		if err == nil {
			assert.True(t, times[0].In(location).Weekday().String() == "Friday" && times[1].In(location).Weekday().String() == "Saturday")
		}

		times, err = p.everyEN("every monday and wednesday at 1:39am", user)
		assert.Nil(t, err)
		if err == nil {
			assert.True(t, times[0].In(location).Weekday().String() == "Monday" &&
				times[1].In(location).Weekday().String() == "Wednesday" && times[0].In(location).Hour() == 1 && times[0].Minute() == 39)
		}

		times, err = p.everyEN("every monday, tuesday and sunday at 11:00am", user)
		assert.Nil(t, err)
		if err == nil {
			assert.True(t, times[0].In(location).Weekday().String() == "Monday" &&
				times[1].In(location).Weekday().String() == "Tuesday" && times[2].In(location).Weekday().String() == "Sunday" &&
				times[0].In(location).Hour() == 11)
		}

		times, err = p.everyEN("every monday, tuesday at 2pm", user)
		assert.Nil(t, err)
		if err == nil {
			assert.True(t, times[0].In(location).Weekday().String() == "Monday" &&
				times[1].In(location).Weekday().String() == "Tuesday" && times[0].In(location).Hour() == 14)
		}

		times, err = p.everyEN("every 1/30 and 9/30 at noon", user)
		assert.Nil(t, err)
		if err == nil {
			assert.True(t, times[0].In(location).Month() == 1 && times[0].In(location).Day() == 30 &&
				times[1].In(location).Month() == 9 && times[1].In(location).Day() == 30 && times[0].In(location).Hour() == 12)
		}
	})
}

func TestFreeForm(t *testing.T) {

	setupAPI := func() *plugintest.API {
		api := &plugintest.API{}
		api.On("LogDebug", mock.Anything, mock.Anything, mock.Anything).Maybe()
		api.On("LogInfo", mock.Anything).Maybe()
		return api
	}

	t.Run("if freeformEN locale is used", func(t *testing.T) {
		api := setupAPI()
		defer api.AssertExpectations(t)

		p := &Plugin{}
		p.API = api

		user := &model.User{
			Email:    "-@-.-",
			Nickname: "TestUser",
			Password: model.NewId(),
			Username: "testuser",
			Roles:    model.SYSTEM_USER_ROLE_ID,
			Locale:   "en",
		}
		location := p.location(user)

		times, err := p.freeFormEN("monday", user)
		assert.Nil(t, err)
		assert.Equal(t, times[0].In(location).Weekday().String(), "Monday")

		times, err = p.freeFormEN("tuesday at 9:34pm", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Weekday().String() == "Tuesday" && times[0].In(location).Hour() == 21 &&
			times[0].In(location).Minute() == 34)

		times, err = p.freeFormEN("wednesday", user)
		assert.Nil(t, err)
		assert.Equal(t, times[0].In(location).Weekday().String(), "Wednesday")

		times, err = p.freeFormEN("thursday at noon", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Weekday().String() == "Thursday" && times[0].In(location).Hour() == 12)

		times, err = p.freeFormEN("friday", user)
		assert.Nil(t, err)
		assert.Equal(t, times[0].In(location).Weekday().String(), "Friday")

		times, err = p.freeFormEN("saturday", user)
		assert.Nil(t, err)
		assert.Equal(t, times[0].In(location).Weekday().String(), "Saturday")

		times, err = p.freeFormEN("sunday at 4:20pm", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Weekday().String() == "Sunday" && times[0].In(location).Hour() == 16 &&
			times[0].In(location).Minute() == 20)

		times, err = p.freeFormEN("today at 3pm", user)
		assert.Nil(t, err)
		if err == nil {
			assert.Equal(t, times[0].In(location).Hour(), 15)
		}

		times, err = p.freeFormEN("tomorrow", user)
		assert.Nil(t, err)
		assert.True(t, times[0].In(location).Weekday().String() == time.Now().In(location).AddDate(0, 0, 1).Weekday().String())

		times, err = p.freeFormEN("tomorrow at 4pm", user)
		assert.Nil(t, err)
		if err == nil {
			assert.True(t, times[0].In(location).Weekday().String() == time.Now().In(location).AddDate(0, 0, 1).Weekday().String() &&
				times[0].In(location).Hour() == 16)
		}

		times, err = p.freeFormEN("everyday", user)
		assert.Nil(t, err)
		if err == nil {
			assert.Equal(t, times[0].In(location).Weekday().String(), time.Now().In(location).AddDate(0, 0, 1).Weekday().String())
		}

		times, err = p.freeFormEN("everyday at 3:23am", user)
		assert.Nil(t, err)
		if err == nil {
			assert.True(t, times[0].In(location).Weekday().String() == time.Now().In(location).AddDate(0, 0, 1).Weekday().String() &&
				times[0].In(location).Hour() == 3 && times[0].In(location).Minute() == 23)
		}

		times, err = p.freeFormEN("mondays", user)
		assert.Nil(t, err)
		if err == nil {
			assert.Equal(t, times[0].In(location).Weekday().String(), "Monday")
		}

		times, err = p.freeFormEN("tuesdays", user)
		assert.Nil(t, err)
		if err == nil {
			assert.Equal(t, times[0].In(location).Weekday().String(), "Tuesday")
		}

		times, err = p.freeFormEN("wednesdays", user)
		assert.Nil(t, err)
		if err == nil {
			assert.Equal(t, times[0].In(location).Weekday().String(), "Wednesday")
		}

		times, err = p.freeFormEN("Thursdays", user)
		assert.Nil(t, err)
		if err == nil {
			assert.Equal(t, times[0].In(location).Weekday().String(), "Thursday")
		}

		times, err = p.freeFormEN("Fridays", user)
		assert.Nil(t, err)
		if err == nil {
			assert.Equal(t, times[0].In(location).Weekday().String(), "Friday")
		}

		times, err = p.freeFormEN("Saturdays", user)
		assert.Nil(t, err)
		if err == nil {
			assert.Equal(t, times[0].In(location).Weekday().String(), "Saturday")
		}

		times, err = p.freeFormEN("Sundays", user)
		assert.Nil(t, err)
		if err == nil {
			assert.Equal(t, times[0].In(location).Weekday().String(), "Sunday")
		}

	})
}
