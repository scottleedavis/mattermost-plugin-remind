package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/model"
)

type Occurrence struct {
	Id string

	Username string

	ReminderId string

	Occurrence time.Time

	Snoozed time.Time

	Repeat string
}

func (p *Plugin) ClearScheduledOccurrence(reminder Reminder, occurrence Occurrence) {

	bytes, err := p.API.KVGet(string(fmt.Sprintf("%v", occurrence.Occurrence)))
	if err != nil {
		p.API.LogError("failed KVGet %s", err)
		return
	}

	var occurrences []Occurrence
	if roErr := json.Unmarshal(bytes, &occurrences); roErr != nil {
		return
	}

	var occurrencesDelta []Occurrence
	for _, o := range occurrences {
		if o.ReminderId != reminder.Id {
			occurrencesDelta = append(occurrencesDelta, occurrence)
		}
	}

	ro, oErr := json.Marshal(occurrencesDelta)
	if oErr != nil {
		p.API.LogError("failed to marshal reminderOccurrences %s", occurrence.Id)
		return
	}

	p.API.KVSet(string(fmt.Sprintf("%v", occurrence.Occurrence)), ro)

}

func (p *Plugin) CreateOccurrences(request *ReminderRequest) error {

	user, uErr := p.API.GetUserByUsername(request.Username)
	if uErr != nil {
		return uErr
	}
	_, locale := p.translation(user)

	switch locale {
	case "en":
		return p.createOccurrencesEN(request)
	default:
		return p.createOccurrencesEN(request)
	}

}

func (p *Plugin) createOccurrencesEN(request *ReminderRequest) error {

	user, uErr := p.API.GetUserByUsername(request.Username)
	if uErr != nil {
		return uErr
	}
	T, _ := p.translation(user)

	if strings.HasPrefix(request.Reminder.When, T("in")) {
		if occurrences, inErr := p.in(request.Reminder.When, user); inErr != nil {
			return inErr
		} else {
			return p.addOccurrences(request, occurrences)
		}
	}

	if strings.HasPrefix(request.Reminder.When, T("at")) {
		if occurrences, inErr := p.at(request.Reminder.When, user); inErr != nil {
			return inErr
		} else {
			return p.addOccurrences(request, occurrences)
		}
	}

	if strings.HasPrefix(request.Reminder.When, T("on")) {
		if occurrences, inErr := p.on(request.Reminder.When, user); inErr != nil {
			return inErr
		} else {
			return p.addOccurrences(request, occurrences)
		}
	}

	if strings.HasPrefix(request.Reminder.When, T("every")) {
		if occurrences, inErr := p.every(request.Reminder.When, user); inErr != nil {
			return inErr
		} else {
			return p.addOccurrences(request, occurrences)
		}
	}

	if occurrences, freeErr := p.freeForm(request.Reminder.When, user); freeErr != nil {
		return freeErr
	} else {
		return p.addOccurrences(request, occurrences)
	}

}

func (p *Plugin) addOccurrences(request *ReminderRequest, occurrences []time.Time) error {

	user, uErr := p.API.GetUserByUsername(request.Username)
	if uErr != nil {
		return uErr
	}
	T, _ := p.translation(user)

	for _, o := range occurrences {

		repeat := ""

		if p.isRepeating(request) {
			repeat = request.Reminder.When
			if strings.HasPrefix(request.Reminder.Target, "@") &&
				request.Reminder.Target != T("me") {

				rUser, _ := p.API.GetUserByUsername(request.Username)

				if tUser, tErr := p.API.GetUserByUsername(request.Reminder.Target[1:]); tErr != nil {
					return tErr
				} else {
					if rUser.Id != tUser.Id {
						return errors.New("repeating reminders for another user not permitted")
					}
				}

			}
		}

		occurrence := &Occurrence{
			Id:         model.NewId(),
			Username:   request.Username,
			ReminderId: request.Reminder.Id,
			Repeat:     repeat,
			Occurrence: o,
			Snoozed:    p.emptyTime,
		}

		request.Reminder.Occurrences = p.upsertOccurrence(occurrence)
	}

	return nil
}

func (p *Plugin) isRepeating(request *ReminderRequest) bool {

	user, uErr := p.API.GetUserByUsername(request.Username)
	if uErr != nil {
		p.API.LogError(uErr.Error())
		return false
	}
	T, _ := p.translation(user)

	return strings.Contains(request.Reminder.When, T("every")) ||
		strings.Contains(request.Reminder.When, T("sundays")) ||
		strings.Contains(request.Reminder.When, T("mondays")) ||
		strings.Contains(request.Reminder.When, T("tuesdays")) ||
		strings.Contains(request.Reminder.When, T("wednesdays")) ||
		strings.Contains(request.Reminder.When, T("thursdays")) ||
		strings.Contains(request.Reminder.When, T("fridays")) ||
		strings.Contains(request.Reminder.When, T("saturdays"))

}

func (p *Plugin) upsertOccurrence(occurrence *Occurrence) []Occurrence {

	bytes, err := p.API.KVGet(string(fmt.Sprintf("%v", occurrence.Occurrence)))
	if err != nil {
		p.API.LogError("failed KVGet %s", err)
		return nil
	}

	var occurrences []Occurrence
	roErr := json.Unmarshal(bytes, &occurrences)
	if roErr != nil {
		p.API.LogDebug("new occurrence " + string(fmt.Sprintf("%v", occurrence.Occurrence)))
	}

	occurrences = append(occurrences, *occurrence)
	ro, __ := json.Marshal(occurrences)

	if __ != nil {
		p.API.LogError("failed to marshal reminderOccurrences %s", occurrence.Id)
		return occurrences
	}

	p.API.KVSet(string(fmt.Sprintf("%v", occurrence.Occurrence)), ro)

	return occurrences

}

func (p *Plugin) upsertSnoozedOccurrence(occurrence *Occurrence) []Occurrence {

	bytes, err := p.API.KVGet(string(fmt.Sprintf("%v", occurrence.Snoozed)))
	if err != nil {
		p.API.LogError("failed KVGet %s", err)
		return nil
	}

	var occurrences []Occurrence
	roErr := json.Unmarshal(bytes, &occurrences)
	if roErr != nil {
		p.API.LogDebug("new snoozed occurrence " + string(fmt.Sprintf("%v", occurrence.Snoozed)))
	}

	occurrences = append(occurrences, *occurrence)
	ro, roErr := json.Marshal(occurrences)
	if roErr != nil {
		p.API.LogError("failed to marshal reminderOccurrences %s", occurrence.Id)
		return occurrences
	}

	p.API.KVSet(string(fmt.Sprintf("%v", occurrence.Snoozed)), ro)

	return occurrences

}

func (p *Plugin) in(when string, user *model.User) (times []time.Time, err error) {

	_, locale := p.translation(user)

	switch locale {
	case "en":
		return p.inEN(when, user)
	default:
		return p.inEN(when, user)
	}

}

func (p *Plugin) inEN(when string, user *model.User) (times []time.Time, err error) {

	T, _ := p.translation(user)

	whenSplit := strings.Split(when, " ")
	value := whenSplit[1]
	units := whenSplit[len(whenSplit)-1]

	switch units {
	case T("seconds"),
		T("second"),
		T("secs"),
		T("sec"),
		T("s"):

		i, e := strconv.Atoi(value)

		if e != nil {
			num, wErr := p.wordToNumber(value, user)
			if wErr != nil {
				p.API.LogError(fmt.Sprintf("%v", wErr))
				return []time.Time{}, wErr
			}
			i = num
		}
		times = append(times, time.Now().UTC().Round(time.Second).Add(time.Second*time.Duration(i)))

		return times, nil

	case T("minutes"),
		T("minute"),
		T("min"):

		i, e := strconv.Atoi(value)

		if e != nil {
			num, wErr := p.wordToNumber(value, user)
			if wErr != nil {
				p.API.LogError(fmt.Sprintf("%v", wErr))
				return []time.Time{}, wErr
			}
			i = num
		}

		times = append(times, time.Now().UTC().Round(time.Second).Add(time.Minute*time.Duration(i)))

		return times, nil

	case T("hours"),
		T("hour"),
		T("hrs"),
		T("hr"):

		i, e := strconv.Atoi(value)

		if e != nil {
			num, wErr := p.wordToNumber(value, user)
			if wErr != nil {
				p.API.LogError(fmt.Sprintf("%v", wErr))
				return []time.Time{}, wErr
			}
			i = num
		}

		times = append(times, time.Now().UTC().Round(time.Second).Add(time.Hour*time.Duration(i)))

		return times, nil

	case T("days"),
		T("day"),
		T("d"):

		i, e := strconv.Atoi(value)

		if e != nil {
			num, wErr := p.wordToNumber(value, user)
			if wErr != nil {
				p.API.LogError(fmt.Sprintf("%v", wErr))
				return []time.Time{}, wErr
			}
			i = num
		}

		times = append(times, time.Now().UTC().Round(time.Second).Add(time.Hour*24*time.Duration(i)))

		return times, nil

	case T("weeks"),
		T("week"),
		T("wks"),
		T("wk"):

		i, e := strconv.Atoi(value)

		if e != nil {
			num, wErr := p.wordToNumber(value, user)
			if wErr != nil {
				p.API.LogError(fmt.Sprintf("%v", wErr))
				return []time.Time{}, wErr
			}
			i = num
		}

		times = append(times, time.Now().UTC().Round(time.Second).Add(time.Hour*24*7*time.Duration(i)))

		return times, nil

	case T("months"),
		T("month"),
		T("m"):

		i, e := strconv.Atoi(value)

		if e != nil {
			num, wErr := p.wordToNumber(value, user)
			if wErr != nil {
				p.API.LogError(fmt.Sprintf("%v", wErr))
				return []time.Time{}, wErr
			}
			i = num
		}

		times = append(times, time.Now().UTC().Round(time.Second).Add(time.Hour*24*30*time.Duration(i)))

		return times, nil

	case T("years"),
		T("year"),
		T("yr"),
		T("y"):

		i, e := strconv.Atoi(value)

		if e != nil {
			num, wErr := p.wordToNumber(value, user)
			if wErr != nil {
				p.API.LogError(fmt.Sprintf("%v", wErr))
				return []time.Time{}, wErr
			}
			i = num
		}

		times = append(times, time.Now().UTC().Round(time.Second).Add(time.Hour*24*365*time.Duration(i)))

		return times, nil

	default:
		return nil, errors.New("could not format 'in'")
	}

}

func (p *Plugin) at(when string, user *model.User) (times []time.Time, err error) {

	_, locale := p.translation(user)

	switch locale {
	case "en":
		return p.atEN(when, user)
	default:
		return p.atEN(when, user)
	}

}

func (p *Plugin) atEN(when string, user *model.User) (times []time.Time, err error) {

	T, _ := p.translation(user)
	location := p.location(user)

	whenTrim := strings.Trim(when, " ")
	whenSplit := strings.Split(whenTrim, " ")
	normalizedWhen := strings.ToLower(whenSplit[1])

	if strings.Contains(when, T("every")) {

		dateTimeSplit := strings.Split(when, " "+T("every")+" ")
		return p.every(T("every")+" "+dateTimeSplit[1]+" "+dateTimeSplit[0], user)

	} else if len(whenSplit) >= 3 &&
		(strings.EqualFold(whenSplit[2], T("pm")) ||
			strings.EqualFold(whenSplit[2], T("am"))) {

		if !strings.Contains(normalizedWhen, ":") {
			if len(normalizedWhen) >= 3 {
				hrs := string(normalizedWhen[:len(normalizedWhen)-2])
				mins := string(normalizedWhen[len(normalizedWhen)-2:])
				normalizedWhen = hrs + ":" + mins
			} else {
				normalizedWhen = normalizedWhen + ":00"
			}
		}
		t, pErr := time.ParseInLocation(time.Kitchen, normalizedWhen+strings.ToUpper(whenSplit[2]), location)
		if pErr != nil {
			p.API.LogError(fmt.Sprintf("%v", pErr))
		}

		now := time.Now().In(location).Round(time.Hour * time.Duration(24))
		occurrence := t.AddDate(now.Year(), int(now.Month())-1, now.Day()-1)
		return []time.Time{p.chooseClosest(user, &occurrence, true).UTC()}, nil

	} else if strings.HasSuffix(normalizedWhen, T("pm")) ||
		strings.HasSuffix(normalizedWhen, T("am")) {

		if !strings.Contains(normalizedWhen, ":") {
			var s string
			var s2 string
			if len(normalizedWhen) == 3 {
				s = normalizedWhen[:len(normalizedWhen)-2]
				s2 = normalizedWhen[len(normalizedWhen)-2:]
			} else if len(normalizedWhen) >= 4 {
				s = normalizedWhen[:len(normalizedWhen)-4]
				s2 = normalizedWhen[len(normalizedWhen)-4:]
			}

			if len(s2) > 2 {
				normalizedWhen = s + ":" + s2
			} else {
				normalizedWhen = s + ":00" + s2
			}

		}
		t, pErr := time.ParseInLocation(time.Kitchen, strings.ToUpper(normalizedWhen), location)
		if pErr != nil {
			p.API.LogError(fmt.Sprintf("%v", pErr))
		}

		now := time.Now().In(location).Round(time.Hour * time.Duration(24))
		occurrence := t.AddDate(now.Year(), int(now.Month())-1, now.Day()-1)

		return []time.Time{p.chooseClosest(user, &occurrence, true).UTC()}, nil

	}

	switch normalizedWhen {

	case T("noon"):

		now := time.Now().In(location)

		noon, pErr := time.ParseInLocation(time.Kitchen, "12:00PM", location)
		if pErr != nil {
			p.API.LogError(fmt.Sprintf("%v", pErr))
			return []time.Time{}, pErr
		}

		noon = noon.AddDate(now.Year(), int(now.Month())-1, now.Day()-1)
		return []time.Time{p.chooseClosest(user, &noon, true).UTC()}, nil

	case T("midnight"):

		now := time.Now().In(location)

		midnight, pErr := time.ParseInLocation(time.Kitchen, "12:00AM", location)
		if pErr != nil {
			p.API.LogError(fmt.Sprintf("%v", pErr))
			return []time.Time{}, pErr
		}

		midnight = midnight.AddDate(now.Year(), int(now.Month())-1, now.Day()-1)
		return []time.Time{p.chooseClosest(user, &midnight, true).UTC()}, nil

	case T("one"),
		T("two"),
		T("three"),
		T("four"),
		T("five"),
		T("six"),
		T("seven"),
		T("eight"),
		T("nine"),
		T("ten"),
		T("eleven"),
		T("twelve"):

		now := time.Now().In(location)
		ampm := string(now.Format(time.Kitchen)[len(now.Format(time.Kitchen))-2:])

		num, wErr := p.wordToNumber(normalizedWhen, user)
		if wErr != nil {
			return []time.Time{}, wErr
		}

		wordTime, _ := time.ParseInLocation(time.Kitchen, strconv.Itoa(num)+":00"+ampm, location)
		wordTime = wordTime.AddDate(now.Year(), int(now.Month())-1, now.Day()-1)

		return []time.Time{p.chooseClosest(user, &wordTime, false).UTC()}, nil

	case T("0"),
		T("1"),
		T("2"),
		T("3"),
		T("4"),
		T("5"),
		T("6"),
		T("7"),
		T("8"),
		T("9"),
		T("10"),
		T("11"),
		T("12"):

		now := time.Now().In(location)
		ampm := string(now.Format(time.Kitchen)[len(now.Format(time.Kitchen))-2:])

		num, wErr := strconv.Atoi(normalizedWhen)
		if wErr != nil {
			return []time.Time{}, wErr
		}

		wordTime, _ := time.ParseInLocation(time.Kitchen, strconv.Itoa(num)+":00"+ampm, location)
		wordTime = wordTime.AddDate(now.Year(), int(now.Month())-1, now.Day()-1)

		return []time.Time{p.chooseClosest(user, &wordTime, false).UTC()}, nil

	default:

		if !strings.Contains(normalizedWhen, ":") && len(normalizedWhen) >= 3 {
			s := normalizedWhen[:len(normalizedWhen)-2]
			normalizedWhen = s + ":" + normalizedWhen[len(normalizedWhen)-2:]
		}

		timeSplit := strings.Split(normalizedWhen, ":")
		hr, _ := strconv.Atoi(timeSplit[0])
		ampm := T("am")
		dayInterval := false

		if hr > 11 {
			ampm = T("pm")
		}
		if hr > 12 {
			hr -= 12
			dayInterval = true
			timeSplit[0] = strconv.Itoa(hr)
			normalizedWhen = strings.Join(timeSplit, ":")
		}

		t, pErr := time.ParseInLocation(time.Kitchen, strings.ToUpper(normalizedWhen+ampm), location)
		if pErr != nil {
			return []time.Time{}, pErr
		}

		now := time.Now().In(location).Round(time.Hour * time.Duration(24))
		occurrence := t.In(location).AddDate(now.Year(), int(now.Month())-1, now.Day()-1)
		return []time.Time{p.chooseClosest(user, &occurrence, dayInterval).UTC()}, nil

	}
}

func (p *Plugin) on(when string, user *model.User) (times []time.Time, err error) {

	_, locale := p.translation(user)

	switch locale {
	case "en":
		return p.onEN(when, user)
	default:
		return p.onEN(when, user)
	}

}

func (p *Plugin) onEN(when string, user *model.User) (times []time.Time, err error) {

	T, _ := p.translation(user)
	location := p.location(user)

	whenTrim := strings.Trim(when, " ")
	whenSplit := strings.Split(whenTrim, " ")

	if len(whenSplit) < 2 {
		return []time.Time{}, errors.New("not enough arguments")
	}

	chronoUnit := strings.ToLower(strings.Join(whenSplit[1:], " "))
	dateTimeSplit := strings.Split(chronoUnit, " "+T("at")+" ")
	chronoDate := dateTimeSplit[0]
	chronoTime := "9:00AM"
	if len(dateTimeSplit) > 1 {
		chronoTime = dateTimeSplit[1]
	}

	dateUnit, ndErr := p.normalizeDate(chronoDate, user)
	if ndErr != nil {
		return []time.Time{}, ndErr
	}

	timeUnit, ntErr := p.normalizeTime(chronoTime, user)
	if ntErr != nil {
		return []time.Time{}, ntErr
	}

	switch dateUnit {
	case T("sunday"),
		T("monday"),
		T("tuesday"),
		T("wednesday"),
		T("thursday"),
		T("friday"),
		T("saturday"):

		todayWeekDayNum := int(time.Now().Weekday())
		weekDayNum := p.weekDayNumber(dateUnit, user)
		day := 0

		if weekDayNum < todayWeekDayNum {
			day = 7 - (todayWeekDayNum - weekDayNum)
		} else if weekDayNum >= todayWeekDayNum {
			day = 7 + (weekDayNum - todayWeekDayNum)
		}

		timeUnitSplit := strings.Split(timeUnit, ":")
		hr, _ := strconv.Atoi(timeUnitSplit[0])
		ampm := strings.ToUpper(T("am"))

		if hr > 11 {
			ampm = strings.ToUpper(T("pm"))
		}
		if hr > 12 {
			hr -= 12
			timeUnitSplit[0] = strconv.Itoa(hr)
		}

		timeUnit = timeUnitSplit[0] + ":" + timeUnitSplit[1] + ampm
		wallClock, pErr := time.ParseInLocation(time.Kitchen, timeUnit, location)
		if pErr != nil {
			return []time.Time{}, pErr
		}

		nextDay := time.Now().In(location).AddDate(0, 0, day)
		occurrence := wallClock.In(location).AddDate(nextDay.Year(), int(nextDay.Month())-1, nextDay.Day()-1)

		return []time.Time{p.chooseClosest(user, &occurrence, false).UTC()}, nil

	case T("mondays"),
		T("tuesdays"),
		T("wednesdays"),
		T("thursdays"),
		T("fridays"),
		T("saturdays"),
		T("sundays"):

		return p.every(
			T("every")+" "+
				dateUnit[:len(dateUnit)-1]+" "+
				T("at")+" "+
				timeUnit[:len(timeUnit)-3],
			user)

	}

	dateSplit := p.regSplit(dateUnit, "T|Z")
	p.API.LogInfo("parsing " + dateSplit[0] + "T" + timeUnit + "Z")

	dateSplit = p.regSplit(dateSplit[0], "-")
	timeSplit := p.regSplit(timeUnit, ":")
	year, _ := strconv.Atoi(dateSplit[0])
	month, _ := strconv.Atoi(dateSplit[1])
	day, _ := strconv.Atoi(dateSplit[2])
	hour, _ := strconv.Atoi(timeSplit[0])
	minute, _ := strconv.Atoi(timeSplit[1])
	second, _ := strconv.Atoi(timeSplit[2])
	t := time.Date(
		year,
		time.Month(month),
		day,
		hour,
		minute,
		second, 0, location)

	return []time.Time{t.UTC()}, nil

}

func (p *Plugin) every(when string, user *model.User) (times []time.Time, err error) {

	_, locale := p.translation(user)

	switch locale {
	case "en":
		return p.everyEN(when, user)
	default:
		return p.everyEN(when, user)
	}

}

func (p *Plugin) everyEN(when string, user *model.User) (times []time.Time, err error) {

	T, _ := p.translation(user)
	location := p.location(user)

	whenTrim := strings.Trim(when, " ")
	whenSplit := strings.Split(whenTrim, " ")

	if len(whenSplit) < 2 {
		return []time.Time{}, errors.New("not enough arguments")
	}

	var everyOther bool
	chronoUnit := strings.ToLower(strings.Join(whenSplit[1:], " "))
	otherSplit := strings.Split(chronoUnit, T("other"))
	if len(otherSplit) == 2 {
		chronoUnit = strings.Trim(otherSplit[1], " ")
		everyOther = true
	}
	dateTimeSplit := strings.Split(chronoUnit, " "+T("at")+" ")
	chronoDate := dateTimeSplit[0]
	chronoTime := "9:00AM"
	if len(dateTimeSplit) > 1 {
		chronoTime = strings.Trim(dateTimeSplit[1], " ")
	}

	days := p.regSplit(chronoDate, "("+T("and")+")|(,)")

	for _, chrono := range days {

		dateUnit, ndErr := p.normalizeDate(strings.Trim(chrono, " "), user)
		if ndErr != nil {
			return []time.Time{}, ndErr
		}

		timeUnit, ntErr := p.normalizeTime(chronoTime, user)
		if ntErr != nil {
			return []time.Time{}, ntErr
		}

		switch dateUnit {
		case T("day"):
			d := 1
			if everyOther {
				d = 2
			}

			timeUnitSplit := strings.Split(timeUnit, ":")
			hr, _ := strconv.Atoi(timeUnitSplit[0])
			ampm := strings.ToUpper(T("am"))

			if hr > 11 {
				ampm = strings.ToUpper(T("pm"))
			}
			if hr > 12 {
				hr -= 12
				timeUnitSplit[0] = strconv.Itoa(hr)
			}

			timeUnit = timeUnitSplit[0] + ":" + timeUnitSplit[1] + ampm
			wallClock, pErr := time.ParseInLocation(time.Kitchen, timeUnit, location)
			if pErr != nil {
				return []time.Time{}, pErr
			}

			nextDay := time.Now().In(location).AddDate(0, 0, d)
			occurrence := wallClock.In(location).AddDate(nextDay.Year(), int(nextDay.Month())-1, nextDay.Day()-1)
			times = append(times, p.chooseClosest(user, &occurrence, false).UTC())

			break
		case T("sunday"),
			T("monday"),
			T("tuesday"),
			T("wednesday"),
			T("thursday"),
			T("friday"),
			T("saturday"):

			todayWeekDayNum := int(time.Now().Weekday())
			weekDayNum := p.weekDayNumber(dateUnit, user)
			day := 0

			if weekDayNum < todayWeekDayNum {
				day = 7 - (todayWeekDayNum - weekDayNum)
			} else if weekDayNum >= todayWeekDayNum {
				day = 7 + (weekDayNum - todayWeekDayNum)
			}

			timeUnitSplit := strings.Split(timeUnit, ":")
			hr, _ := strconv.Atoi(timeUnitSplit[0])
			ampm := strings.ToUpper(T("am"))

			if hr > 11 {
				ampm = strings.ToUpper(T("pm"))
			}
			if hr > 12 {
				hr -= 12
				timeUnitSplit[0] = strconv.Itoa(hr)
			}

			timeUnit = timeUnitSplit[0] + ":" + timeUnitSplit[1] + ampm
			wallClock, pErr := time.ParseInLocation(time.Kitchen, timeUnit, location)
			if pErr != nil {
				return []time.Time{}, pErr
			}

			nextDay := time.Now().In(location).AddDate(0, 0, day)
			occurrence := wallClock.In(location).AddDate(nextDay.Year(), int(nextDay.Month())-1, nextDay.Day()-1)
			times = append(times, p.chooseClosest(user, &occurrence, false).UTC())
			break
		default:

			dateSplit := p.regSplit(dateUnit, "T|Z")

			p.API.LogInfo("parsing " + dateSplit[0] + "T" + timeUnit + "Z")

			dateSplit = p.regSplit(dateSplit[0], "-")
			timeSplit := p.regSplit(timeUnit, ":")
			year, _ := strconv.Atoi(dateSplit[0])
			month, _ := strconv.Atoi(dateSplit[1])
			day, _ := strconv.Atoi(dateSplit[2])
			hour, _ := strconv.Atoi(timeSplit[0])
			minute, _ := strconv.Atoi(timeSplit[1])
			second, _ := strconv.Atoi(timeSplit[2])
			t := time.Date(
				year,
				time.Month(month),
				day,
				hour,
				minute,
				second, 0, location)

			times = append(times, t.UTC())

		}

	}

	return times, nil

}

func (p *Plugin) freeForm(when string, user *model.User) (times []time.Time, err error) {

	_, locale := p.translation(user)

	switch locale {
	case "en":
		return p.freeFormEN(when, user)
	default:
		return p.freeFormEN(when, user)
	}

}

func (p *Plugin) freeFormEN(when string, user *model.User) (times []time.Time, err error) {

	T, _ := p.translation(user)
	location := p.location(user)

	whenTrim := strings.Trim(when, " ")
	chronoUnit := strings.ToLower(whenTrim)
	dateTimeSplit := strings.Split(chronoUnit, " "+T("at")+" ")
	chronoTime := "9:00AM"
	chronoDate := dateTimeSplit[0]

	if len(dateTimeSplit) > 1 {
		chronoTime = dateTimeSplit[1]
	}
	dateUnit, ndErr := p.normalizeDate(chronoDate, user)
	if ndErr != nil {
		return []time.Time{}, ndErr
	}
	timeUnit, ntErr := p.normalizeTime(chronoTime, user)
	if ntErr != nil {
		return []time.Time{}, ntErr
	}
	timeUnit = chronoTime

	switch dateUnit {
	case T("today"):
		return p.at(T("at")+" "+timeUnit, user)
	case T("tomorrow"):
		return p.on(
			T("on")+" "+
				time.Now().In(location).Add(time.Hour*24).Weekday().String()+" "+
				T("at")+" "+
				timeUnit,
			user)
	case T("everyday"):
		return p.every(
			T("every")+" "+
				T("day")+" "+
				T("at")+" "+
				timeUnit,
			user)
	case T("mondays"),
		T("tuesdays"),
		T("wednesdays"),
		T("thursdays"),
		T("fridays"),
		T("saturdays"),
		T("sundays"):
		return p.every(
			T("every")+" "+
				dateUnit[:len(dateUnit)-1]+" "+
				T("at")+" "+
				timeUnit,
			user)
	case T("monday"),
		T("tuesday"),
		T("wednesday"),
		T("thursday"),
		T("friday"),
		T("saturday"),
		T("sunday"):
		return p.on(
			T("on")+" "+
				dateUnit+" "+
				T("at")+" "+
				timeUnit,
			user)
	default:
		return p.on(
			T("on")+" "+
				dateUnit[:len(dateUnit)-1]+" "+
				T("at")+" "+
				timeUnit,
			user)
	}

}

func (p *Plugin) formatWhen(username string, when string, occurrence string, snoozed bool) string {

	user, uErr := p.API.GetUserByUsername(username)
	if uErr != nil {
		return ""
	}
	_, locale := p.translation(user)

	switch locale {
	case "en":
		return p.formatWhenEN(username, when, occurrence, snoozed)
	default:
		return p.formatWhenEN(username, when, occurrence, snoozed)
	}
}

func (p *Plugin) formatWhenEN(username string, when string, occurrence string, snoozed bool) string {

	user, uErr := p.API.GetUserByUsername(username)
	if uErr != nil {
		p.API.LogError(uErr.Error())
		return ""
	}
	T, _ := p.translation(user)

	if strings.HasPrefix(when, T("in")) {

		t, _ := time.Parse(time.RFC3339, occurrence)
		endDate := ""
		if time.Now().YearDay() == t.YearDay() {
			endDate = T("today")
		} else if time.Now().YearDay() == t.YearDay()-1 {
			endDate = T("tomorrow")
		} else {
			endDate = t.Weekday().String() + ", " + t.Month().String() + " " + p.daySuffixFromInt(user, t.Day())
		}
		prefix := ""
		if !snoozed {
			prefix = when + " " + T("at") + " "
		}
		return prefix + t.Format(time.Kitchen) + " " + endDate + "."
	}

	if strings.HasPrefix(when, T("at")) {

		t, _ := time.Parse(time.RFC3339, occurrence)
		endDate := ""
		if time.Now().YearDay() == t.YearDay() {
			endDate = T("today")
		} else if time.Now().YearDay() == t.YearDay()-1 {
			endDate = T("tomorrow")
		} else {
			endDate = t.Weekday().String() + ", " + t.Month().String() + " " + p.daySuffixFromInt(user, t.Day())
		}
		prefix := ""
		if !snoozed {
			prefix = T("at") + " "
		}
		return prefix + t.Format(time.Kitchen) + " " + endDate + "."

	}

	if strings.HasPrefix(when, T("on")) {

		t, _ := time.Parse(time.RFC3339, occurrence)
		endDate := ""
		if time.Now().YearDay() == t.YearDay() {
			endDate = T("today")
		} else if time.Now().YearDay() == t.YearDay()-1 {
			endDate = T("tomorrow")
		} else {
			endDate = t.Weekday().String() + ", " + t.Month().String() + " " + p.daySuffixFromInt(user, t.Day())
		}
		prefix := ""
		if !snoozed {
			prefix = T("at") + " "
		}
		return prefix + t.Format(time.Kitchen) + " " + endDate + "."

	}

	if strings.HasPrefix(when, T("every")) {

		t, _ := time.Parse(time.RFC3339, occurrence)
		repeatDate := strings.Trim(strings.Split(when, T("at"))[0], " ")
		repeatDate = strings.Replace(repeatDate, T("every"), "", -1)
		repeatDate = strings.Title(strings.ToLower(repeatDate))
		repeatDate = T("every") + repeatDate
		prefix := ""
		if !snoozed {
			prefix = T("at") + " "
		}
		return prefix + t.Format(time.Kitchen) + " " + repeatDate + "."

	}

	t, _ := time.Parse(time.RFC3339, occurrence)
	endDate := ""
	if time.Now().YearDay() == t.YearDay() {
		endDate = T("today")
	} else if time.Now().YearDay() == t.YearDay()-1 {
		endDate = T("tomorrow")
	} else {
		endDate = t.Weekday().String() + ", " + t.Month().String() + " " + p.daySuffixFromInt(user, t.Day())
	}
	prefix := ""
	if !snoozed {
		prefix = T("at") + " "
	}
	return prefix + t.Format(time.Kitchen) + " " + endDate + "."

}

func (p *Plugin) chooseClosest(user *model.User, chosen *time.Time, dayInterval bool) time.Time {

	location := p.location(user)

	if dayInterval {
		if chosen.Before(time.Now().In(location)) {
			return chosen.In(location).Round(time.Second).Add(time.Hour * 24 * time.Duration(1))
		} else {
			return *chosen
		}
	} else {
		if chosen.Before(time.Now().In(location)) {
			if chosen.Add(time.Hour * 12 * time.Duration(1)).Before(time.Now()) {
				return chosen.In(location).Round(time.Second).Add(time.Hour * 24 * time.Duration(1))
			} else {
				return chosen.In(location).Round(time.Second).Add(time.Hour * 12 * time.Duration(1))
			}
		} else {
			return *chosen
		}
	}
}
