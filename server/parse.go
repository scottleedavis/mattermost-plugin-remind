package main

import (
	"strings"
	"fmt"
	"errors"
)

func (p *Plugin) ParseRequest(request ReminderRequest) (string, string, string, error) {

	commandSplit := strings.Split(request.Payload, " ")

	p.API.LogError("parseRequest " + fmt.Sprintf("%v", request))
	p.API.LogError(request.Payload)

	if strings.HasPrefix(request.Payload, "me") ||
		strings.HasPrefix(request.Payload, "~") ||
		strings.HasPrefix(request.Payload, "@") {

		p.API.LogError("found target")

		var message string
		var when string
		var firstIndex int
		var lastIndex int

		firstIndex = strings.Index(request.Payload, "\"")
		lastIndex = strings.LastIndex(request.Payload, "\"")

		if firstIndex > -1 && lastIndex > -1 && firstIndex != lastIndex {

			message := request.Payload[firstIndex:lastIndex]
			when = strings.Replace(request.Payload, message, "", -1)
			when = strings.Replace(when, commandSplit[1], "", -1)
			p.API.LogError("quotes when " + fmt.Sprintf("%v", firstIndex) + " " + fmt.Sprintf("%v", lastIndex) + " " + when)

			return commandSplit[0], when, message, nil
		}

		p.API.LogError("no quotes when " + fmt.Sprintf("%v", firstIndex) + " " + fmt.Sprintf("%v", lastIndex))

		message = "foo"

		////////
		// TODO determine when
		/////////

		return commandSplit[0], "in 2 seconds", message, nil
	}
	err := errors.New("Unrecognized Target")

	return "", "", "", err
}
