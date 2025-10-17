// Package ping provides a simple example matcher that responds to the /ping command
// with a single "pong" message. It demonstrates basic command matching and replying.
package ping

import (
	"errors"
	"regexp"

	matcher "github.com/br0-space/bot-matcher"
	telegramclient "github.com/br0-space/bot-telegramclient"
)

// identifier is the unique name of this matcher.
const identifier = "ping"

// pattern matches /ping, optionally with a bot username suffix and arguments.
var pattern = regexp.MustCompile(`(?i)^/(ping)(@\w+)?($| )`)

// help describes how to use this matcher and can be rendered in help messages.
var help = []matcher.HelpStruct{{
	Command:     "ping",
	Description: `Responds with "pong"`,
	Usage:       "/ping",
	Example:     "/ping",
}}

// template is the text sent back to the user when the matcher processes a ping.
const template = `pong`

// Matcher is a ping matcher that embeds the base matcher and implements Process.
type Matcher struct {
	matcher.Matcher
}

// MakeMatcher constructs a new ping.Matcher with its identifier, pattern, and help configured.
func MakeMatcher() Matcher {
	return Matcher{
		Matcher: matcher.MakeMatcher(identifier, pattern, help),
	}
}

// Process handles an incoming message. If the message matches the /ping command
// it returns a single reply with "pong"; otherwise it returns an error and no replies.
func (m Matcher) Process(messageIn telegramclient.WebhookMessageStruct) ([]telegramclient.MessageStruct, error) {
	if !m.DoesMatch(messageIn) {
		return nil, errors.New("message does not match")
	}

	return []telegramclient.MessageStruct{
		telegramclient.Reply(template, messageIn.ID),
	}, nil
}
