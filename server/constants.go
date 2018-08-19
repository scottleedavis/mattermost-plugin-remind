package main

const Version = "0.0.1"

const CommandTrigger = "remind"

const ExceptionText = "Sorry, I didn’t quite get that. I’m easily confused. " +
	"Perhaps try the words in a different order? This usually works: " +
	"`/remind [@someone or ~channel] [what] [when]`.\n"

const HelpText = ":wave: Need some help with `/remind`?\n" +
	"Use `/remind` to set a reminder for yourself, someone else, or for a channel. Some examples include:\n" +
	"* `/remind me to drink water at 3pm every day`\n" +
	"* `/remind me on June 1st to wish Linda happy birthday`\n" +
	"* `/remind ~team-alpha to update the project status every Monday at 9am`\n" +
	"* `/remind @jessica about the interview in 3 hours`\n" +
	"* `/remind @peter tomorrow \"Please review the office seating plan\"`\n" +
	"Or, use `/remind list` to see the list of all your reminders.\n" +
	"Have a bug to report or a feature request?  [Submit your issue here](https://gitreports.com/issue/scottleedavis/mattermost-plugin-remind)."
