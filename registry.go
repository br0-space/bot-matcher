package matcher

import (
	"fmt"
	"sync"

	logger "github.com/br0-space/bot-logger"
	telegramclient "github.com/br0-space/bot-telegramclient"
)

const errorTemplate = "⚠️ *Error in matcher \"%s\"*\n\n%s"

type Registry struct {
	log      logger.Interface
	telegram telegramclient.ClientInterface
	matchers []Interface
}

// NewRegistry creates a new Registry using the provided logger and Telegram client.
// It initializes empty matcher storage.
func NewRegistry(
	logger logger.Interface,
	telegram telegramclient.ClientInterface,
) *Registry {
	return &Registry{
		log:      logger,
		telegram: telegram,
		matchers: []Interface{},
	}
}

// Register adds a matcher to the registry.
func (r *Registry) Register(matcher Interface) {
	r.log.Debug("Registering matcher", matcher.Identifier())

	r.matchers = append(r.matchers, matcher)
}

// Process routes an incoming message to all registered matchers concurrently.
// It checks whether each matcher is enabled, evaluates DoesMatch, executes the matcher,
// reports errors to the user as a Markdown reply, sends all returned messages,
// and waits for all matchers to finish.
func (r *Registry) Process(messageIn telegramclient.WebhookMessageStruct) {
	r.log.Debugf("Processing message from %s: %s", messageIn.From.Username, messageIn.Text)

	var waitGroup sync.WaitGroup
	waitGroup.Add(len(r.matchers))

	for _, m := range r.matchers {
		go func(m Interface) {
			defer waitGroup.Done()

			chatID := messageIn.Chat.ID
			if !r.shouldRunMatcher(m, chatID) {
				return
			}

			messagesOut := r.executeMatcher(m, messageIn)
			r.sendMessages(chatID, messagesOut)
		}(m)
	}

	waitGroup.Wait()
}

// shouldRunMatcher encapsulates the decision logic and logging to determine if a matcher
// should be executed for a particular chat.
func (r *Registry) shouldRunMatcher(m Interface, chatID int64) bool {
	if !m.IsEnabled() {
		r.log.Debugf("Matcher %s will not be executed: disabled", m.Identifier())

		return false
	}

	r.log.Debugf("Matcher %s will be executed for chat %d", m.Identifier(), chatID)

	return true
}

// executeMatcher runs DoesMatch and Process and normalizes/augments the output
// by appending a Markdown error reply if Process returned an error.
func (r *Registry) executeMatcher(m Interface, messageIn telegramclient.WebhookMessageStruct) []telegramclient.MessageStruct {
	if !m.DoesMatch(messageIn) {
		return nil
	}

	messagesOut, err := m.Process(messageIn)
	if messagesOut == nil {
		messagesOut = []telegramclient.MessageStruct{}
	}

	if err != nil {
		r.log.Errorf("Error in matcher %s: %s", m.Identifier(), err)
		messagesOut = append(
			messagesOut,
			telegramclient.MarkdownReply(
				fmt.Sprintf(
					errorTemplate,
					m.Identifier(),
					telegramclient.EscapeMarkdown(err.Error()),
				),
				messageIn.ID,
			),
		)
	}

	return messagesOut
}

// sendMessages delivers all messages to the given chat ID and logs errors individually.
func (r *Registry) sendMessages(chatID int64, messagesOut []telegramclient.MessageStruct) {
	for _, messageOut := range messagesOut {
		if err := r.telegram.SendMessage(chatID, messageOut); err != nil {
			r.log.Error("Error while sending message:", err)
		}
	}
}
