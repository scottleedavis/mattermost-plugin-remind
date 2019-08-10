package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/model"
)

type Occurrence struct {
	Hostname string

	Id string

	Username string

	ReminderId string

	Repeat string

	Occurrence time.Time

	Snoozed time.Time
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

func (p *Plugin) deleteSnoozedOccurrence(occurrence Occurrence) {

	bytes, err := p.API.KVGet(string(fmt.Sprintf("%v", occurrence.Snoozed)))
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
		if o.Id != occurrence.Id {
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

func (p *Plugin) deleteOccurrence(occurrence Occurrence) {
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
		if o.Id != occurrence.Id {
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

				rUser, rErr := p.API.GetUserByUsername(request.Username)
				if rErr != nil {
					return rErr
				}
				target := request.Reminder.Target
				if len(target) > 0 {
					target = target[1:]
				}
				if tUser, tErr := p.API.GetUserByUsername(target); tErr != nil {
					return tErr
				} else {
					if rUser.Id != tUser.Id {
						return errors.New("repeating reminders for another user not permitted")
					}
				}

			}
		}

		hostname, _ := os.Hostname()
		occurrence := Occurrence{
			Hostname:   hostname,
			Id:         model.NewId(),
			Username:   request.Username,
			ReminderId: request.Reminder.Id,
			Repeat:     repeat,
			Occurrence: o,
			Snoozed:    p.emptyTime,
		}

		request.Reminder.Occurrences = append(request.Reminder.Occurrences, occurrence)

		p.upsertOccurrence(&occurrence)

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
	ro, roErr := json.Marshal(occurrences)
	if roErr != nil {
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

	when = strings.Trim(when, " ")
	whenSplit := strings.Split(when, " ")
	if len(whenSplit) < 2 {
		return []time.Time{}, errors.New("empty when split")
	}
	value := whenSplit[1]
	units := whenSplit[len(whenSplit)-1]
	if len(whenSplit) == 2 {
		if strings.HasSuffix(units, T("seconds")) {
			value = strings.Trim(units, T("seconds"))
			units = T("seconds")
		} else if strings.HasSuffix(units, T("second")) {
			value = strings.Trim(units, T("second"))
			units = T("second")
		} else if strings.HasSuffix(units, T("secs")) {
			value = strings.Trim(units, T("secs"))
			units = T("secs")
		} else if strings.HasSuffix(units, T("sec")) {
			value = strings.Trim(units, T("sec"))
			units = T("sec")
		} else if strings.HasSuffix(units, T("s")) {
			value = strings.Trim(units, T("s"))
			units = T("s")
		} else if strings.HasSuffix(units, T("minutes")) {
			value = strings.Trim(units, T("minutes"))
			units = T("minutes")
		} else if strings.HasSuffix(units, T("minute")) {
			value = strings.Trim(units, T("minute"))
			units = T("minute")
		} else if strings.HasSuffix(units, T("min")) {
			value = strings.Trim(units, T("min"))
			units = T("min")
		} else if strings.HasSuffix(units, T("hours")) {
			value = strings.Trim(units, T("hours"))
			units = T("hours")
		} else if strings.HasSuffix(units, T("hour")) {
			value = strings.Trim(units, T("hour"))
			units = T("hour")
		} else if strings.HasSuffix(units, T("hrs")) {
			value = strings.Trim(units, T("hrs"))
			units = T("hrs")
		} else if strings.HasSuffix(units, T("hr")) {
			value = strings.Trim(units, T("hr"))
			units = T("hr")
		} else if strings.HasSuffix(units, T("days")) {
			value = strings.Trim(units, T("days"))
			units = T("days")
		} else if strings.HasSuffix(units, T("day")) {
			value = strings.Trim(units, T("day"))
			units = T("day")
		} else if strings.HasSuffix(units, T("d")) {
			value = strings.Trim(units, T("d"))
			units = T("d")
		} else if strings.HasSuffix(units, T("weeks")) {
			value = strings.Trim(units, T("weeks"))
			units = T("weeks")
		} else if strings.HasSuffix(units, T("week")) {
			value = strings.Trim(units, T("week"))
			units = T("week")
		} else if strings.HasSuffix(units, T("wks")) {
			value = strings.Trim(units, T("wks"))
			units = T("wks")
		} else if strings.HasSuffix(units, T("wk")) {
			value = strings.Trim(units, T("wk"))
			units = T("wk")
		} else if strings.HasSuffix(units, T("months")) {
			value = strings.Trim(units, T("months"))
			units = T("months")
		} else if strings.HasSuffix(units, T("month")) {
			value = strings.Trim(units, T("month"))
			units = T("month")
		} else if strings.HasSuffix(units, T("m")) {
			value = strings.Trim(units, T("m"))
			units = T("m")
		} else if strings.HasSuffix(units, T("years")) {
			value = strings.Trim(units, T("years"))
			units = T("years")
		} else if strings.HasSuffix(units, T("year")) {
			value = strings.Trim(units, T("year"))
			units = T("year")
		} else if strings.HasSuffix(units, T("yr")) {
			value = strings.Trim(units, T("yr"))
			units = T("yr")
		} else if strings.HasSuffix(units, T("y")) {
			value = strings.Trim(units, T("y"))
			units = T("y")
		}
	}

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
		T("mins"),
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
		if len(dateTimeSplit) < 2 {
			return []time.Time{}, errors.New("empty date time split")
		}
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
			} else if len(normalizedWhen) == 4 {
				s = normalizedWhen[:len(normalizedWhen)-2]
				s2 = normalizedWhen[len(normalizedWhen)-2:]
			} else if len(normalizedWhen) > 4 {
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
	} else {
		timeTest := strings.Split(chronoDate, " ")
		if len(timeTest) > 1 {
			check := timeTest[len(timeTest)-1]
			if strings.Contains(check, ":") {
				chronoDate = strings.Join(timeTest[:len(timeTest)-1], " ")
				chronoTime = check
			}
		}
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

		todayWeekDayNum := int(time.Now().In(location).Weekday())
		weekDayNum := p.weekDayNumber(dateUnit, user)
		day := 0

		if weekDayNum < todayWeekDayNum {
			day = 7 - (todayWeekDayNum - weekDayNum)
		} else if weekDayNum == todayWeekDayNum {
			todayTime := time.Now().In(location)
			if wallClock.Hour() >= todayTime.Hour() && wallClock.Minute() > todayTime.Minute() {
				day = weekDayNum - todayWeekDayNum
			} else {
				day = 7 + (weekDayNum - todayWeekDayNum)
			}
		} else {
			day = weekDayNum - todayWeekDayNum
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

	if strings.ToLower(whenSplit[0]) == T("everyday") {
		whenSplit[0] = T("day")
		whenSplit = append([]string{T("every")}, whenSplit...)
	}
	chronoUnit := strings.ToLower(strings.Join(whenSplit[1:], " "))
	var everyOther bool
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

	if chronoDate == T("weekday") || chronoDate == T("weekdays") {
		chronoDate = T("monday") + "," + T("tuesday") + "," + T("wednesday") + "," + T("thursday") + "," + T("friday")
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

			day := 0

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

			todayTime := time.Now().In(location)
			if wallClock.Hour() >= todayTime.Hour() && wallClock.Minute() > todayTime.Minute() {
				if everyOther {
					day = 1
				} else {
					day = 0
				}
			} else {
				if everyOther {
					day = 2
				} else {
					day = 1
				}
			}

			nextDay := time.Now().In(location).AddDate(0, 0, day)
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

			todayWeekDayNum := int(time.Now().In(location).Weekday())
			weekDayNum := p.weekDayNumber(dateUnit, user)
			day := 0

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

			if weekDayNum < todayWeekDayNum {
				day = 7 - (todayWeekDayNum - weekDayNum)
			} else if weekDayNum == todayWeekDayNum {
				todayTime := time.Now().In(location)
				if wallClock.Hour() >= todayTime.Hour() && wallClock.Minute() > todayTime.Minute() {
					day = weekDayNum - todayWeekDayNum
				} else {
					day = 7 + (weekDayNum - todayWeekDayNum)
				}
			} else {
				day = weekDayNum - todayWeekDayNum
			}

			nextDay := time.Now().In(location).AddDate(0, 0, day)
			occurrence := wallClock.In(location).AddDate(nextDay.Year(), int(nextDay.Month())-1, nextDay.Day()-1)
			times = append(times, p.chooseClosest(user, &occurrence, false).UTC())
			break
		default:

			dateSplit := p.regSplit(dateUnit, "T|Z")
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
	chronoTime := "9:00AM"
	chronoDate := chronoUnit

	if strings.Contains(chronoUnit, T("at")) {
		dateTimeSplit := strings.Split(chronoUnit, " "+T("at")+" ")
		chronoDate = dateTimeSplit[0]
		if len(dateTimeSplit) > 1 {
			chronoTime = dateTimeSplit[1]
		}
	} else {
		dateTimeSplit := strings.Split(chronoDate, " ")
		if len(dateTimeSplit) == 2 {
			chronoDate = dateTimeSplit[0]
			chronoTime = dateTimeSplit[1]
		} else {
			_, ntErr := p.normalizeTime(chronoDate, user)
			if ntErr == nil {
				return p.at(T("at")+" "+chronoDate, user)
			}
		}
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
	location := p.location(user)
	now := time.Now().In(location)
	t, _ := time.Parse(time.RFC3339, occurrence)

	if strings.HasPrefix(when, T("in")) {

		endDate := ""
		if now.YearDay() == t.YearDay() {
			endDate = T("today")
		} else if now.YearDay() == t.YearDay()-1 {
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

		endDate := ""
		if now.YearDay() == t.YearDay() {
			endDate = T("today")
		} else if now.YearDay() == t.YearDay()-1 {
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

		endDate := ""
		if now.YearDay() == t.YearDay() {
			endDate = T("today")
		} else if now.YearDay() == t.YearDay()-1 {
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

	endDate := ""
	if now.YearDay() == t.YearDay() {
		endDate = T("today")
	} else if now.YearDay() == t.YearDay()-1 {
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
	now := time.Now().In(location)

	if dayInterval {
		if chosen.Before(now) {
			return chosen.In(location).Round(time.Second).Add(time.Hour * 24 * time.Duration(1))
		} else {
			return *chosen
		}
	} else {
		if chosen.Before(time.Now().In(location)) {
			if chosen.Add(time.Hour * 12 * time.Duration(1)).Before(now) {
				return chosen.In(location).Round(time.Second).Add(time.Hour * 24 * time.Duration(1))
			} else {
				return chosen.In(location).Round(time.Second).Add(time.Hour * 12 * time.Duration(1))
			}
		} else {
			return *chosen
		}
	}
}
