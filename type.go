package matcher

import telegramclient "github.com/br0-space/bot-telegramclient"

type Interface interface {
	IsEnabled() bool
	Identifier() string
	Help() []HelpStruct
	DoesMatch(messageIn telegramclient.WebhookMessageStruct) bool
	GetCommandMatch(messageIn telegramclient.WebhookMessageStruct) []string
	GetInlineMatches(messageIn telegramclient.WebhookMessageStruct) []string
	Process(messageIn telegramclient.WebhookMessageStruct) ([]telegramclient.MessageStruct, error)
	HandleError(messageIn telegramclient.WebhookMessageStruct, identifier string, err error)
}

type HelpStruct struct {
	Command     string
	Description string
	Usage       string
	Example     string
}
