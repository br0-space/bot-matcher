package matcher_test

import (
	"sync"
	"testing"

	logger "github.com/br0-space/bot-logger"
	matcher "github.com/br0-space/bot-matcher"
	"github.com/br0-space/bot-matcher/examples/null"
	"github.com/br0-space/bot-matcher/examples/ping"
	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/stretchr/testify/assert"
)

// fakeTelegramClient is a minimal test double that records sent messages.
type fakeTelegramClient struct {
	mu      sync.Mutex
	sentTo  []int64
	sentMsg []telegramclient.MessageStruct
}

// SendMessage records the target chat ID and message payload for inspection in tests.
func (f *fakeTelegramClient) SendMessage(chatID int64, msg telegramclient.MessageStruct) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	f.sentTo = append(f.sentTo, chatID)
	f.sentMsg = append(f.sentMsg, msg)

	return nil
}

// TestRegistry_Process_WithPingAndNull verifies that only the ping matcher responds and the null matcher never does.
func TestRegistry_Process_WithPingAndNull(t *testing.T) {
	t.Parallel()

	// Arrange
	client := &fakeTelegramClient{}
	reg := matcher.NewRegistry(logger.New(), client)

	// Register both example matchers.
	reg.Register(ping.MakeMatcher())
	reg.Register(null.MakeMatcher())

	// Incoming message that should trigger ping but never null.
	msg := telegramclient.TestWebhookMessage("/ping")

	// Act
	reg.Process(msg)

	// Assert
	// Exactly one message was sent by the ping matcher.
	if assert.Len(t, client.sentTo, 1) {
		assert.Equal(t, msg.Chat.ID, client.sentTo[0])
	}

	if assert.Len(t, client.sentMsg, 1) {
		// ping replies with a single "pong" reply to the incoming message ID
		expected := telegramclient.MessageStruct{
			ReplyToMessageID: 123, // from TestWebhookMessage
			Text:             "pong",
		}
		assert.Equal(t, expected.ReplyToMessageID, client.sentMsg[0].ReplyToMessageID)
		assert.Equal(t, expected.Text, client.sentMsg[0].Text)
	}
}
