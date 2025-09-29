package configurable

import (
	"errors"

	matcher "github.com/br0-space/bot-matcher"
	telegramclient "github.com/br0-space/bot-telegramclient"
)

const identifier = "configurable"

// Matcher provides a command matcher whose command and reply are loaded from config.
// The config structure is defined in config.go.
type Matcher struct {
	matcher.WithCustomConfigType[Config]
}

// MakeMatcher constructs a new configurable matcher using values from config/configurable.yaml.
// It loads the command, reply, and description from the config, builds the matching pattern accordingly,
// and wires the base matcher with that pattern and a generated help entry.
func MakeMatcher() Matcher {
	cfgs := matcher.LoadMatcherConfig[Config](identifier)
	cfg := cfgs[0]

	pattern := cfg.Pattern()
	help := cfg.Help()

	return Matcher{
		WithCustomConfigType: matcher.MakeMatcherWithCustomConfigType(identifier, pattern, help, cfg),
	}
}

// Process checks whether the message matches the configured command and replies with the
// configured reply text. If the message does not match, it returns an error.
func (m Matcher) Process(messageIn telegramclient.WebhookMessageStruct) ([]telegramclient.MessageStruct, error) {
	if !m.DoesMatch(messageIn) {
		return nil, errors.New("message does not match")
	}

	replyText := m.Config().Reply()

	return []telegramclient.MessageStruct{
		telegramclient.Reply(replyText, messageIn.ID),
	}, nil
}
