package main

import (
	"strings"
	"fmt"
	"errors"
)

func (p *Plugin) ParseRequest(request ReminderRequest) (string, string, string, error) {

	commandSplit := strings.Split(request.Payload, " ")

	p.API.LogDebug("parseRequest " + fmt.Sprintf("%v", request))
	p.API.LogDebug(request.Payload)

	if strings.HasPrefix(request.Payload, "me") ||
		strings.HasPrefix(request.Payload, "~") ||
		strings.HasPrefix(request.Payload, "@") {

		p.API.LogDebug("found target")

		firstIndex := strings.Index(request.Payload, "\"")
		lastIndex := strings.LastIndex(request.Payload, "\"")

		if firstIndex > -1 && lastIndex > -1 && firstIndex != lastIndex { // has quotes

			message := request.Payload[firstIndex-1:lastIndex]
			when := strings.Replace(request.Payload, message, "", -1)
			when = strings.Replace(when, commandSplit[0], "", -1)
			p.API.LogError("quotes when " + fmt.Sprintf("%v", firstIndex) + " " + fmt.Sprintf("%v", lastIndex) + " " + when)

			return commandSplit[0], when, message, nil
		}

		p.API.LogDebug("no quotes when " + fmt.Sprintf("%v", firstIndex) + " " + fmt.Sprintf("%v", lastIndex))

		when, wErr := p.findWhen(request.Payload)
		if wErr != nil {
			return "", "", "", wErr
		}

		message := strings.Replace(request.Payload, when, "", -1)
		message = strings.Replace(message, commandSplit[0], "", -1)
		message = strings.Trim(message, " \"")

		return commandSplit[0], when, message, nil
	}

	return "", "", "", errors.New("unrecognized Target")
}

func (p *Plugin) findWhen(payload string) (string, error) {

	inSplit := strings.Split(payload, " in ")
	if len(inSplit) == 2 {
		return "in "+inSplit[len(inSplit)-1], nil
	}

	//TODO the additional when states

	return "", errors.New("unable to find when")
}
