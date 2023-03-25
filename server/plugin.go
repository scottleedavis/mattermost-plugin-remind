package main

import (
	"io/ioutil"
	"time"

	"github.com/gorilla/mux"
	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
	"github.com/pkg/errors"
)

const (
	botUserName    = "remindbot"
	botDisplayName = "Remindbot"
)

type Plugin struct {
	plugin.MattermostPlugin
	client *pluginapi.Client

	router      *mux.Router
	botUserId   string
	running     bool
	emptyTime   time.Time
	defaultTime time.Time

	ServerConfig *model.Config

	readFile func(path string) ([]byte, error)
	locales  map[string]string
}

func NewPlugin() *Plugin {
	return &Plugin{
		readFile: ioutil.ReadFile,
		locales:  make(map[string]string),
	}
}

func (p *Plugin) OnActivate() error {
	p.client = pluginapi.NewClient(p.API, p.Driver)

	p.ServerConfig = p.API.GetConfig()
	p.router = p.InitAPI()
	p.emptyTime = time.Time{}.AddDate(1, 1, 1)

	botID, err := p.client.Bot.EnsureBot(&model.Bot{
		Username:    botUserName,
		DisplayName: botDisplayName,
		Description: "Created by the GitHub plugin.",
	}, pluginapi.ProfileImagePath("assets/icon.png"))
	if err != nil {
		return errors.Wrap(err, "failed to ensure remind bot")
	}
	p.botUserId = botID

	err = p.registerCommand()
	if err != nil {
		return errors.Wrap(err, "failed to register command")
	}

	if err := p.TranslationsPreInit(); err != nil {
		return errors.Wrap(err, "failed to initialize translations")
	}
	p.Run()

	return nil
}

func (p *Plugin) OnDeactivate() error {
	p.Stop()
	return nil
}
