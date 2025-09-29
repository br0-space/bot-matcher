// Package ping_test contains tests for the ping example matcher.
package ping_test

import (
	"testing"

	"github.com/br0-space/bot-matcher/examples/ping"
	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var expectedReply = []telegramclient.MessageStruct{{
	ChatID:                0,
	ReplyToMessageID:      123,
	Text:                  "pong",
	Photo:                 "",
	Caption:               "",
	ParseMode:             "",
	DisableWebPagePreview: false,
	DisableNotification:   false,
}}

var tests = []struct {
	in              string
	expectedReplies []telegramclient.MessageStruct
}{
	{"", nil},
	{"foobar", nil},
	{"ping", nil},
	{"/pings", nil},
	{" /ping", nil},
	{"/ping", expectedReply},
	{"/ping foo", expectedReply},
	{"/ping@bot", expectedReply},
	{"/ping@bot foo", expectedReply},
}

// provideMatcher returns a new ping.Matcher instance used by tests.
func provideMatcher() ping.Matcher {
	return ping.MakeMatcher()
}

// newTestMessage creates a minimal webhook message with the provided text for testing.
func newTestMessage(text string) telegramclient.WebhookMessageStruct {
	return telegramclient.TestWebhookMessage(text)
}

// TestMatcher_DoesMatch ensures that the matcher correctly identifies inputs it should respond to.
func TestMatcher_DoesMatch(t *testing.T) {
	t.Parallel()

	for _, tt := range tests {
		doesMatch := provideMatcher().DoesMatch(newTestMessage(tt.in))
		assert.Equal(t, tt.expectedReplies != nil, doesMatch, tt.in)
	}
}

// TestMatcher_Process validates the replies and error behavior returned by Process for each input case.
func TestMatcher_Process(t *testing.T) {
	t.Parallel()

	for _, tt := range tests {
		replies, err := provideMatcher().Process(newTestMessage(tt.in))
		if tt.expectedReplies == nil {
			require.Error(t, err, tt.in)
			assert.Nil(t, replies, tt.in)
		} else {
			require.NoError(t, err, tt.in)
			assert.NotNil(t, replies, tt.in)
			assert.Equal(t, tt.expectedReplies, replies, tt.in)
		}
	}
}
