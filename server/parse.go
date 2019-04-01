package main

import (
	"errors"
	"fmt"
	"strings"
)

func (p *Plugin) ParseRequest(request *ReminderRequest) error {


	user, _ := p.API.GetUserByUsername(request.Username)
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
			when = strings.Replace(when, commandSplit[0], "", -1)
			when = strings.Trim(when, " ")

			message = strings.Replace(message, "\"", "", -1)

			request.Reminder.When = when
			request.Reminder.Message = message
			return nil
		}

		if wErr := a.findWhen(request); wErr != nil {
			return wErr
		}

		message := strings.Replace(request.Payload, request.Reminder.When, "", -1)
		message = strings.Replace(message, commandSplit[0], "", -1)
		message = strings.Trim(message, " \"")

		request.Reminder.Message = message

		return nil

	}

	return errors.New("unrecognized target")
}

func (p *Plugin) findWhen(request *ReminderRequest) error {

	user, _ := p.API.GetUserByUsername(request.Username)
	_, locale := p.translation(user)

	switch locale {
	case "en":
		return a.findWhenEN(request)
	default:
		return a.findWhenEN(request)
	}

}

func (p *Plugin) findWhenEN(request *ReminderRequest) error {

	user, _ := p.API.GetUserByUsername(request.Username)
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
	_, dErr := a.normalizeDate(lastWord, user)
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

		_, dErr = a.normalizeDate(lastWord, user)
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

func (p *Plugin) normalizeDate(text string, user *model.User) (string, error) {

	// cfg := a.Config()
	// location := a.location(user)
	// T, _ := a.translation(user)

	// date := strings.ToLower(text)
	// if strings.EqualFold(T("app.reminder.chrono.day"), date) {
	// 	return date, nil
	// } else if strings.EqualFold(T("app.reminder.chrono.today"), date) {
	// 	return date, nil
	// } else if strings.EqualFold(T("app.reminder.chrono.everyday"), date) {
	// 	return date, nil
	// } else if strings.EqualFold(T("app.reminder.chrono.tomorrow"), date) {
	// 	return date, nil
	// }

	// switch date {
	// case T("app.reminder.chrono.mon"),
	// 	T("app.reminder.chrono.monday"):
	// 	return T("app.reminder.chrono.monday"), nil
	// case T("app.reminder.chrono.tues"),
	// 	T("app.reminder.chrono.tuesday"):
	// 	return T("app.reminder.chrono.tuesday"), nil
	// case T("app.reminder.chrono.wed"),
	// 	T("app.reminder.chrono.wednes"),
	// 	T("app.reminder.chrono.wednesday"):
	// 	return T("app.reminder.chrono.wednesday"), nil
	// case T("app.reminder.chrono.thur"),
	// 	T("app.reminder.chrono.thursday"):
	// 	return T("app.reminder.chrono.thursday"), nil
	// case T("app.reminder.chrono.fri"),
	// 	T("app.reminder.chrono.friday"):
	// 	return T("app.reminder.chrono.friday"), nil
	// case T("app.reminder.chrono.sat"),
	// 	T("app.reminder.chrono.satur"),
	// 	T("app.reminder.chrono.saturday"):
	// 	return T("app.reminder.chrono.saturday"), nil
	// case T("app.reminder.chrono.sun"),
	// 	T("app.reminder.chrono.sunday"):
	// 	return T("app.reminder.chrono.sunday"), nil
	// case T("app.reminder.chrono.mondays"),
	// 	T("app.reminder.chrono.tuesdays"),
	// 	T("app.reminder.chrono.wednesdays"),
	// 	T("app.reminder.chrono.thursdays"),
	// 	T("app.reminder.chrono.fridays"),
	// 	T("app.reminder.chrono.saturdays"),
	// 	T("app.reminder.chrono.sundays"):
	// 	return date, nil
	// }

	// if strings.Contains(date, T("app.reminder.chrono.jan")) ||
	// 	strings.Contains(date, T("app.reminder.chrono.january")) ||
	// 	strings.Contains(date, T("app.reminder.chrono.feb")) ||
	// 	strings.Contains(date, T("app.reminder.chrono.february")) ||
	// 	strings.Contains(date, T("app.reminder.chrono.mar")) ||
	// 	strings.Contains(date, T("app.reminder.chrono.march")) ||
	// 	strings.Contains(date, T("app.reminder.chrono.apr")) ||
	// 	strings.Contains(date, T("app.reminder.chrono.april")) ||
	// 	strings.Contains(date, T("app.reminder.chrono.may")) ||
	// 	strings.Contains(date, T("app.reminder.chrono.june")) ||
	// 	strings.Contains(date, T("app.reminder.chrono.july")) ||
	// 	strings.Contains(date, T("app.reminder.chrono.aug")) ||
	// 	strings.Contains(date, T("app.reminder.chrono.august")) ||
	// 	strings.Contains(date, T("app.reminder.chrono.sept")) ||
	// 	strings.Contains(date, T("app.reminder.chrono.september")) ||
	// 	strings.Contains(date, T("app.reminder.chrono.oct")) ||
	// 	strings.Contains(date, T("app.reminder.chrono.october")) ||
	// 	strings.Contains(date, T("app.reminder.chrono.nov")) ||
	// 	strings.Contains(date, T("app.reminder.chrono.november")) ||
	// 	strings.Contains(date, T("app.reminder.chrono.dec")) ||
	// 	strings.Contains(date, T("app.reminder.chrono.december")) {

	// 	date = strings.Replace(date, ",", "", -1)
	// 	parts := strings.Split(date, " ")

	// 	switch len(parts) {
	// 	case 1:
	// 		break
	// 	case 2:
	// 		if len(parts[1]) > 2 {
	// 			parts[1] = a.daySuffix(user, parts[1])
	// 		}
	// 		if _, err := strconv.Atoi(parts[1]); err != nil {
	// 			if wn, wErr := a.wordToNumber(parts[1], user); wErr == nil {
	// 				parts[1] = strconv.Itoa(wn)
	// 			}
	// 		}

	// 		parts = append(parts, fmt.Sprintf("%v", time.Now().Year()))

	// 		break
	// 	case 3:
	// 		if len(parts[1]) > 2 {
	// 			parts[1] = a.daySuffix(user, parts[1])
	// 		}

	// 		if _, err := strconv.Atoi(parts[1]); err != nil {
	// 			if wn, wErr := a.wordToNumber(parts[1], user); wErr == nil {
	// 				parts[1] = strconv.Itoa(wn)
	// 			} else {
	// 				mlog.Error(wErr.Error())
	// 			}

	// 			if _, pErr := strconv.Atoi(parts[2]); pErr != nil {
	// 				return "", pErr
	// 			}
	// 		}

	// 		break
	// 	default:
	// 		return "", errors.New("unrecognized date format")
	// 	}

	// 	switch parts[0] {
	// 	case T("app.reminder.chrono.jan"),
	// 		T("app.reminder.chrono.january"):
	// 		parts[0] = "01"
	// 		break
	// 	case T("app.reminder.chrono.feb"),
	// 		T("app.reminder.chrono.february"):
	// 		parts[0] = "02"
	// 		break
	// 	case T("app.reminder.chrono.mar"),
	// 		T("app.reminder.chrono.march"):
	// 		parts[0] = "03"
	// 		break
	// 	case T("app.reminder.chrono.apr"),
	// 		T("app.reminder.chrono.april"):
	// 		parts[0] = "04"
	// 		break
	// 	case T("app.reminder.chrono.may"):
	// 		parts[0] = "05"
	// 		break
	// 	case T("app.reminder.chrono.june"):
	// 		parts[0] = "06"
	// 		break
	// 	case T("app.reminder.chrono.july"):
	// 		parts[0] = "07"
	// 		break
	// 	case T("app.reminder.chrono.aug"),
	// 		T("app.reminder.chrono.august"):
	// 		parts[0] = "08"
	// 		break
	// 	case T("app.reminder.chrono.sept"),
	// 		T("app.reminder.chrono.september"):
	// 		parts[0] = "09"
	// 		break
	// 	case T("app.reminder.chrono.oct"),
	// 		T("app.reminder.chrono.october"):
	// 		parts[0] = "10"
	// 		break
	// 	case T("app.reminder.chrono.nov"),
	// 		T("app.reminder.chrono.november"):
	// 		parts[0] = "11"
	// 		break
	// 	case T("app.reminder.chrono.dec"),
	// 		T("app.reminder.chrono.december"):
	// 		parts[0] = "12"
	// 		break
	// 	default:
	// 		return "", errors.New("month not found")
	// 	}

	// 	if len(parts[1]) < 2 {
	// 		parts[1] = "0" + parts[1]
	// 	}
	// 	return parts[2] + "-" + parts[0] + "-" + parts[1] + "T00:00:00Z", nil

	// } else if match, _ := regexp.MatchString("^(([0-9]{2}|[0-9]{1})(-|/)([0-9]{2}|[0-9]{1})((-|/)([0-9]{4}|[0-9]{2}))?)", date); match {

	// 	date := a.regSplit(date, "-|/")

	// 	switch len(date) {
	// 	case 2:
	// 		year := time.Now().Year()
	// 		month, mErr := strconv.Atoi(date[0])
	// 		if mErr != nil {
	// 			return "", mErr
	// 		}
	// 		day, dErr := strconv.Atoi(date[1])
	// 		if dErr != nil {
	// 			return "", dErr
	// 		}

	// 		if *cfg.DisplaySettings.ExperimentalTimezone {
	// 			return time.Date(year, time.Month(month), day, 0, 0, 0, 0, location).Format(time.RFC3339), nil
	// 		}

	// 		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local).Format(time.RFC3339), nil

	// 	case 3:
	// 		if len(date[2]) == 2 {
	// 			date[2] = "20" + date[2]
	// 		}
	// 		year, yErr := strconv.Atoi(date[2])
	// 		if yErr != nil {
	// 			return "", yErr
	// 		}
	// 		month, mErr := strconv.Atoi(date[0])
	// 		if mErr != nil {
	// 			return "", mErr
	// 		}
	// 		day, dErr := strconv.Atoi(date[1])
	// 		if dErr != nil {
	// 			return "", dErr
	// 		}

	// 		if *cfg.DisplaySettings.ExperimentalTimezone {
	// 			return time.Date(year, time.Month(month), day, 0, 0, 0, 0, location).Format(time.RFC3339), nil
	// 		}

	// 		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local).Format(time.RFC3339), nil

	// 	default:
	// 		return "", errors.New("unrecognized date")
	// 	}

	// } else if match, _ := regexp.MatchString("^(([0-9]{2}|[0-9]{1})(.)([0-9]{2}|[0-9]{1})((.)([0-9]{4}|[0-9]{2}))?)", date); match {

	// 	date := a.regSplit(date, "\\.")

	// 	switch len(date) {
	// 	case 2:
	// 		year := time.Now().Year()
	// 		month, mErr := strconv.Atoi(date[1])
	// 		if mErr != nil {
	// 			return "", mErr
	// 		}
	// 		day, dErr := strconv.Atoi(date[0])
	// 		if dErr != nil {
	// 			return "", dErr
	// 		}

	// 		if *cfg.DisplaySettings.ExperimentalTimezone {
	// 			return time.Date(year, time.Month(month), day, 0, 0, 0, 0, location).Format(time.RFC3339), nil
	// 		}

	// 		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local).Format(time.RFC3339), nil

	// 	case 3:
	// 		if len(date[2]) == 2 {
	// 			date[2] = "20" + date[2]
	// 		}
	// 		year, yErr := strconv.Atoi(date[2])
	// 		if yErr != nil {
	// 			return "", yErr
	// 		}
	// 		month, mErr := strconv.Atoi(date[1])
	// 		if mErr != nil {
	// 			return "", mErr
	// 		}
	// 		day, dErr := strconv.Atoi(date[0])
	// 		if dErr != nil {
	// 			return "", dErr
	// 		}

	// 		if *cfg.DisplaySettings.ExperimentalTimezone {
	// 			return time.Date(year, time.Month(month), day, 0, 0, 0, 0, location).Format(time.RFC3339), nil
	// 		}

	// 		return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local).Format(time.RFC3339), nil

	// 	default:
	// 		return "", errors.New("unrecognized date")
	// 	}

	// } else { //single number day

	// 	var dayInt int
	// 	day := a.daySuffix(user, date)

	// 	if d, nErr := strconv.Atoi(day); nErr != nil {
	// 		if wordNum, wErr := a.wordToNumber(date, user); wErr != nil {
	// 			return "", wErr
	// 		} else {
	// 			day = strconv.Itoa(wordNum)
	// 			dayInt = wordNum
	// 		}
	// 	} else {
	// 		dayInt = d
	// 	}

	// 	month := time.Now().Month()
	// 	year := time.Now().Year()

	// 	var t time.Time
	// 	if *cfg.DisplaySettings.ExperimentalTimezone {
	// 		t = time.Date(year, month, dayInt, 0, 0, 0, 0, location)
	// 	} else {
	// 		t = time.Date(year, month, dayInt, 0, 0, 0, 0, time.Local)
	// 	}

	// 	if t.Before(time.Now()) {
	// 		t = t.AddDate(0, 1, 0)
	// 	}

	// 	return t.Format(time.RFC3339), nil

	// }

}
