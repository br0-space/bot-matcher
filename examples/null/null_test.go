// Package null_test contains tests for the null example matcher.
package null_test

import (
	"testing"

	"github.com/br0-space/bot-matcher/examples/null"
	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var tests = []struct {
	in              string
	expectedReplies []telegramclient.MessageStruct
}{
	{"", nil},
	{"foobar", nil},
	{"null", nil},
	{"/nulls", nil},
	{" /null", nil},
	{"/null", nil},
	{"/null foo", nil},
	{"/null@bot", nil},
	{"/null@bot foo", nil},
}

// provideMatcher returns a new null.Matcher instance used by tests.
func provideMatcher() null.Matcher {
	return null.MakeMatcher()
}

// newTestMessage creates a minimal webhook message with the provided text for testing.
func newTestMessage(text string) telegramclient.WebhookMessageStruct {
	return telegramclient.TestWebhookMessage(text)
}

// TestMatcher_DoesMatch ensures that the null matcher never matches any input.
func TestMatcher_DoesMatch(t *testing.T) {
	t.Parallel()

	for _, tt := range tests {
		doesMatch := provideMatcher().DoesMatch(newTestMessage(tt.in))
		assert.False(t, doesMatch, tt.in)
	}
}

// TestMatcher_Process verifies that Process always returns an error and no replies for the null matcher.
func TestMatcher_Process(t *testing.T) {
	t.Parallel()

	for _, tt := range tests {
		replies, err := provideMatcher().Process(newTestMessage(tt.in))
		require.Error(t, err, tt.in)
		assert.Nil(t, replies, tt.in)
	}
}
