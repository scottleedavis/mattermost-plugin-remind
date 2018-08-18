package main

import (
	"fmt"
	"strings"
	"time"
	"encoding/json"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/google/uuid"
)

const Version = "0.0.1"

const CommandTrigger = "remind"

const ExceptionText = "Sorry, I didn’t quite get that. I’m easily confused. " +
            "Perhaps try the words in a different order? This usually works: " +
            "`/remind [@someone or ~channel] [what] [when]`.\n";

const HelpText = ":wave: Need some help with `/remind`?\n" +
            "Use `/remind` to set a reminder for yourself, someone else, or for a channel. Some examples include:\n" +
            "* `/remind me to drink water at 3pm every day`\n" +
            "* `/remind me on June 1st to wish Linda happy birthday`\n" +
            "* `/remind ~team-alpha to update the project status every Monday at 9am`\n" +
            "* `/remind @jessica about the interview in 3 hours`\n" +
            "* `/remind @peter tomorrow \"Please review the office seating plan\"`\n" +
            "Or, use `/remind list` to see the list of all your reminders.\n" +
            "Have a bug to report or a feature request?  [Submit your issue here](https://gitreports.com/issue/scottleedavis/mattermost-plugin-remind).";

func (p *Plugin) registerCommand(teamId string) error {

	if err := p.API.RegisterCommand(&model.Command{
		TeamId:           teamId,
		Trigger:          CommandTrigger,
		AutoComplete:     true,
		AutoCompleteHint: "[@someone or ~channel] [what] [when]",
		AutoCompleteDesc: "Enables or disables the demo plugin hooks.",
		DisplayName:      "Remind Plugin Command",
		Description:      "Set a reminder",
	}); err != nil {
		p.API.LogError(
			"failed to register command",
			"error", err.Error(),
		)
	}
	
	p.run()

	return nil
}

func (p *Plugin) runner() {

    go func() {
		<-time.NewTimer(time.Second).C
		p.triggerReminders()
		p.runner()
	}()
}

func (p *Plugin) run() {
	
	if !p.running {
		p.running = true
		p.runner()
	}
}

func (p *Plugin) triggerReminders() {

	bytes, err := p.API.KVGet(string(fmt.Sprintf("%v", time.Now().Round(time.Second))))

	if err != nil {
		p.API.LogError("failed KVGet %s", err)
	} else {
		if string(bytes[:]) != "" {

			// p.API.LogError( "value: "+string(bytes[:]) )			
			var reminderOccurrences []ReminderOccurrence

			ro_err := json.Unmarshal(bytes, &reminderOccurrences)
			if ro_err == nil {
				p.API.LogError("existing "+fmt.Sprintf("%v",reminderOccurrences))















				// TODO loop through array of occurrences, and trigger DM between remind user & user















			}

		}
	}


}

func (p *Plugin) parseRequest(request ReminderRequest) (string, string, string, error) {
	return "me", "in 2 seconds", "foo bar", nil
}

func (p *Plugin) upsertOccurrence(reminderOccurrence ReminderOccurrence) {

	bytes, err := p.API.KVGet(string(fmt.Sprintf("%v", reminderOccurrence.Occurrence)))
	if err != nil {
		p.API.LogError("failed KVGet %s", err)
		return
	}

	var reminderOccurrences []ReminderOccurrence

	ro_err := json.Unmarshal(bytes, &reminderOccurrences)
	if ro_err != nil {
		p.API.LogError("new occurruence " + string(fmt.Sprintf("%v", reminderOccurrence.Occurrence)))
	} else {
		p.API.LogError("existing "+fmt.Sprintf("%v",reminderOccurrences))
	}

	reminderOccurrences = append(reminderOccurrences, reminderOccurrence)
	ro,__ := json.Marshal(reminderOccurrences)

	if __ != nil {
		p.API.LogError("failed to marshal reminderOccurrences %s", reminderOccurrence.Id)
		return
	}

	p.API.KVSet(string(fmt.Sprintf("%v", reminderOccurrence.Occurrence)),ro)

}

func (p *Plugin) createOccurrences(request ReminderRequest) ([]ReminderOccurrence) {

	var ReminderOccurrences []ReminderOccurrence

	// switch the when patterns

	// handle seconds as proof of concept


	guid, gerr := uuid.NewRandom()
	if gerr != nil {
		p.API.LogError("Failed to generate guid")
		return []ReminderOccurrence{}
	}

	occurrence := time.Now().Round(time.Second).Add(time.Second * time.Duration(5))
	reminderOccurrence := ReminderOccurrence{guid.String(),request.Username, request.Reminder.Id, occurrence, time.Time{}, ""}
	ReminderOccurrences = append(ReminderOccurrences, reminderOccurrence)

	p.upsertOccurrence(reminderOccurrence)

	return ReminderOccurrences
}

func (p *Plugin) upsertReminder(request ReminderRequest) {

	user, u_err := p.API.GetUserByUsername(request.Username)
	
	if u_err != nil {
		p.API.LogError("failed to query user %s", request.Username)
		return
	}

	bytes, b_err := p.API.KVGet(user.Username)
	if b_err != nil {
		p.API.LogError("failed KVGet %s", b_err)
		return
	}

	var reminders []Reminder
	err := json.Unmarshal(bytes, &reminders)

	if err != nil {
		p.API.LogError("new reminder " + user.Username)
	} else {
		p.API.LogError("existing "+fmt.Sprintf("%v",reminders))
	}

	reminders = append(reminders, request.Reminder)
	ro,__ := json.Marshal(reminders)

	if __ != nil {
		p.API.LogError("failed to marshal reminders %s", user.Username)
		return
	}

	p.API.KVSet(user.Username,ro)
}

func (p *Plugin) scheduleReminder(request ReminderRequest) (string, error) {

	var when string
	var target string
	var message string
	var useTo bool
	useTo = false
	var useToString string
	if useTo {
		useToString = " to"
	} else {
		useToString = ""
	}

	guid, gerr := uuid.NewRandom()
	if gerr != nil {
		p.API.LogError("Failed to generate guid")
	}

	target, when, message, perr := p.parseRequest(request)
	if perr != nil {
		return ExceptionText, nil
	}

	request.Reminder.Id = guid.String()
	request.Reminder.Username = request.Username
	request.Reminder.Target = target
	request.Reminder.Message = message
	request.Reminder.When = when
	request.Reminder.Occurrences = p.createOccurrences(request)

	p.API.KVDelete(request.Username)

	p.upsertReminder(request)

	response := ":thumbsup: I will remind " + target + useToString + " \"" + request.Reminder.Message + "\" " + when;
   
	return response, nil
}

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {

	p.API.LogError("ExecuteCommand")

	user, err := p.API.GetUser(args.UserId)
	
	if err != nil {
		p.API.LogError("failed to query user %s", args.UserId)
	}

	p.run()

	if strings.HasSuffix(args.Command, "help") {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text: fmt.Sprintf(HelpText),
		}, nil
	}

	if strings.HasSuffix(args.Command, "list") {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text: fmt.Sprintf("todo"),
		}, nil
	}

	if strings.HasSuffix(args.Command, "version") {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text: fmt.Sprintf(Version),
		}, nil
	}

	if strings.HasSuffix(args.Command, "debug") {
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text: fmt.Sprintf("* %s\n * %s\n * %s\n * %s\n * %s\n * %s\n", 
				args.Command, 
				args.TeamId,
				args.SiteURL,
				user.Username, 
				user.Id,  
				user.Timezone["automaticTimezone"]),
		}, nil
	}

	payload := strings.Trim(strings.Replace(args.Command, "/"+CommandTrigger, "", -1),"")
	commandSplit := strings.Split(payload," ")

	if len(commandSplit) == 0 {	

		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text: fmt.Sprintf(ExceptionText),
		}, nil
	}

	if commandSplit[1] == "me" ||
		strings.HasPrefix(commandSplit[1][:1], "@") ||
		strings.HasPrefix(commandSplit[1][:1], "~") {

		request := ReminderRequest{user.Username, payload, Reminder{}}
		response, err := p.scheduleReminder(request)

		if err != nil {
			return &model.CommandResponse{
				ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
				Text: fmt.Sprintf(ExceptionText),
			}, nil
		}

		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text: fmt.Sprintf("%s",response),
		}, nil
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text: fmt.Sprintf(ExceptionText),
	}, nil

}

