// Package configurable_test contains tests for the configurable example matcher.
package configurable_test

import (
	"os"
	"path/filepath"
	"testing"

	matcher "github.com/br0-space/bot-matcher"
	"github.com/br0-space/bot-matcher/examples/configurable"
	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// helper to write the config/configurable.yaml file relative to the package CWD.
func writeConfigFile(t *testing.T, yaml string) {
	t.Helper()

	dir := "config"
	require.NoError(t, os.MkdirAll(dir, 0o755))

	path := filepath.Join(dir, "configurable.yml")
	require.NoError(t, os.WriteFile(path, []byte(yaml), 0o600))

	t.Cleanup(func() {
		_ = os.Remove(path)
		_ = os.Remove(dir)
	})
}

// newTestMessage creates a minimal webhook message with the provided text for testing.
func newTestMessage(text string) telegramclient.WebhookMessageStruct {
	return telegramclient.TestWebhookMessage(text)
}

// TestMatcher_Defaults verifies default behavior when no command/reply are configured.
func TestMatcher_Defaults(t *testing.T) { //nolint:paralleltest
	// empty config file; all values default in accessors
	writeConfigFile(t, "")

	m := configurable.MakeMatcher()

	// DoesMatch cases mirror the ping tests, but for /configurable
	type tc struct {
		in      string
		matches bool
	}

	cases := []tc{
		{"", false},
		{"foobar", false},
		{"configurable", false},
		{"/configurables", false},
		{" /configurable", false},
		{"/configurable", true},
		{"/configurable foo", true},
		{"/configurable@bot", true},
		{"/configurable@bot foo", true},
	}

	for _, c := range cases {
		assert.Equal(t, c.matches, m.DoesMatch(newTestMessage(c.in)), c.in)
	}

	// Process should reply with the default reply text
	replies, err := m.Process(newTestMessage("/configurable"))
	require.NoError(t, err)

	expected := []telegramclient.MessageStruct{{
		ChatID:                0,
		ReplyToMessageID:      123,
		Text:                  "unconfigured reply",
		Photo:                 "",
		Caption:               "",
		ParseMode:             "",
		DisableWebPagePreview: false,
		DisableNotification:   false,
	}}
	assert.Equal(t, expected, replies)

	// Help should reflect defaults
	help := m.Help()
	require.Len(t, help, 1)
	assert.Equal(t, matcher.HelpStruct{
		Command:     "configurable",
		Description: "Responds with a configured reply",
		Usage:       "/configurable",
		Example:     "/configurable",
	}, help[0])
}

// TestMatcher_CustomCommandAndReply verifies behavior with explicit configuration.
func TestMatcher_CustomCommandAndReply(t *testing.T) { //nolint:paralleltest
	// Provide explicit command, reply and description
	writeConfigFile(t, "command: hello\nreply: world\ndescription: Says world\n")

	m := configurable.MakeMatcher()

	// The matcher should honor the configured command and reply.
	// It should match the configured command, not the default one.
	assert.False(t, m.DoesMatch(newTestMessage("/configurable")))
	assert.True(t, m.DoesMatch(newTestMessage("/hello")))

	// Process should return the configured reply
	replies, err := m.Process(newTestMessage("/hello"))
	require.NoError(t, err)

	expected := []telegramclient.MessageStruct{{
		ChatID:                0,
		ReplyToMessageID:      123,
		Text:                  "world",
		Photo:                 "",
		Caption:               "",
		ParseMode:             "",
		DisableWebPagePreview: false,
		DisableNotification:   false,
	}}
	assert.Equal(t, expected, replies)

	// Help should mirror the configured description and command/usage/example
	help := m.Help()
	require.Len(t, help, 1)
	assert.Equal(t, matcher.HelpStruct{
		Command:     "hello",
		Description: "Says world",
		Usage:       "/hello",
		Example:     "/hello",
	}, help[0])
}
