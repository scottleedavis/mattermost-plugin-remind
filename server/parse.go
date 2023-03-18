package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mattermost/mattermost-server/v6/model"
)

func (p *Plugin) ParseRequest(request *ReminderRequest) error {

	user, uErr := p.API.GetUserByUsername(request.Username)
	if uErr != nil {
		return uErr
	}
	T, _ := p.translation(user)

	commandSplit := strings.Split(request.Payload, " ")

	if strings.HasPrefix(request.Payload, T("me")) ||
		strings.HasPrefix(request.Payload, "~") ||
		strings.HasPrefix(request.Payload, "@") {

		request.Reminder.Target = commandSplit[0]

		firstIndex := strings.Index(request.Payload, "\"")
		lastIndex := strings.LastIndex(request.Payload, "\"")

		if firstIndex > -1 && lastIndex > -1 && firstIndex != lastIndex { // has quotes

			message := request.Payload[firstIndex : lastIndex+1]

			when := strings.Replace(request.Payload, message, "", -1)
			when = strings.Replace(when, request.Reminder.Target, "", 1)
			when = strings.Trim(when, " ")

			message = strings.Replace(message, "\"", "", -1)

			request.Reminder.When = when
			request.Reminder.Message = message
			return nil
		}

		if wErr := p.findWhen(request); wErr != nil {
			return wErr
		}

		toIndex := strings.Index(request.Reminder.When, T("to")+" ")
		if toIndex > -1 {
			request.Reminder.When = request.Reminder.When[0:toIndex]
		}

		message := strings.Replace(request.Payload, request.Reminder.When, "", -1)
		message = strings.Replace(message, request.Reminder.Target, "", 1)
		message = strings.Trim(message, " \"")

		if message == "" {
			return errors.New("no message parsed")
		}
		request.Reminder.Message = message

		return nil

	} else {
		request.Reminder.Target = T("me")

		firstIndex := strings.Index(request.Payload, "\"")
		lastIndex := strings.LastIndex(request.Payload, "\"")

		if firstIndex > -1 && lastIndex > -1 && firstIndex != lastIndex { // has quotes

			message := request.Payload[firstIndex : lastIndex+1]

			when := strings.Replace(request.Payload, message, "", -1)
			when = strings.Trim(when, " ")

			message = strings.Replace(message, "\"", "", -1)

			request.Reminder.When = when
			request.Reminder.Message = message
			return nil
		}

		if wErr := p.findWhen(request); wErr != nil {
			return wErr
		}

		toIndex := strings.Index(request.Reminder.When, T("to")+" ")
		if toIndex > -1 {
			request.Reminder.When = request.Reminder.When[0:toIndex]
		}

		message := strings.Replace(request.Payload, request.Reminder.When, "", -1)
		message = strings.Trim(message, " \"")

		if message == "" {
			return errors.New("no message parsed")
		}
		request.Reminder.Message = message

		return nil

	}

}

func (p *Plugin) findWhen(request *ReminderRequest) error {

	user, uErr := p.API.GetUserByUsername(request.Username)
	if uErr != nil {
		return uErr
	}
	_, locale := p.translation(user)

	switch locale {
	case "en":
		return p.findWhenEN(request)
	default:
		return p.findWhenEN(request)
	}

}

func (p *Plugin) findWhenEN(request *ReminderRequest) error {

	user, uErr := p.API.GetUserByUsername(request.Username)
	if uErr != nil {
		return uErr
	}
	T, _ := p.translation(user)

	inIndex := strings.Index(request.Payload, " "+T("in")+" ")
	if inIndex > -1 {
		request.Reminder.When = strings.Trim(request.Payload[inIndex:], " ")
		return nil
	}

	everyIndex := strings.Index(request.Payload, " "+T("every")+" ")
	atIndex := strings.Index(request.Payload, " "+T("at")+" ")
	if (everyIndex > -1 && atIndex == -1) || (atIndex > everyIndex) && everyIndex != -1 {
		request.Reminder.When = strings.Trim(request.Payload[everyIndex:], " ")
		return nil
	}

	onIndex := strings.Index(request.Payload, " "+T("on")+" ")
	if onIndex > -1 {
		request.Reminder.When = strings.Trim(request.Payload[onIndex:], " ")
		return nil
	}

	everydayIndex := strings.Index(request.Payload, " "+T("everyday")+" ")
	atIndex = strings.Index(request.Payload, " "+T("at")+" ")
	if (everydayIndex > -1 && atIndex >= -1) && (atIndex > everydayIndex) {
		request.Reminder.When = strings.Trim(request.Payload[everydayIndex:], " ")
		return nil
	}

	todayIndex := strings.Index(request.Payload, " "+T("today")+" ")
	atIndex = strings.Index(request.Payload, " "+T("at")+" ")
	if (todayIndex > -1 && atIndex >= -1) && (atIndex > todayIndex) {
		request.Reminder.When = strings.Trim(request.Payload[todayIndex:], " ")
		return nil
	}

	tomorrowIndex := strings.Index(request.Payload, " "+T("tomorrow")+" ")
	atIndex = strings.Index(request.Payload, " "+T("at")+" ")
	if (tomorrowIndex > -1 && atIndex >= -1) && (atIndex > tomorrowIndex) {
		request.Reminder.When = strings.Trim(request.Payload[tomorrowIndex:], " ")
		return nil
	}

	mondayIndex := strings.Index(request.Payload, " "+T("monday")+" ")
	atIndex = strings.Index(request.Payload, " "+T("at")+" ")
	if (mondayIndex > -1 && atIndex >= -1) && (atIndex > mondayIndex) {
		request.Reminder.When = strings.Trim(request.Payload[mondayIndex:], " ")
		return nil
	}

	tuesdayIndex := strings.Index(request.Payload, " "+T("tuesday")+" ")
	atIndex = strings.Index(request.Payload, " "+T("at")+" ")
	if (tuesdayIndex > -1 && atIndex >= -1) && (atIndex > tuesdayIndex) {
		request.Reminder.When = strings.Trim(request.Payload[tuesdayIndex:], " ")
		return nil
	}

	wednesdayIndex := strings.Index(request.Payload, " "+T("wednesday")+" ")
	atIndex = strings.Index(request.Payload, " "+T("at")+" ")
	if (wednesdayIndex > -1 && atIndex >= -1) && (atIndex > wednesdayIndex) {
		request.Reminder.When = strings.Trim(request.Payload[wednesdayIndex:], " ")
		return nil
	}

	thursdayIndex := strings.Index(request.Payload, " "+T("thursday")+" ")
	atIndex = strings.Index(request.Payload, " "+T("at")+" ")
	if (thursdayIndex > -1 && atIndex >= -1) && (atIndex > thursdayIndex) {
		request.Reminder.When = strings.Trim(request.Payload[thursdayIndex:], " ")
		return nil
	}

	fridayIndex := strings.Index(request.Payload, " "+T("friday")+" ")
	atIndex = strings.Index(request.Payload, " "+T("at")+" ")
	if (fridayIndex > -1 && atIndex >= -1) && (atIndex > fridayIndex) {
		request.Reminder.When = strings.Trim(request.Payload[fridayIndex:], " ")
		return nil
	}

	saturdayIndex := strings.Index(request.Payload, " "+T("saturday")+" ")
	atIndex = strings.Index(request.Payload, " "+T("at")+" ")
	if (saturdayIndex > -1 && atIndex >= -1) && (atIndex > saturdayIndex) {
		request.Reminder.When = strings.Trim(request.Payload[saturdayIndex:], " ")
		return nil
	}

	sundayIndex := strings.Index(request.Payload, " "+T("sunday")+" ")
	atIndex = strings.Index(request.Payload, " "+T("at")+" ")
	if (sundayIndex > -1 && atIndex >= -1) && (atIndex > sundayIndex) {
		request.Reminder.When = strings.Trim(request.Payload[sundayIndex:], " ")
		return nil
	}

	atIndex = strings.Index(request.Payload, " "+T("at")+" ")
	everyIndex = strings.Index(request.Payload, " "+T("every")+" ")
	if (atIndex > -1 && everyIndex >= -1) || (everyIndex > atIndex) && atIndex != -1 {
		request.Reminder.When = strings.Trim(request.Payload[atIndex:], " ")
		return nil
	}

	textSplit := strings.Split(request.Payload, " ")

	if len(textSplit) == 1 {
		request.Reminder.When = textSplit[0]
		return nil
	}

	lastWord := textSplit[len(textSplit)-2] + " " + textSplit[len(textSplit)-1]
	_, dErr := p.normalizeDate(lastWord, user)
	if dErr == nil {
		request.Reminder.When = lastWord
		return nil
	} else {

		lastWord = textSplit[len(textSplit)-1]

		switch lastWord {
		case T("tomorrow"):
			request.Reminder.When = lastWord
			return nil
		case T("everyday"),
			T("mondays"),
			T("tuesdays"),
			T("wednesdays"),
			T("thursdays"),
			T("fridays"),
			T("saturdays"),
			T("sundays"):
			request.Reminder.When = lastWord
		default:
			break
		}

		_, tErr := p.normalizeTime(lastWord, user)
		if tErr == nil {
			request.Reminder.When = lastWord
			return nil
		}

		_, dErr = p.normalizeDate(lastWord, user)
		if dErr == nil {
			request.Reminder.When = lastWord
			return nil
		} else {
			if len(textSplit) < 3 {
				return errors.New("unable to find when")
			}
			var firstWord string
			switch textSplit[1] {
			case T("at"):
				firstWord = textSplit[2]
				request.Reminder.When = textSplit[1] + " " + firstWord
				return nil
			case T("in"),
				T("on"):
				if len(textSplit) < 4 {
					return errors.New("unable to find when")
				}
				firstWord = textSplit[2] + " " + textSplit[3]
				request.Reminder.When = textSplit[1] + " " + firstWord
				return nil
			case T("tomorrow"),
				T("monday"),
				T("tuesday"),
				T("wednesday"),
				T("thursday"),
				T("friday"),
				T("saturday"),
				T("sunday"):
				firstWord = textSplit[1]
				request.Reminder.When = firstWord
				return nil
			default:
				break
			}
		}

	}

	return errors.New("unable to find when")

}

func (p *Plugin) normalizeTime(text string, user *model.User) (string, error) {

	T, _ := p.translation(user)
	location := p.location(user)

	switch text {
	case T("noon"):
		return "12:00:00", nil
	case T("midnight"):
		return "00:00:00", nil
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

		num, wErr := p.wordToNumber(text, user)
		if wErr != nil {
			p.API.LogError(fmt.Sprintf("%v", wErr))
			return "", wErr
		}

		wordTime := time.Now().In(location).Round(time.Hour).Add(time.Hour * time.Duration(num+2))

		dateTimeSplit := p.regSplit(p.chooseClosest(user, &wordTime, false).Format(time.RFC3339), "T|Z")

		switch len(dateTimeSplit) {
		case 2:
			tzSplit := strings.Split(dateTimeSplit[1], "-")
			return tzSplit[0], nil
		case 3:
			break
		default:
			return "", errors.New("unrecognized dateTime format")
		}

		return dateTimeSplit[1], nil

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
		T("12"),
		T("13"),
		T("14"),
		T("15"),
		T("16"),
		T("17"),
		T("18"),
		T("19"),
		T("20"),
		T("21"),
		T("22"),
		T("23"):

		if len(text) == 1 {
			text = "0" + text
		}

		return text + ":00:00", nil

	default:
		break
	}

	t := text
	if match, _ := regexp.MatchString("(1[012]|[1-9]):[0-5][0-9](\\s)?(?i)(am|pm)", t); match { // 12:30PM, 12:30 pm

		t = strings.ToUpper(strings.Replace(t, " ", "", -1))
		test, tErr := time.Parse(time.Kitchen, t)
		if tErr != nil {
			return "", tErr
		}

		dateTimeSplit := p.regSplit(test.Format(time.RFC3339), "T|Z")

		switch len(dateTimeSplit) {
		case 2:
			tzSplit := strings.Split(dateTimeSplit[1], "-")
			return tzSplit[0], nil
		case 3:
			break
		default:
			return "", errors.New("unrecognized dateTime format")
		}

		return dateTimeSplit[1], nil
	} else if match, _ := regexp.MatchString("(1[012]|[1-9]):[0-5][0-9]", t); match { // 12:30

		nowkit := time.Now().In(location).Format(time.Kitchen)
		ampm := string(nowkit[len(nowkit)-2:])
		timeUnitSplit := strings.Split(t, ":")
		hr, _ := strconv.Atoi(timeUnitSplit[0])

		if hr > 11 {
			ampm = strings.ToUpper(T("pm"))
		}
		if hr > 12 {
			hr -= 12
			timeUnitSplit[0] = strconv.Itoa(hr)
		}

		t = timeUnitSplit[0] + ":" + timeUnitSplit[1] + ampm

		test, tErr := time.ParseInLocation(time.Kitchen, t, location)
		if tErr != nil {
			return "", tErr
		}

		dateTimeSplit := p.regSplit(p.chooseClosest(user, &test, false).Format(time.RFC3339), "T|Z")

		switch len(dateTimeSplit) {
		case 2:
			tzSplit := strings.Split(dateTimeSplit[1], "-")
			return tzSplit[0], nil
		case 3:
			break
		default:
			return "", errors.New("unrecognized dateTime format")
		}

		return dateTimeSplit[1], nil

	} else if match, _ := regexp.MatchString("(1[012]|[1-9])(\\s)?(?i)(am|pm)", t); match { // 5PM, 7 am

		nowkit := time.Now().In(location).Format(time.Kitchen)
		ampm := string(nowkit[len(nowkit)-2:])

		if strings.HasSuffix(t, "pm") {
			ampm = "PM"
		} else if strings.HasSuffix(t, "am") {
			ampm = "AM"
		}

		timeSplit := p.regSplit(t, "(?i)(am|pm)")

		test, tErr := time.ParseInLocation(time.Kitchen, timeSplit[0]+":00"+ampm, location)
		if tErr != nil {
			return "", tErr
		}

		dateTimeSplit := p.regSplit(p.chooseClosest(user, &test, false).Format(time.RFC3339), "T|Z")

		switch len(dateTimeSplit) {
		case 2:
			tzSplit := strings.Split(dateTimeSplit[1], "-")
			return tzSplit[0], nil
		case 3:
			break
		default:
			return "", errors.New("unrecognized dateTime format")
		}

		return dateTimeSplit[1], nil
	} else if match, _ := regexp.MatchString("(1[012]|[1-9])[0-5][0-9]", t); match { // 1200

		return t[:len(t)-2] + ":" + t[len(t)-2:] + ":00", nil

	}

	return "", errors.New("unable to normalize time")
}

func (p *Plugin) normalizeDate(text string, user *model.User) (string, error) {

	location := p.location(user)
	T, _ := p.translation(user)

	date := strings.ToLower(text)
	if strings.EqualFold(T("day"), date) {
		return date, nil
	} else if strings.EqualFold(T("today"), date) {
		return date, nil
	} else if strings.EqualFold(T("everyday"), date) {
		return date, nil
	} else if strings.EqualFold(T("tomorrow"), date) {
		return date, nil
	}

	switch date {
	case T("mon"),
		T("monday"):
		return T("monday"), nil
	case T("tues"),
		T("tuesday"):
		return T("tuesday"), nil
	case T("wed"),
		T("wednes"),
		T("wednesday"):
		return T("wednesday"), nil
	case T("thur"),
		T("thursday"):
		return T("thursday"), nil
	case T("fri"),
		T("friday"):
		return T("friday"), nil
	case T("sat"),
		T("satur"),
		T("saturday"):
		return T("saturday"), nil
	case T("sun"),
		T("sunday"):
		return T("sunday"), nil
	case T("mondays"),
		T("tuesdays"),
		T("wednesdays"),
		T("thursdays"),
		T("fridays"),
		T("saturdays"),
		T("sundays"):
		return date, nil
	}

	if strings.Contains(date, T("jan")) ||
		strings.Contains(date, T("january")) ||
		strings.Contains(date, T("feb")) ||
		strings.Contains(date, T("february")) ||
		strings.Contains(date, T("mar")) ||
		strings.Contains(date, T("march")) ||
		strings.Contains(date, T("apr")) ||
		strings.Contains(date, T("april")) ||
		strings.Contains(date, T("may")) ||
		strings.Contains(date, T("june")) ||
		strings.Contains(date, T("july")) ||
		strings.Contains(date, T("aug")) ||
		strings.Contains(date, T("august")) ||
		strings.Contains(date, T("sept")) ||
		strings.Contains(date, T("september")) ||
		strings.Contains(date, T("oct")) ||
		strings.Contains(date, T("october")) ||
		strings.Contains(date, T("nov")) ||
		strings.Contains(date, T("november")) ||
		strings.Contains(date, T("dec")) ||
		strings.Contains(date, T("december")) {

		date = strings.Replace(date, ",", "", -1)
		parts := strings.Split(date, " ")

		switch len(parts) {
		case 1:
			break
		case 2:
			if len(parts[1]) > 2 {
				parts[1] = p.daySuffix(user, parts[1])
			}
			if _, err := strconv.Atoi(parts[1]); err != nil {
				if wn, wErr := p.wordToNumber(parts[1], user); wErr == nil {
					parts[1] = strconv.Itoa(wn)
				}
			}

			parts = append(parts, fmt.Sprintf("%v", time.Now().In(location).Year()))

			break
		case 3:
			if len(parts[1]) > 2 {
				parts[1] = p.daySuffix(user, parts[1])
			}

			if _, err := strconv.Atoi(parts[1]); err != nil {
				if wn, wErr := p.wordToNumber(parts[1], user); wErr == nil {
					parts[1] = strconv.Itoa(wn)
				} else {
					p.API.LogError(wErr.Error())
				}

				if _, pErr := strconv.Atoi(parts[2]); pErr != nil {
					return "", pErr
				}
			}

			break
		default:
			return "", errors.New("unrecognized date format")
		}

		var err error
		parts[0], err = p.monthNumber(parts[0], user)
		if err != nil {
			return "", err
		}

		if len(parts) < 3 {
			return "", errors.New("unrecognized date format")
		}

		mon, mErr := strconv.Atoi(parts[0])
		day, dErr := strconv.Atoi(parts[1])
		year, yErr := strconv.Atoi(parts[2])

		if mErr != nil {
			return "", mErr
		}
		if dErr != nil {
			return "", dErr
		}
		if yErr != nil {
			return "", yErr
		}
		timeNow := time.Now().In(location)
		parseTime := time.Date(year, time.Month(mon), day, 0, 0, 0, 0, location)
		if timeNow.After(parseTime) {
			parts[2] = fmt.Sprintf("%v", time.Now().In(location).Year()+1)
		}

		if len(parts[1]) < 2 {
			parts[1] = "0" + parts[1]
		}
		return parts[2] + "-" + parts[0] + "-" + parts[1] + "T00:00:00Z", nil

	} else if match, _ := regexp.MatchString("^([0-9]{4}-[0-9]{2}-[0-9]{2})", date); match {

		date := p.regSplit(date, "-")

		switch len(date) {
		case 3:
			year, yErr := strconv.Atoi(date[0])
			if yErr != nil {
				return "", yErr
			}
			month, mErr := strconv.Atoi(date[1])
			if mErr != nil {
				return "", mErr
			}
			day, dErr := strconv.Atoi(date[2])
			if dErr != nil {
				return "", dErr
			}

			return time.Date(year, time.Month(month), day, 0, 0, 0, 0, location).Format(time.RFC3339), nil

		default:
			return "", errors.New("unrecognized date")
		}

	} else if match, _ := regexp.MatchString("^(([0-9]{2}|[0-9]{1})(-|/)([0-9]{2}|[0-9]{1})((-|/)([0-9]{4}|[0-9]{2}))?)", date); match {

		date := p.regSplit(date, "-|/")

		switch len(date) {
		case 2:
			year := time.Now().In(location).Year()
			month, mErr := strconv.Atoi(date[0])
			if mErr != nil {
				return "", mErr
			}
			day, dErr := strconv.Atoi(date[1])
			if dErr != nil {
				return "", dErr
			}

			return time.Date(year, time.Month(month), day, 0, 0, 0, 0, location).Format(time.RFC3339), nil

		case 3:
			if len(date[2]) == 2 {
				date[2] = "20" + date[2]
			}
			year, yErr := strconv.Atoi(date[2])
			if yErr != nil {
				return "", yErr
			}
			month, mErr := strconv.Atoi(date[0])
			if mErr != nil {
				return "", mErr
			}
			day, dErr := strconv.Atoi(date[1])
			if dErr != nil {
				return "", dErr
			}

			return time.Date(year, time.Month(month), day, 0, 0, 0, 0, location).Format(time.RFC3339), nil

		default:
			return "", errors.New("unrecognized date")
		}

	} else if match, _ := regexp.MatchString("^(([0-9]{2}|[0-9]{1})(.)([0-9]{2}|[0-9]{1})((.)([0-9]{4}|[0-9]{2}))?)", date); match {

		date := p.regSplit(date, "\\.")

		switch len(date) {
		case 2:
			year := time.Now().In(location).Year()
			month, mErr := strconv.Atoi(date[1])
			if mErr != nil {
				return "", mErr
			}
			day, dErr := strconv.Atoi(date[0])
			if dErr != nil {
				return "", dErr
			}

			return time.Date(year, time.Month(month), day, 0, 0, 0, 0, location).Format(time.RFC3339), nil

		case 3:
			if len(date[2]) == 2 {
				date[2] = "20" + date[2]
			}
			year, yErr := strconv.Atoi(date[2])
			if yErr != nil {
				return "", yErr
			}
			month, mErr := strconv.Atoi(date[1])
			if mErr != nil {
				return "", mErr
			}
			day, dErr := strconv.Atoi(date[0])
			if dErr != nil {
				return "", dErr
			}

			return time.Date(year, time.Month(month), day, 0, 0, 0, 0, location).Format(time.RFC3339), nil

		default:
			return "", errors.New("unrecognized date")
		}

	} else { //single number day

		var dayInt int
		day := p.daySuffix(user, date)

		if d, nErr := strconv.Atoi(day); nErr != nil {
			if wordNum, wErr := p.wordToNumber(date, user); wErr != nil {
				return "", wErr
			} else {
				day = strconv.Itoa(wordNum)
				dayInt = wordNum
			}
		} else {
			dayInt = d
		}

		month := time.Now().In(location).Month()
		year := time.Now().In(location).Year()
		t := time.Date(year, month, dayInt, 0, 0, 0, 0, location)

		if t.Before(time.Now().In(location)) {
			t = t.AddDate(0, 1, 0)
		}

		return t.Format(time.RFC3339), nil

	}

}

func (p *Plugin) daySuffixFromInt(user *model.User, day int) string {

	T, _ := p.translation(user)

	daySuffixes := []string{
		T("0th"),
		T("1st"),
		T("2nd"),
		T("3rd"),
		T("4th"),
		T("5th"),
		T("6th"),
		T("7th"),
		T("8th"),
		T("9th"),
		T("10th"),
		T("11th"),
		T("12th"),
		T("13th"),
		T("14th"),
		T("15th"),
		T("16th"),
		T("17th"),
		T("18th"),
		T("19th"),
		T("20th"),
		T("21st"),
		T("22nd"),
		T("23rd"),
		T("24th"),
		T("25th"),
		T("26th"),
		T("27th"),
		T("28th"),
		T("29th"),
		T("30th"),
		T("31st"),
	}
	return daySuffixes[day]

}

func (p *Plugin) daySuffix(user *model.User, day string) string {

	T, _ := p.translation(user)

	daySuffixes := []string{
		T("0th"),
		T("1st"),
		T("2nd"),
		T("3rd"),
		T("4th"),
		T("5th"),
		T("6th"),
		T("7th"),
		T("8th"),
		T("9th"),
		T("10th"),
		T("11th"),
		T("12th"),
		T("13th"),
		T("14th"),
		T("15th"),
		T("16th"),
		T("17th"),
		T("18th"),
		T("19th"),
		T("20th"),
		T("21st"),
		T("22nd"),
		T("23rd"),
		T("24th"),
		T("25th"),
		T("26th"),
		T("27th"),
		T("28th"),
		T("29th"),
		T("30th"),
		T("31st"),
	}
	for _, suffix := range daySuffixes {
		if suffix == day {
			day = day[:len(day)-2]
			break
		}
	}
	return day
}

func (p *Plugin) monthNumber(month string, user *model.User) (string, error) {

	T, _ := p.translation(user)

	switch month {
	case T("jan"),
		T("january"):
		return "01", nil
	case T("feb"),
		T("february"):
		return "02", nil
	case T("mar"),
		T("march"):
		return "03", nil
	case T("apr"),
		T("april"):
		return "04", nil
	case T("may"):
		return "05", nil
	case T("june"):
		return "06", nil
	case T("july"):
		return "07", nil
	case T("aug"),
		T("august"):
		return "08", nil
	case T("sept"),
		T("september"):
		return "09", nil
	case T("oct"),
		T("october"):
		return "10", nil
	case T("nov"),
		T("november"):
		return "11", nil
	case T("dec"),
		T("december"):
		return "12", nil
	default:
		return "", errors.New("month not found")
	}
}

func (p *Plugin) weekDayNumber(day string, user *model.User) int {

	T, _ := p.translation(user)

	switch day {
	case T("sunday"):
		return 0
	case T("monday"):
		return 1
	case T("tuesday"):
		return 2
	case T("wednesday"):
		return 3
	case T("thursday"):
		return 4
	case T("friday"):
		return 5
	case T("saturday"):
		return 6
	default:
		return -1
	}
}

func (p *Plugin) regSplit(text string, delimeter string) []string {

	reg := regexp.MustCompile(delimeter)
	indexes := reg.FindAllStringIndex(text, -1)
	laststart := 0
	result := make([]string, len(indexes)+1)
	for i, element := range indexes {
		result[i] = text[laststart:element[0]]
		laststart = element[1]
	}
	result[len(indexes)] = text[laststart:]
	return result
}

func (p *Plugin) wordToNumber(word string, user *model.User) (int, error) {

	T, _ := p.translation(user)

	var sum int
	var temp int
	var previous int

	numbers := make(map[string]int)
	onumbers := make(map[string]int)
	tnumbers := make(map[string]int)

	numbers[T("zero")] = 0
	numbers[T("one")] = 1
	numbers[T("two")] = 2
	numbers[T("three")] = 3
	numbers[T("four")] = 4
	numbers[T("five")] = 5
	numbers[T("six")] = 6
	numbers[T("seven")] = 7
	numbers[T("eight")] = 8
	numbers[T("nine")] = 9
	numbers[T("ten")] = 10
	numbers[T("eleven")] = 11
	numbers[T("twelve")] = 12
	numbers[T("thirteen")] = 13
	numbers[T("fourteen")] = 14
	numbers[T("fifteen")] = 15
	numbers[T("sixteen")] = 16
	numbers[T("seventeen")] = 17
	numbers[T("eighteen")] = 18
	numbers[T("nineteen")] = 19

	tnumbers[T("twenty")] = 20
	tnumbers[T("thirty")] = 30
	tnumbers[T("forty")] = 40
	tnumbers[T("fifty")] = 50
	tnumbers[T("sixty")] = 60
	tnumbers[T("seventy")] = 70
	tnumbers[T("eighty")] = 80
	tnumbers[T("ninety")] = 90

	onumbers[T("hundred")] = 100
	onumbers[T("thousand")] = 1000
	onumbers[T("million")] = 1000000
	onumbers[T("billion")] = 1000000000

	numbers[T("first")] = 1
	numbers[T("second")] = 2
	numbers[T("third")] = 3
	numbers[T("fourth")] = 4
	numbers[T("fifth")] = 5
	numbers[T("sixth")] = 6
	numbers[T("seventh")] = 7
	numbers[T("eighth")] = 8
	numbers[T("nineth")] = 9
	numbers[T("tenth")] = 10
	numbers[T("eleventh")] = 11
	numbers[T("twelveth")] = 12
	numbers[T("thirteenth")] = 13
	numbers[T("fourteenth")] = 14
	numbers[T("fifteenth")] = 15
	numbers[T("sixteenth")] = 16
	numbers[T("seventeenth")] = 17
	numbers[T("eighteenth")] = 18
	numbers[T("nineteenth")] = 19

	tnumbers[T("twenteth")] = 20
	tnumbers[T("twentyfirst")] = 21
	tnumbers[T("twentysecond")] = 22
	tnumbers[T("twentythird")] = 23
	tnumbers[T("twentyfourth")] = 24
	tnumbers[T("twentyfifth")] = 25
	tnumbers[T("twentysixth")] = 26
	tnumbers[T("twentyseventh")] = 27
	tnumbers[T("twentyeight")] = 28
	tnumbers[T("twentynineth")] = 29
	tnumbers[T("thirteth")] = 30
	tnumbers[T("thirtyfirst")] = 31

	splitted := strings.Split(strings.ToLower(word), " ")

	for _, split := range splitted {
		if numbers[split] != 0 {
			temp = numbers[split]
			sum = sum + temp
			previous = previous + temp
		} else if onumbers[split] != 0 {
			if sum != 0 {
				sum = sum - previous
			}
			sum = sum + previous*onumbers[split]
			temp = 0
			previous = 0
		} else if tnumbers[split] != 0 {
			temp = tnumbers[split]
			sum = sum + temp
		}
	}

	if sum == 0 {
		return 0, errors.New("couldn't format number")
	}

	return sum, nil
}
