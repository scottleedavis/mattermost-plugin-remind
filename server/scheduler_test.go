package main

import (
	"testing"
	// "github.com/mattermost/mattermost-server/model"
	// "github.com/mattermost/mattermost-server/plugin/plugintest"
	// "github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/mock"
)

func TestScheduleReminders(t *testing.T) {
	// th := Setup()
	// defer th.TearDown()

	// th.Server.InitReminders()
	// defer th.Server.StopReminders()
	// user, _ := th.App.GetUserByUsername(model.REMIND_BOTNAME)
	// T := utils.GetUserTranslations(user.Locale)

	// request := &model.ReminderRequest{}
	// request.UserId = user.Id

	// request.Payload = "me foo in 1 seconds"
	// response, err := th.App.ScheduleReminder(request)
	// if err != nil {
	// 	t.Fatal(UNABLE_TO_SCHEDULE_REMINDER)
	// }

	// t2 := time.Now().Add(2 * time.Second).Format(time.Kitchen)
	// var responseParameters = map[string]interface{}{
	// 	"Target":  "You",
	// 	"UseTo":   "",
	// 	"Message": "foo",
	// 	"When":    "in 1 seconds at " + t2 + " today.",
	// }
	// expectedResponse := T("app.reminder.response", responseParameters)
	// assert.Equal(t, response, expectedResponse)

	// request.Payload = "@bob foo in 1 seconds"
	// request.Occurrences = model.Occurrences{}
	// response, err = th.App.ScheduleReminder(request)
	// if err != nil {
	// 	t.Fatal(UNABLE_TO_SCHEDULE_REMINDER)
	// }
	// t2 = time.Now().Add(time.Second).Format(time.Kitchen)
	// responseParameters = map[string]interface{}{
	// 	"Target":  "@bob",
	// 	"UseTo":   "",
	// 	"Message": "foo",
	// 	"When":    "in 1 seconds at " + t2 + " today.",
	// }
	// expectedResponse = T("app.reminder.response", responseParameters)
	// assert.Equal(t, response, expectedResponse)

	// request.Payload = "~off-topic foo in 1 seconds"
	// request.Occurrences = model.Occurrences{}
	// response, err = th.App.ScheduleReminder(request)
	// if err != nil {
	// 	t.Fatal(UNABLE_TO_SCHEDULE_REMINDER)
	// }
	// t2 = time.Now().Add(time.Second).Format(time.Kitchen)

	// responseParameters = map[string]interface{}{
	// 	"Target":  "~off-topic",
	// 	"UseTo":   "",
	// 	"Message": "foo",
	// 	"When":    "in 1 seconds at " + t2 + " today.",
	// }
	// expectedResponse = T("app.reminder.response", responseParameters)
	// assert.Equal(t, response, expectedResponse)

	// request.Payload = "me \"foo foo foo\" in 1 seconds"
	// request.Occurrences = model.Occurrences{}
	// response, err = th.App.ScheduleReminder(request)
	// if err != nil {
	// 	t.Fatal(UNABLE_TO_SCHEDULE_REMINDER)
	// }
	// t2 = time.Now().Add(time.Second).Format(time.Kitchen)

	// responseParameters = map[string]interface{}{
	// 	"Target":  "You",
	// 	"UseTo":   "",
	// 	"Message": "foo foo foo",
	// 	"When":    "in 1 seconds at " + t2 + " today.",
	// }
	// expectedResponse = T("app.reminder.response", responseParameters)
	// assert.Equal(t, response, expectedResponse)

	// request.Payload = "me foo in 24 hours"
	// request.Occurrences = model.Occurrences{}
	// response, err = th.App.ScheduleReminder(request)
	// if err != nil {
	// 	t.Fatal(UNABLE_TO_SCHEDULE_REMINDER)
	// }
	// t2 = time.Now().Add(time.Hour * time.Duration(24)).Format(time.Kitchen)

	// responseParameters = map[string]interface{}{
	// 	"Target":  "You",
	// 	"UseTo":   "",
	// 	"Message": "foo",
	// 	"When":    "in 24 hours at " + t2 + " tomorrow.",
	// }
	// expectedResponse = T("app.reminder.response", responseParameters)
	// assert.Equal(t, response, expectedResponse)

	// request.Payload = "me foo in 3 days"
	// request.Occurrences = model.Occurrences{}
	// response, err = th.App.ScheduleReminder(request)
	// if err != nil {
	// 	t.Fatal(UNABLE_TO_SCHEDULE_REMINDER)
	// }
	// t3 := time.Now().AddDate(0, 0, 3)
	// responseParameters = map[string]interface{}{
	// 	"Target":  "You",
	// 	"UseTo":   "",
	// 	"Message": "foo",
	// 	"When":    "in 3 days at " + t3.Format(time.Kitchen) + " " + t3.Weekday().String() + ", " + t3.Month().String() + " " + th.App.daySuffixFromInt(user, t3.Day()) + ".",
	// }
	// expectedResponse = T("app.reminder.response", responseParameters)
	// assert.Equal(t, response, expectedResponse)

	// request.Payload = "me foo at 2:04 pm"
	// request.Occurrences = model.Occurrences{}
	// response, err = th.App.ScheduleReminder(request)
	// if err != nil {
	// 	t.Fatal(UNABLE_TO_SCHEDULE_REMINDER)
	// }
	// responseParameters = map[string]interface{}{
	// 	"Target":  "You",
	// 	"UseTo":   "",
	// 	"Message": "foo",
	// 	"When":    "at 2:04PM tomorrow.",
	// }
	// expectedResponse = T("app.reminder.response", responseParameters)
	// responseParameters = map[string]interface{}{
	// 	"Target":  "You",
	// 	"UseTo":   "",
	// 	"Message": "foo",
	// 	"When":    "at 2:04PM today.",
	// }
	// expectedResponse2 := T("app.reminder.response", responseParameters)
	// assert.True(t, response == expectedResponse || response == expectedResponse2)

	// request.Payload = "me foo on monday at 12:30PM"
	// request.Occurrences = model.Occurrences{}
	// response, err = th.App.ScheduleReminder(request)
	// if err != nil {
	// 	t.Fatal(UNABLE_TO_SCHEDULE_REMINDER)
	// }
	// t3, _ = time.Parse(time.RFC3339, request.Occurrences[0].Occurrence)
	// responseParameters = map[string]interface{}{
	// 	"Target":  "You",
	// 	"UseTo":   "",
	// 	"Message": "foo",
	// 	"When":    "at 12:30PM Monday, " + t3.Month().String() + " " + th.App.daySuffixFromInt(user, t3.Day()) + ".",
	// }
	// expectedResponse = T("app.reminder.response", responseParameters)
	// assert.Equal(t, response, expectedResponse)

	// request.Payload = "me foo every wednesday at 12:30PM"
	// request.Occurrences = model.Occurrences{}
	// response, err = th.App.ScheduleReminder(request)
	// if err != nil {
	// 	t.Fatal(UNABLE_TO_SCHEDULE_REMINDER)
	// }
	// responseParameters = map[string]interface{}{
	// 	"Target":  "You",
	// 	"UseTo":   "",
	// 	"Message": "foo",
	// 	"When":    "at 12:30PM every Wednesday.",
	// }
	// expectedResponse = T("app.reminder.response", responseParameters)
	// assert.Equal(t, response, expectedResponse)

	// request.Payload = "me tuesday foo"
	// request.Occurrences = model.Occurrences{}
	// response, err = th.App.ScheduleReminder(request)
	// if err != nil {
	// 	t.Fatal(UNABLE_TO_SCHEDULE_REMINDER)
	// }
	// t3, _ = time.Parse(time.RFC3339, request.Occurrences[0].Occurrence)
	// responseParameters = map[string]interface{}{
	// 	"Target":  "You",
	// 	"UseTo":   "",
	// 	"Message": "foo",
	// 	"When":    "at 9:00AM Tuesday, " + t3.Month().String() + " " + th.App.daySuffixFromInt(user, t3.Day()) + ".",
	// }
	// expectedResponse = T("app.reminder.response", responseParameters)
	// assert.Equal(t, response, expectedResponse)

	// request.Payload = "me tomorrow foo"
	// request.Occurrences = model.Occurrences{}
	// response, err = th.App.ScheduleReminder(request)
	// if err != nil {
	// 	t.Fatal(UNABLE_TO_SCHEDULE_REMINDER)
	// }
	// t3, _ = time.Parse(time.RFC3339, request.Occurrences[0].Occurrence)
	// responseParameters = map[string]interface{}{
	// 	"Target":  "You",
	// 	"UseTo":   "",
	// 	"Message": "foo",
	// 	"When":    "at 9:00AM tomorrow.",
	// }
	// expectedResponse = T("app.reminder.response", responseParameters)
	// responseParameters = map[string]interface{}{
	// 	"Target":  "You",
	// 	"UseTo":   "",
	// 	"Message": "foo",
	// 	"When":    "at 9:00AM " + t3.Weekday().String() + ", " + t3.Month().String() + " " + th.App.daySuffixFromInt(user, t3.Day()) + ".",
	// }
	// expectedResponse2 = T("app.reminder.response", responseParameters)
	// assert.True(t, response == expectedResponse || response == expectedResponse2)

}
