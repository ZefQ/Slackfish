package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/neversun/Slackfish/settings"
	slack "github.com/neversun/Slackfish/slack"
	qml "gopkg.in/qml.v1"
)

const (
	// Appname is the package name
	Appname = "harbour-slackfish"
)

// SlackfishControl is exported to QML and enables interaction with API
type SlackfishControl struct {
	Root            qml.Object
	Slack           *slack.Model
	ChannelsModel   *slack.Channels
	SettingsModel   *settings.Settings
	MessagesModel   *slack.Messages
	UsersModel      *slack.Users
	ImChannelsModel *slack.IMs
}

func main() {
	if err := qml.SailfishRun(run); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ss := &settings.Settings{}

	slackfish := SlackfishControl{
		Slack:           &slack.Slack,
		ChannelsModel:   &slack.Slack.Channels,
		SettingsModel:   ss,
		MessagesModel:   &slack.Slack.Messages,
		UsersModel:      &slack.Slack.Users,
		ImChannelsModel: &slack.Slack.IMs,
	}

	engine := qml.SailfishNewEngine()

	path, err := getPath()
	if err != nil {
		panic(err)
	}
	dataDir := filepath.Join(path, ".local", "share", Appname)
	set := settings.Settings{Location: filepath.Join(dataDir, "settings.yml")}
	slackfish.SettingsModel = &set

	if err = slackfish.SettingsModel.Load(); err != nil {
		log.Printf("WARN: %+v\n", err)
	}

	// TODO: implement translation
	// engine.Translator("/usr/share/harbour-slackfish/qml/i18n")

	context := engine.Context()
	context.SetVar("slackfishctrl", &slackfish)
	// publish direct bindings to public members of slackfish
	context.SetVars(&slackfish)

	controls, err := engine.SailfishSetSource("qml/" + Appname + ".qml")
	if err != nil {
		return err
	}

	window := controls.SailfishCreateWindow()
	slackfish.Root = window.Root()

	window.SailfishShow()
	window.Wait()

	return nil
}

func getPath() (string, error) {
	path := os.Getenv("XDG_DATA_HOME")
	if len(path) == 0 {
		path = os.Getenv("HOME")
		if len(path) == 0 {
			return "", errors.New("No XDG_DATA or HOME env set")
		}
	}
	return path, nil
}
