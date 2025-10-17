package configurable_test

import (
	"testing"

	matcher "github.com/br0-space/bot-matcher"
	"github.com/br0-space/bot-matcher/examples/configurable"
	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConfig_Command verifies the Command accessor returns the default and custom values.
func TestConfig_Command(t *testing.T) {
	t.Parallel()

	// default
	var c configurable.Config
	assert.Equal(t, "configurable", c.Command())
}

// TestConfig_Reply verifies the Reply accessor returns the default and custom values.
func TestConfig_Reply(t *testing.T) {
	t.Parallel()

	// default
	var c configurable.Config
	assert.Equal(t, "unconfigured reply", c.Reply())
}

// Helper for building a minimal message. Marked as a test helper for better failure reporting.
func msg(t *testing.T, text string) telegramclient.WebhookMessageStruct {
	t.Helper()

	return telegramclient.TestWebhookMessage(text)
}

// TestConfig_Pattern_Default ensures the default pattern matches /configurable forms only.
func TestConfig_Pattern_Default(t *testing.T) {
	t.Parallel()

	var c configurable.Config

	re := c.Pattern()

	cases := []struct {
		in   string
		want bool
	}{
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

	for _, cse := range cases {
		assert.Equal(t, cse.want, re.MatchString(msg(t, cse.in).TextOrCaption()), cse.in)
	}
}

// TestConfig_Help verifies Help builds entries from the description with sensible defaults for other fields.
func TestConfig_Help(t *testing.T) {
	t.Parallel()

	// defaults
	var c configurable.Config

	help := c.Help()
	require.Len(t, help, 1)
	assert.Equal(t, matcher.HelpStruct{
		Command:     "configurable",
		Description: "Responds with a configured reply",
		Usage:       "/configurable",
		Example:     "/configurable",
	}, help[0])

	// custom description; command remains default due to encapsulation
	c.Description = "Say hello"

	help = c.Help()
	require.Len(t, help, 1)
	assert.Equal(t, matcher.HelpStruct{
		Command:     "configurable",
		Description: "Say hello",
		Usage:       "/configurable",
		Example:     "/configurable",
	}, help[0])
}
