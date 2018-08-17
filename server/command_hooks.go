package main

import (
	"fmt"
	"strings"
	"time"
	"encoding/json"
	// "math/rand"
    // "os/exec"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
	"github.com/google/uuid"
	// "github.com/nu7hatch/gouuid"

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
		// p.triggerReminders()
		p.runner()
	}()
}

func (p *Plugin) run() {
	
	if !p.running {
		p.running = true
		p.runner()
	}
}

// func (p *Plugin) remindersToJson(reminders []Reminder) (string) {
// 	b, _ := json.Marshal(reminders)
// 	return string(b)
// }


// func (p *Plugin) remindersFromJson(data []byte) ([]Reminder, error) {
// 	// b, err := ioutil.ReadAll(data)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }

// 	var reminders []Reminder
// 	err := json.Unmarshal(data, &reminders)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return reminders, nil
// }

func (p *Plugin) triggerReminders() {

	// bytes, err := p.API.KVGet(string(fmt.Sprintf("%v", time.Now().Round(time.Second))))
	bytes, err := p.API.KVGet("skawtus")
	if err != nil {
		p.API.LogError("failed KVGet %s", err)
	} else {
		p.API.LogError( "value: "+string(bytes[:]) )			
	}


}

func (p *Plugin) upsertReminder(request ReminderRequest) ([]Reminder, error) {

	user, u_err := p.API.GetUserByUsername(request.Username)
	
	if u_err != nil {
		p.API.LogError("failed to query user %s", request.Username)
		return []Reminder{}, u_err
	}

	bytes, b_err := p.API.KVGet(user.Username)
	if b_err != nil {
		p.API.LogError("failed KVGet %s", b_err)
		return []Reminder{}, b_err
	}

	// reminder := request.Reminder
	// reminder := Reminder{user.Username, "me", "foo in 2 seconds", []time.Time{}, time.Time{}}
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
		return []Reminder{}, __
	}

	p.API.KVSet(user.Username,ro)

	return reminders, nil

}

// func (p *Plugin) createOccurrences(request ReminderRequest) ([]ReminderOccurrence) {

// }

func (p *Plugin) scheduleReminder(request ReminderRequest) (string, error) {

    // out, err := exec.Command("uuidgen").Output()
    // if err != nil {
    // 	p.API.LogError("here"+fmt.Sprintf("%v", err))
    //     return "to do", err
    // }

	// request.Reminder.Id = fmt.Sprintf("%s", out) //fmt.Sprintf("%x",rand.Intn(100))

	// p.API.LogError("reminder id "+fmt.Sprintf("%v",request.Reminder.Id))

	guid, gerr := uuid.NewRandom()
	if gerr != nil {
		p.API.LogError("Failed to generate guid")
	}

	request.Reminder.Id = guid.String()

	request.Reminder.Username = request.Username

	request.Reminder.Target = "me"

	request.Reminder.Message = "super foo bar"

	// request.Reminder.Occurrences := p.createOccurrences(request)


	reminders, err := p.upsertReminder(request)
	if err != nil {
		p.API.LogError("failed to query user reminders "+request.Username)
	}


	return fmt.Sprintf("%v",reminders), nil
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

