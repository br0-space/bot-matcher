package matcher

import (
	telegramclient "github.com/br0-space/bot-telegramclient"
)

type Interface interface {
	IsEnabled() bool
	Identifier() string
	Help() []HelpStruct
	DoesMatch(messageIn telegramclient.WebhookMessageStruct) bool
	CommandMatch(messageIn telegramclient.WebhookMessageStruct) []string
	InlineMatches(messageIn telegramclient.WebhookMessageStruct) []string
	Process(messageIn telegramclient.WebhookMessageStruct) ([]telegramclient.MessageStruct, error)
	HandleError(messageIn telegramclient.WebhookMessageStruct, identifier string, err error)
}
