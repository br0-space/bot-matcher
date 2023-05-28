package matcher

import (
	"fmt"
	logger "github.com/br0-space/bot-logger"
	telegramclient "github.com/br0-space/bot-telegramclient"
	"sync"
)

const errorTemplate = "⚠️ *Error in matcher \"%s\"*\n\n%s"

type Registry struct {
	log      logger.Interface
	telegram telegramclient.ClientInterface
	matchers []Interface
}

func NewRegistry(
	logger logger.Interface,
	telegram telegramclient.ClientInterface,
) *Registry {
	return &Registry{
		log:      logger,
		telegram: telegram,
	}
}

func (r *Registry) Register(matcher Interface) {
	r.log.Debug("Registering matcher", matcher.Identifier())

	r.matchers = append(r.matchers, matcher)
}

func (r *Registry) Process(messageIn telegramclient.WebhookMessageStruct) {
	r.log.Debugf("Processing message from %s: %s", messageIn.From.Username, messageIn.Text)

	// Create a wait group for synchronization
	var waitGroup sync.WaitGroup

	// We need to wait until all matchers are executed
	waitGroup.Add(len(r.matchers))

	// Launch a goroutine for each matcher
	for _, m := range r.matchers {
		go func(m Interface) {
			defer waitGroup.Done()

			if !m.IsEnabled() {
				return
			}

			if !m.DoesMatch(messageIn) {
				return
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

			for _, messageOut := range messagesOut {
				if err := r.telegram.SendMessage(messageIn.Chat.ID, messageOut); err != nil {
					r.log.Error("Error while sending message:", err)
				}
			}
		}(m)
	}

	waitGroup.Wait()

}
