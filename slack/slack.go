package slack

import (
	"fmt"

	slackApi "github.com/nlopes/slack"
)

var API *slackApi.Client
var slackRtm *slackApi.RTM

var messageID = 0
var token string

type Slack struct {
}

// SendMessage sends a message to channel name or group/DM ID
func (s *Slack) SendMessage(channel string, text string) {
	msg := slackApi.OutgoingMessage{
		ID:      messageID,
		Channel: channel,
		Text:    text,
		Type:    "message",
	}
	messageID++

	slackRtm.SendMessage(&msg)
}

// Connect establishes a connection to slack API
func (s *Slack) Connect(tkn string) {
	token = tkn
	API = slackApi.New(tkn)

	slackApi.SetLogger(logger)
	API.SetDebug(true)

	slackRtm = API.NewRTM()
	go slackRtm.ManageConnection()

	go processEvents(s)
}

func processEvents(s *Slack) {
	for {
		select {
		case msg := <-slackRtm.IncomingEvents:
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slackApi.HelloEvent:
				// Ignore hello

			case *slackApi.ConnectedEvent:
				fmt.Println("Infos:", ev.Info)
				fmt.Println("Connection counter:", ev.ConnectionCount)
				// Replace #general with your Channel ID
				slackRtm.SendMessage(slackRtm.NewOutgoingMessage("Hello world", "#general"))

			case *slackApi.MessageEvent:
				fmt.Printf("Message: %v\n", ev)

			case *slackApi.PresenceChangeEvent:
				fmt.Printf("Presence Change: %v\n", ev)

			case *slackApi.LatencyReport:
				fmt.Printf("Current latency: %v\n", ev.Value)

			case *slackApi.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slackApi.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break

			default:

				fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}
}