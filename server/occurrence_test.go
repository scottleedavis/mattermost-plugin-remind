package main

import (
	"testing"
	// "github.com/mattermost/mattermost-server/model"
	// "github.com/mattermost/mattermost-server/plugin/plugintest"
	// "github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/mock"
)

func TestIn(t *testing.T) {
	// th := Setup()
	// defer th.TearDown()

	// th.Server.InitReminders()
	// defer th.Server.StopReminders()
	// user, _ := th.App.GetUserByUsername(model.REMIND_BOTNAME)
	// //T := utils.GetUserTranslations(user.Locale)

	// when := "in one second"
	// times, iErr := th.App.in(when, user)
	// if iErr != nil {
	// 	t.Fatal("in one second doesn't parse")
	// }
	// var duration time.Duration
	// duration = times[0].Round(time.Second).Sub(time.Now().Round(time.Second))
	// assert.Equal(t, duration, time.Second)

	// when = "in 712 minutes"
	// times, iErr = th.App.in(when, user)
	// if iErr != nil {
	// 	t.Fatal("in 712 minutes doesn't parse")
	// }
	// duration = times[0].Round(time.Second).Sub(time.Now().Round(time.Second))
	// assert.Equal(t, duration, time.Minute*time.Duration(712))

	// when = "in three hours"
	// times, iErr = th.App.in(when, user)
	// if iErr != nil {
	// 	t.Fatal("in three hours doesn't parse")
	// }
	// duration = times[0].Round(time.Second).Sub(time.Now().Round(time.Second))
	// assert.Equal(t, duration, time.Hour*time.Duration(3))

	// when = "in 24 hours"
	// times, iErr = th.App.in(when, user)
	// if iErr != nil {
	// 	t.Fatal("in 24 hours doesn't parse")
	// }
	// duration = times[0].Round(time.Second).Sub(time.Now().Round(time.Second))
	// assert.Equal(t, duration, time.Hour*time.Duration(24))

	// when = "in 2 days"
	// times, iErr = th.App.in(when, user)
	// if iErr != nil {
	// 	t.Fatal("in 2 days doesn't parse")
	// }
	// duration = times[0].Round(time.Second).Sub(time.Now().Round(time.Second))
	// assert.Equal(t, duration, time.Hour*time.Duration(24)*time.Duration(2))

	// when = "in 90 weeks"
	// times, iErr = th.App.in(when, user)
	// if iErr != nil {
	// 	t.Fatal("in 90 weeks doesn't parse")
	// }
	// duration = times[0].Round(time.Second).Sub(time.Now().Round(time.Second))
	// assert.Equal(t, duration, time.Hour*time.Duration(24)*time.Duration(7)*time.Duration(90))

	// when = "in 4 months"
	// times, iErr = th.App.in(when, user)
	// if iErr != nil {
	// 	t.Fatal("in 4 months doesn't parse")
	// }
	// duration = times[0].Round(time.Second).Sub(time.Now().Round(time.Second))
	// assert.Equal(t, duration, time.Hour*time.Duration(24)*time.Duration(30)*time.Duration(4))

	// when = "in one year"
	// times, iErr = th.App.in(when, user)
	// if iErr != nil {
	// 	t.Fatal("in one year doesn't parse")
	// }
	// duration = times[0].Round(time.Second).Sub(time.Now().Round(time.Second))
	// assert.Equal(t, duration, time.Hour*time.Duration(24)*time.Duration(365))

}

func TestAt(t *testing.T) {
	// th := Setup()
	// defer th.TearDown()

	// th.Server.InitReminders()
	// defer th.Server.StopReminders()
	// user, _ := th.App.GetUserByUsername(model.REMIND_BOTNAME)
	// //T := utils.GetUserTranslations(user.Locale)

	// when := "at noon"
	// times, iErr := th.App.at(when, user)
	// if iErr != nil {
	// 	t.Fatal("at noon doesn't parse")
	// }
	// assert.Equal(t, times[0].Hour(), 12)

	// when = "at midnight"
	// times, iErr = th.App.at(when, user)
	// if iErr != nil {
	// 	t.Fatal("at midnight doesn't parse")
	// }
	// assert.Equal(t, times[0].Hour(), 0)

	// when = "at two"
	// times, iErr = th.App.at(when, user)
	// if iErr != nil {
	// 	t.Fatal("at two doesn't parse")
	// }
	// assert.True(t, times[0].Hour() == 2 || times[0].Hour() == 14)

	// when = "at 7"
	// times, iErr = th.App.at(when, user)
	// if iErr != nil {
	// 	t.Fatal("at 7 doesn't parse")
	// }
	// assert.True(t, times[0].Hour() == 7 || times[0].Hour() == 19)

	// when = "at 12:30pm"
	// times, iErr = th.App.at(when, user)
	// if iErr != nil {
	// 	t.Fatal("at 12:30pm doesn't parse")
	// }
	// assert.True(t, times[0].Hour() == 12 && times[0].Minute() == 30)

	// when = "at 7:12 pm"
	// times, iErr = th.App.at(when, user)
	// if iErr != nil {
	// 	t.Fatal("at 7:12 pm doesn't parse")
	// }
	// assert.True(t, times[0].Hour() == 19 && times[0].Minute() == 12)

	// when = "at 8:05 PM"
	// times, iErr = th.App.at(when, user)
	// if iErr != nil {
	// 	t.Fatal("at 8:05 PM doesn't parse")
	// }
	// assert.True(t, times[0].Hour() == 20 && times[0].Minute() == 5)

	// when = "at 9:52 am"
	// times, iErr = th.App.at(when, user)
	// if iErr != nil {
	// 	t.Fatal("at 9:52 am doesn't parse")
	// }
	// assert.True(t, times[0].Hour() == 9 && times[0].Minute() == 52)

	// when = "at 9:12"
	// times, iErr = th.App.at(when, user)
	// if iErr != nil {
	// 	t.Fatal("at 9:12 doesn't parse")
	// }
	// assert.True(t, (times[0].Hour() == 9 || times[0].Hour() == 21) && times[0].Minute() == 12)

	// when = "at 17:15"
	// times, iErr = th.App.at(when, user)
	// if iErr != nil {
	// 	t.Fatal("at 17:15 doesn't parse")
	// }
	// assert.True(t, times[0].Hour() == 17 && times[0].Minute() == 15)

	// when = "at 930am"
	// times, iErr = th.App.at(when, user)
	// if iErr != nil {
	// 	t.Fatal("at 930am doesn't parse")
	// }
	// assert.True(t, times[0].Hour() == 9 && times[0].Minute() == 30)

	// when = "at 1230 am"
	// times, iErr = th.App.at(when, user)
	// if iErr != nil {
	// 	t.Fatal("at 1230 am doesn't parse")
	// }
	// assert.True(t, times[0].Hour() == 0 && times[0].Minute() == 30)

	// when = "at 5PM"
	// times, iErr = th.App.at(when, user)
	// if iErr != nil {
	// 	t.Fatal("at 5PM doesn't parse")
	// }
	// assert.True(t, times[0].Hour() == 17 && times[0].Minute() == 0)

	// when = "at 4 am"
	// times, iErr = th.App.at(when, user)
	// if iErr != nil {
	// 	t.Fatal("at 4 am doesn't parse")
	// }
	// assert.True(t, times[0].Hour() == 4 && times[0].Minute() == 0)

	// when = "at 1400"
	// times, iErr = th.App.at(when, user)
	// if iErr != nil {
	// 	t.Fatal("at 1400 doesn't parse")
	// }
	// assert.True(t, times[0].Hour() == 14 && times[0].Minute() == 0)

	// when = "at 11:00 every Thursday"
	// times, iErr = th.App.at(when, user)
	// if iErr != nil {
	// 	t.Fatal("at 11:00 every Thursday doesn't parse")
	// }
	// assert.True(t, (times[0].Hour() == 11 || times[0].Hour() == 23) && times[0].Weekday().String() == "Thursday")

	// //TODO fix this test
	// //when = "at 3pm every day"
	// //times, iErr = th.App.at(when, user)
	// //if iErr != nil {
	// //	t.Fatal("at 3pm every day doesn't parse")
	// //}
	// //if times[0].Hour() != 15 {
	// //	t.Fatal("at 3pm every day isn't correct")
	// //}

}

func TestOn(t *testing.T) {
	// th := Setup()
	// defer th.TearDown()

	// th.Server.InitReminders()
	// defer th.Server.StopReminders()
	// user, _ := th.App.GetUserByUsername(model.REMIND_BOTNAME)
	// //T := utils.GetUserTranslations(user.Locale)

	// when := "on Monday"
	// times, iErr := th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on Monday doesn't parse")
	// }
	// assert.Equal(t, times[0].Weekday().String(), "Monday")

	// when = "on Tuesday"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on Tuesday doesn't parse")
	// }
	// assert.Equal(t, times[0].Weekday().String(), "Tuesday")

	// when = "on Wednesday"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on Wednesday doesn't parse")
	// }
	// assert.Equal(t, times[0].Weekday().String(), "Wednesday")

	// when = "on Thursday"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on Thursday doesn't parse")
	// }
	// assert.Equal(t, times[0].Weekday().String(), "Thursday")

	// when = "on Friday"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on Friday doesn't parse")
	// }
	// assert.Equal(t, times[0].Weekday().String(), "Friday")

	// when = "on Mondays"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on Mondays doesn't parse")
	// }
	// assert.Equal(t, times[0].Weekday().String(), "Monday")

	// when = "on Tuesdays at 11:15"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on Tuesdays at 11:15 doesn't parse")
	// }
	// assert.Equal(t, times[0].Weekday().String(), "Tuesday")

	// when = "on Wednesdays"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on Wednesdays doesn't parse")
	// }
	// assert.Equal(t, times[0].Weekday().String(), "Wednesday")

	// when = "on Thursdays"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on Thursdays doesn't parse")
	// }
	// assert.Equal(t, times[0].Weekday().String(), "Thursday")

	// when = "on Fridays"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on Fridays doesn't parse")
	// }
	// assert.Equal(t, times[0].Weekday().String(), "Friday")

	// when = "on mon"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on mon doesn't parse")
	// }
	// assert.Equal(t, times[0].Weekday().String(), "Monday")

	// when = "on wED"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on wED doesn't parse")
	// }
	// assert.Equal(t, times[0].Weekday().String(), "Wednesday")

	// when = "on tuesday at noon"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on tuesday at noon doesn't parse")
	// }
	// assert.True(t, times[0].Weekday().String() == "Tuesday" && times[0].Hour() == 12)

	// //TODO fix this test
	// //when = "on sunday at 3:42am"
	// //times, iErr = th.App.on(when, user)
	// //if iErr != nil {
	// //	mlog.Error(iErr.Error())
	// //	t.Fatal("on sunday at 3:42am doesn't parse")
	// //}
	// //assert.True(t, times[0].Weekday().String() == "Sunday" && times[0].Hour() == 3 && times[0].Minute() == 42)

	// when = "on December 15"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on December 15 doesn't parse")
	// }
	// assert.True(t, times[0].Month().String() == "December" && times[0].Day() == 15)

	// when = "on jan 12"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on jan 12 doesn't parse")
	// }
	// assert.True(t, times[0].Month().String() == "January" && times[0].Day() == 12)

	// when = "on July 12th"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on July 12th doesn't parse")
	// }
	// assert.True(t, times[0].Month().String() == "July" && times[0].Day() == 12)

	// when = "on March 22nd"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on March 22nd doesn't parse")
	// }
	// assert.True(t, times[0].Month().String() == "March" && times[0].Day() == 22)

	// when = "on March 17 at 5:41pm"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on March 17 at 5:41pm doesn't parse")
	// }
	// if times[0].Month().String() != "March" && times[0].Day() != 17 && times[0].Hour() != 17 && times[0].Minute() != 41 {
	// 	t.Fatal("on March 17 at 5:41pm isn't correct")
	// }
	// assert.True(t, times[0].Month().String() == "March" && times[0].Day() == 17 && times[0].Hour() == 17 && times[0].Minute() == 41)

	// when = "on September 7th 2020"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on September 7th 2019 doesn't parse")
	// }
	// assert.True(t, times[0].Month().String() == "September" && times[0].Day() == 7)

	// when = "on April 17 2020"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on April 17 2020 doesn't parse")
	// }
	// assert.True(t, times[0].Month().String() == "April" && times[0].Day() == 17)

	// //TODO fix this test
	// //when = "on April 9 2020 at 11am"
	// //times, iErr = th.App.on(when, user)
	// //if iErr != nil {
	// //	mlog.Error(iErr.Error())
	// //	t.Fatal("on April 9 2020 at 11am doesn't parse")
	// //}
	// //assert.True(t, times[0].Month().String() == "April" && times[0].Day() == 9 && times[0].Hour() == 11)

	// when = "on auguSt tenth 2019"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on auguSt tenth 2019 doesn't parse")
	// }
	// assert.True(t, times[0].Month().String() == "August" && times[0].Day() == 10)

	// when = "on 7"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on 7 doesn't parse")
	// }
	// assert.Equal(t, times[0].Day(), 7)

	// when = "on 7th"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on 7th doesn't parse")
	// }
	// assert.Equal(t, times[0].Day(), 7)

	// when = "on seven"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on seven doesn't parse")
	// }
	// assert.Equal(t, times[0].Day(), 7)

	// when = "on 1/17/20"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on 1/17/20 doesn't parse")
	// }
	// assert.True(t, times[0].Year() == 2020 && times[0].Month() == 1 && times[0].Day() == 17)

	// when = "on 12/17/2020"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on 12/17/2020 doesn't parse")
	// }
	// assert.True(t, times[0].Year() == 2020 && times[0].Month() == 12 && times[0].Day() == 17)

	// when = "on 17.1.20"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on 17.1.20 doesn't parse")
	// }
	// assert.True(t, times[0].Year() == 2020 && times[0].Month() == 1 && times[0].Day() == 17)

	// when = "on 17.12.2020"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on 17.12.2020 doesn't parse")
	// }
	// assert.True(t, times[0].Year() == 2020 && times[0].Month() == 12 && times[0].Day() == 17)

	// when = "on 12/1"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on 12/1 doesn't parse")
	// }
	// assert.True(t, times[0].Month() == 12 && times[0].Day() == 1)

	// when = "on 5-17-20"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on 5-17-20 doesn't parse")
	// }
	// assert.True(t, times[0].Month() == 5 && times[0].Day() == 17)

	// when = "on 12-5-2020"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on 12-5-2020 doesn't parse")
	// }
	// assert.True(t, times[0].Month() == 12 && times[0].Day() == 5)

	// when = "on 12-12"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on 12-12 doesn't parse")
	// }
	// assert.True(t, times[0].Month() == 12 && times[0].Day() == 12)

	// when = "on 1-1 at midnight"
	// times, iErr = th.App.on(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("on 1-1 at midnight doesn't parse")
	// }
	// assert.True(t, times[0].Month() == 1 && times[0].Day() == 1 && times[0].Hour() == 0)

}

func TestEvery(t *testing.T) {
	// th := Setup()
	// defer th.TearDown()

	// th.Server.InitReminders()
	// defer th.Server.StopReminders()
	// user, _ := th.App.GetUserByUsername(model.REMIND_BOTNAME)
	// //T := utils.GetUserTranslations(user.Locale)

	// when := "every Thursday"
	// times, iErr := th.App.every(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("every Thursday doesn't parse")
	// }
	// assert.Equal(t, times[0].Weekday().String(), "Thursday")

	// when = "every day"
	// times, iErr = th.App.every(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("every day doesn't parse")
	// }
	// assert.Equal(t, times[0].Weekday().String(), time.Now().AddDate(0, 0, 1).Weekday().String())

	// when = "every 12/18/2022"
	// times, iErr = th.App.every(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("every 12/18 doesn't parse")
	// }
	// assert.True(t, times[0].Month() == 12 && times[0].Year() == 2022)

	// when = "every January 25"
	// times, iErr = th.App.every(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("every January 25 doesn't parse")
	// }
	// assert.True(t, times[0].Month() == 1 && times[0].Day() == 25)

	// when = "every other Wednesday"
	// times, iErr = th.App.every(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("every other Wednesday doesn't parse")
	// }
	// assert.Equal(t, times[0].Weekday().String(), "Wednesday")

	// when = "every day at 11:32am"
	// times, iErr = th.App.every(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("every day at 11:32am doesn't parse")
	// }
	// assert.True(t, times[0].Hour() == 11 && times[0].Minute() == 32)

	// //TODO fix this test
	// //when = "every 5/5 at 7"
	// //times, iErr = th.App.every(when, user)
	// //if iErr != nil {
	// //	mlog.Error(iErr.Error())
	// //	t.Fatal("every 5/5 at 7 doesn't parse")
	// //}
	// //assert.True(t, times[0].Month() == 5 && times[0].Day() == 5 && (times[0].Hour() == 7 || times[0].Hour() == 19))

	// when = "every 7/20 at 1100"
	// times, iErr = th.App.every(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("every 7/20 at 1100 doesn't parse")
	// }
	// assert.True(t, times[0].Month() == 7 && times[0].Day() == 20 && (times[0].Hour() == 11 || times[0].Hour() == 23))

	// when = "every Monday at 7:32am"
	// times, iErr = th.App.every(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("every Monday at 7:32am doesn't parse")
	// }
	// assert.True(t, times[0].Weekday().String() == "Monday" && (times[0].Hour() == 7 || times[0].Hour() == 32))

	// when = "every monday and wednesday"
	// times, iErr = th.App.every(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("every monday and wednesday doesn't parse")
	// }
	// assert.True(t, times[0].Weekday().String() == "Monday" && times[1].Weekday().String() == "Wednesday")

	// when = "every wednesday, thursday"
	// times, iErr = th.App.every(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("every  wednesday, thursday doesn't parse")
	// }
	// assert.True(t, times[0].Weekday().String() == "Wednesday" && times[1].Weekday().String() == "Thursday")

	// when = "every other friday and saturday"
	// times, iErr = th.App.every(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("every  wednesday, thursday doesn't parse")
	// }
	// assert.True(t, times[0].Weekday().String() == "Friday" && times[1].Weekday().String() == "Saturday")

	// when = "every monday and wednesday at 1:39am"
	// times, iErr = th.App.every(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("every monday and wednesday at 1:39am doesn't parse")
	// }
	// assert.True(t, times[0].Weekday().String() == "Monday" && times[1].Weekday().String() == "Wednesday" && times[0].Hour() == 1 && times[0].Minute() == 39)

	// when = "every monday, tuesday and sunday at 11:00am"
	// times, iErr = th.App.every(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("every monday, tuesday and sunday at 11:00 doesn't parse")
	// }
	// assert.True(t, times[0].Weekday().String() == "Monday" && times[1].Weekday().String() == "Tuesday" && times[2].Weekday().String() == "Sunday" && times[0].Hour() == 11)

	// //TODO fix this test
	// //when = "every monday, tuesday at 2pm"
	// //times, iErr = th.App.every(when, user)
	// //if iErr != nil {
	// //	mlog.Error(iErr.Error())
	// //	t.Fatal("every monday, tuesday at 2pm doesn't parse")
	// //}
	// //assert.True(t, times[0].Weekday().String() == "Monday" && times[1].Weekday().String() == "Tuesday" && times[0].Hour() == 14)

	// when = "every 1/30 and 9/30 at noon"
	// times, iErr = th.App.every(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("every 1/30 and 9/30 at noon doesn't parse")
	// }
	// assert.True(t, times[0].Month() == 1 && times[0].Day() == 30 && times[1].Month() == 9 && times[1].Day() == 30 && times[0].Hour() == 12)

}

func TestFreeForm(t *testing.T) {
	// th := Setup()
	// defer th.TearDown()

	// th.Server.InitReminders()
	// defer th.Server.StopReminders()
	// user, _ := th.App.GetUserByUsername(model.REMIND_BOTNAME)
	// //T := utils.GetUserTranslations(user.Locale)

	// when := "monday"
	// times, iErr := th.App.freeForm(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("monday doesn't parse")
	// }
	// assert.Equal(t, times[0].Weekday().String(), "Monday")

	// when = "tuesday at 9:34pm"
	// times, iErr = th.App.freeForm(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("tuesday at 9:34pm doesn't parse")
	// }
	// assert.True(t, times[0].Weekday().String() == "Tuesday" && times[0].Hour() == 21 && times[0].Minute() == 34)

	// when = "wednesday"
	// times, iErr = th.App.freeForm(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("wednesday doesn't parse")
	// }
	// if times[0].Weekday().String() != "Wednesday" {
	// 	t.Fatal("wednesday isn't correct")
	// }
	// assert.Equal(t, times[0].Weekday().String(), "Wednesday")

	// when = "thursday at noon"
	// times, iErr = th.App.freeForm(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("thursday at noon doesn't parse")
	// }
	// assert.True(t, times[0].Weekday().String() == "Thursday" && times[0].Hour() == 12)

	// when = "friday"
	// times, iErr = th.App.freeForm(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("friday doesn't parse")
	// }
	// assert.Equal(t, times[0].Weekday().String(), "Friday")

	// when = "saturday"
	// times, iErr = th.App.freeForm(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("saturday doesn't parse")
	// }
	// assert.Equal(t, times[0].Weekday().String(), "Saturday")

	// when = "sunday at 4:20pm"
	// times, iErr = th.App.freeForm(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("sunday at 4:20pm doesn't parse")
	// }
	// assert.True(t, times[0].Weekday().String() == "Sunday" && times[0].Hour() == 16 && times[0].Minute() == 20)

	// when = "today at 3pm"
	// times, iErr = th.App.freeForm(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("today at 3pm doesn't parse")
	// }
	// assert.Equal(t, times[0].Hour(), 15)

	// when = "tomorrow"
	// times, iErr = th.App.freeForm(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("tomorrow doesn't parse")
	// }
	// assert.True(t, times[0].Weekday().String() == time.Now().AddDate(0, 0, 1).Weekday().String())

	// //TODO fix this test
	// //when = "tomorrow at 4pm"
	// //times, iErr = th.App.freeForm(when, user)
	// //if iErr != nil {
	// //	mlog.Error(iErr.Error())
	// //	t.Fatal("tomorrow at 4pm doesn't parse")
	// //}
	// //assert.True(t, times[0].Weekday().String() == time.Now().AddDate(0, 0, 1).Weekday().String() && times[0].Hour() == 16)

	// when = "everyday"
	// times, iErr = th.App.freeForm(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("everyday doesn't parse")
	// }
	// assert.Equal(t, times[0].Weekday().String(), time.Now().AddDate(0, 0, 1).Weekday().String())

	// //TODO fix this test
	// //when = "everyday at 3:23am"
	// //times, iErr = th.App.freeForm(when, user)
	// //if iErr != nil {
	// //	mlog.Error(iErr.Error())
	// //	t.Fatal("everyday at 3:23am doesn't parse")
	// //}
	// //assert.True(t, times[0].Weekday().String() == time.Now().AddDate(0, 0, 1).Weekday().String() && times[0].Hour() == 3 && times[0].Minute() == 23 )

	// when = "mondays"
	// times, iErr = th.App.freeForm(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("mondays doesn't parse")
	// }
	// assert.Equal(t, times[0].Weekday().String(), "Monday")

	// when = "tuesdays"
	// times, iErr = th.App.freeForm(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("tuesdays doesn't parse")
	// }
	// assert.Equal(t, times[0].Weekday().String(), "Tuesday")

	// when = "wednesdays"
	// times, iErr = th.App.freeForm(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("wednesdays doesn't parse")
	// }
	// assert.Equal(t, times[0].Weekday().String(), "Wednesday")

	// when = "thursdays"
	// times, iErr = th.App.freeForm(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("thursdays doesn't parse")
	// }
	// assert.Equal(t, times[0].Weekday().String(), "Thursday")

	// when = "fridays"
	// times, iErr = th.App.freeForm(when, user)
	// if iErr != nil {
	// 	mlog.Error(iErr.Error())
	// 	t.Fatal("fridays doesn't parse")
	// }
	// assert.Equal(t, times[0].Weekday().String(), "Friday")

}
