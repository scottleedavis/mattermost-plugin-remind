package main

import (
	"fmt"
	"strings"
	// "time"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
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


func runSchedule() {
	runOnce()
}

func runOnce() {

	fmt.Println("runOnce")	

	// maxThreads := 100

	// if echoSem == nil {
	// 	// We want one additional thread allowed so we never reach channel lockup
	// 	echoSem = make(chan bool, maxThreads+1)
	// }
	// echoSem <- true
	// go(func() {
	// 	defer func() { <-echoSem }()
	// 	// post := &model.Post{}
	// 	// post.ChannelId = args.ChannelId
	// 	// post.RootId = args.RootId
	// 	// post.ParentId = args.ParentId
	// 	// post.Message = message
	// 	// post.UserId = args.UserId

	// 	// time.Sleep(time.Duration(delay) * time.Second)
	// 	time.Sleep(1 * time.Second)

	// 	// if _, err := a.CreatePostMissingChannel(post, true); err != nil {
	// 	// 	mlog.Error(fmt.Sprintf("Unable to create /echo post, err=%v", err))
	// 	// }
	//     fmt.Printf("Current Unix Time: %v\n", time.Now().Unix())
	//     runSchedule()

	// })

// 	   // Timers represent a single event in the future. You
//     // tell the timer how long you want to wait, and it
//     // provides a channel that will be notified at that
//     // time. This timer will wait 2 seconds.
// 	fmt.Println("Timer 1 starting")	
//     // timer1 := time.NewTimer(2 * time.Second)

//     // The `<-timer1.C` blocks on the timer's channel `C`
//     // until it sends a value indicating that the timer
//     // expired.
//     // go func() {
//     //     <-timer1.C
//     // 	fmt.Println("Timer 1 expired")	
//     // }()
// timer1.Stop()

    // If you just wanted to wait, you could have used
    // `time.Sleep`. One reason a timer may be useful is
    // that you can cancel the timer before it expires.
    // Here's an example of that.
    // timer2 := time.NewTimer(time.Second)
    // go func() {
    //     <-timer2.C
    //     fmt.Println("Timer 2 expired")
    // }()
    // stop2 := timer2.Stop()
    // if stop2 {
    //     fmt.Println("Timer 2 stopped")
    // }
}

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
	
	// p.API.LogError("registerCommand %s \n", teamId)

	// runSchedule()

	return nil
}

func (p *Plugin) emitStatusChange() {
	// p.API.PublishWebSocketEvent("status_change", map[string]interface{}{
	// 	"enabled": true,
	// }, &model.WebsocketBroadcast{})
	fmt.Println("hahahahahahha")
	p.API.LogError("emitStatusChange")
	// runSchedule()
}


func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {

	p.API.LogError("ExecuteCommand")

	user, err := p.API.GetUser(args.UserId)
	
	if err != nil {
		p.API.LogError("failed to query user %s", args.UserId)
	}

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
			Text: fmt.Sprintf("* %s\n * %s\n * %s\n * %s\n * %s\n * (%s)\n * %s\n", 
				args.Command, 
				args.TeamId,
				args.SiteURL,
				args.Session,
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

	if commandSplit[1] == "me" {

		p.emitStatusChange()
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text: fmt.Sprintf("todo"),
		}, nil
	}

	if strings.HasPrefix(commandSplit[1][:1], "@"){

		p.emitStatusChange()
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:  fmt.Sprintf("todo"),
		}, nil
	}

	if strings.HasPrefix(commandSplit[1][:1], "~") {

		p.emitStatusChange()
		return &model.CommandResponse{
			ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
			Text:  fmt.Sprintf("todo"),
		}, nil
	}

	return &model.CommandResponse{
		ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text: fmt.Sprintf(ExceptionText),
	}, nil

}

