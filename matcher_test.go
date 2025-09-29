// Package matcher_test contains unit tests for the matcher package.
package matcher_test

import (
	"reflect"
	"regexp"
	"testing"
	"unsafe"

	matcher "github.com/br0-space/bot-matcher"
	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setConfigEnabled sets the unexported field `enabled` of matcher.Config in tests.
func setConfigEnabled(cfg *matcher.Config, val *bool) {
	v := reflect.ValueOf(cfg).Elem().FieldByName("enabled")
	reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem().Set(reflect.ValueOf(val))
}

// TestMakeMatcher_Identifier_Help verifies Identifier and Help accessors of the base matcher.
func TestMakeMatcher_Identifier_Help(t *testing.T) {
	t.Parallel()

	pattern := regexp.MustCompile(`(?i)^/test(?:@[-_a-z0-9]+)?(?:\s+.*)?$`)
	helps := []matcher.HelpStruct{{Command: "/test", Description: "desc", Usage: "/test <arg>", Example: "/test foo"}}

	m := matcher.MakeMatcher("test", pattern, helps)

	assert.Equal(t, "test", m.Identifier())
	assert.Equal(t, helps, m.Help())
}

// TestMatcher_IsEnabled_DefaultAndWithConfig documents the default enabled state and config-driven overrides.
func TestMatcher_IsEnabled_DefaultAndWithConfig(t *testing.T) {
	t.Parallel()

	pattern := regexp.MustCompile(`test`)
	m := matcher.MakeMatcher("id", pattern, nil)

	// default: no config -> enabled
	assert.True(t, m.IsEnabled())

	// config with nil flag -> enabled
	m2 := m.WithConfig(&matcher.Config{})
	assert.True(t, m2.IsEnabled())

	// enabled = true
	trueVal := true
	cfgTrue := &matcher.Config{}
	setConfigEnabled(cfgTrue, &trueVal)
	m3 := m.WithConfig(cfgTrue)
	assert.True(t, m3.IsEnabled())

	// enabled = false
	falseVal := false
	cfgFalse := &matcher.Config{}
	setConfigEnabled(cfgFalse, &falseVal)
	m4 := m.WithConfig(cfgFalse)
	assert.False(t, m4.IsEnabled())
}

// TestMatcher_DoesMatch ensures the regex-based matcher detects valid command invocations.
func TestMatcher_DoesMatch(t *testing.T) {
	t.Parallel()

	pattern := regexp.MustCompile(`(?i)^/echo(?:@[-_a-z0-9]+)?(?:\s+.*)?$`)
	m := matcher.MakeMatcher("echo", pattern, nil)
	msg := telegramclient.TestWebhookMessage

	assert.True(t, m.DoesMatch(msg("/echo")))
	assert.True(t, m.DoesMatch(msg("/echo foo")))
	assert.True(t, m.DoesMatch(msg("/echo@mybot foo")))
	assert.False(t, m.DoesMatch(msg("echo")))
	assert.False(t, m.DoesMatch(msg("/echos")))
}

// TestMatcher_CommandMatch validates captured arguments and edge cases for CommandMatch.
func TestMatcher_CommandMatch(t *testing.T) {
	t.Parallel()

	// pattern with a capture group
	pattern := regexp.MustCompile(`(?i)^/say(?:@[-_a-z0-9]+)?\s+(.*)$`)
	m := matcher.MakeMatcher("say", pattern, nil)
	msg := telegramclient.TestWebhookMessage

	// no match -> nil
	assert.Nil(t, m.CommandMatch(msg("/say")))

	// capture content
	assert.Equal(t, []string{"hello world"}, m.CommandMatch(msg("/say hello world")))
	assert.Equal(t, []string{"x"}, m.CommandMatch(msg("/say@bot x")))

	// pattern without capture group returns empty slice on match
	noCap := matcher.MakeMatcher("plain", regexp.MustCompile(`^start$`), nil)
	assert.Equal(t, []string{}, noCap.CommandMatch(msg("start")))
}

// TestMatcher_InlineMatches checks that InlineMatches returns trimmed matches and an empty slice when there are none.
func TestMatcher_InlineMatches(t *testing.T) {
	t.Parallel()

	// collect and trim matches
	pattern := regexp.MustCompile(`foo\s+`)
	m := matcher.MakeMatcher("inline", pattern, nil)
	msg := telegramclient.TestWebhookMessage

	matches := m.InlineMatches(msg("foo  bar foo \n baz"))
	assert.Equal(t, []string{"foo", "foo"}, matches)

	// no matches -> empty slice, not nil
	empty := m.InlineMatches(msg("nothing here"))
	require.NotNil(t, empty)
	assert.Empty(t, empty)
}

// TestMatcher_HandleError_NoPanic ensures that HandleError does not panic when logging errors.
func TestMatcher_HandleError_NoPanic(t *testing.T) {
	t.Parallel()

	m := matcher.MakeMatcher("err", regexp.MustCompile(`.`), nil)
	msg := telegramclient.TestWebhookMessage("text")

	assert.NotPanics(t, func() {
		m.HandleError(msg, "err", assert.AnError)
	})
}

// nonEmbeddingCfg is a sample typed config that does NOT embed matcher.Config.
type nonEmbeddingCfg struct {
	Value string
}

// embeddingCfg embeds matcher.Config and exposes its pointer via the same helper
// interface used by the configurable example.
type embeddingCfg struct {
	matcher.Config

	Flag string
}

func (c embeddingCfg) GetEmbeddedMatcherConfigPtr() *matcher.Config { return &c.Config }

// TestWithCustomConfigType_BasicWiring verifies MakeMatcherWithCustomConfigType wires
// identifier, help, regex, and returns typed config while defaulting to enabled
// when the typed config does not embed matcher.Config.
func TestWithCustomConfigType_BasicWiring(t *testing.T) {
	t.Parallel()

	pattern := regexp.MustCompile(`^abc$`)
	helps := []matcher.HelpStruct{{Command: "abc", Description: "d", Usage: "/abc", Example: "/abc"}}
	cfg := nonEmbeddingCfg{Value: "v"}

	m := matcher.MakeMatcherWithCustomConfigType("id", pattern, helps, cfg)

	// base matcher features are accessible
	assert.Equal(t, "id", m.Identifier())
	assert.Equal(t, helps, m.Help())
	assert.True(t, m.DoesMatch(telegramclient.TestWebhookMessage("abc")))
	assert.False(t, m.DoesMatch(telegramclient.TestWebhookMessage("ab")))

	// typed config is returned as-is
	assert.Equal(t, cfg, m.Config())

	// no embedded matcher.Config -> defaults to enabled
	assert.True(t, m.IsEnabled())
}

// TestWithCustomConfigType_EmbeddedMatcherConfig_WiresBase verifies that when the
// typed config embeds matcher.Config and exposes GetEmbeddedMatcherConfigPtr, the
// base IsEnabled behavior is driven by that flag.
func TestWithCustomConfigType_EmbeddedMatcherConfig_WiresBase(t *testing.T) {
	t.Parallel()

	pattern := regexp.MustCompile(`.`)
	helps := []matcher.HelpStruct(nil)

	falseVal := false
	cfgFalse := embeddingCfg{}
	setConfigEnabled(&cfgFalse.Config, &falseVal)
	mFalse := matcher.MakeMatcherWithCustomConfigType("x", pattern, helps, cfgFalse)
	assert.False(t, mFalse.IsEnabled())
	assert.Equal(t, cfgFalse, mFalse.Config())

	trueVal := true
	cfgTrue := embeddingCfg{}
	setConfigEnabled(&cfgTrue.Config, &trueVal)
	mTrue := matcher.MakeMatcherWithCustomConfigType("x", pattern, helps, cfgTrue)
	assert.True(t, mTrue.IsEnabled())
	assert.Equal(t, cfgTrue, mTrue.Config())
}

// TestWithCustomConfigType_WithTypedConfig_Immutable verifies that WithTypedConfig
// returns a modified copy with updated base wiring and does not mutate the original.
func TestWithCustomConfigType_WithTypedConfig_Immutable(t *testing.T) {
	t.Parallel()

	pattern := regexp.MustCompile(`.`)
	m0 := matcher.MakeMatcherWithCustomConfigType("x", pattern, nil, embeddingCfg{})
	assert.True(t, m0.IsEnabled()) // default nil -> enabled

	falseVal := false
	cfgFalse := embeddingCfg{}
	setConfigEnabled(&cfgFalse.Config, &falseVal)

	m1 := m0.WithTypedConfig(cfgFalse)

	// original stays enabled, new one reflects disabled
	assert.True(t, m0.IsEnabled())
	assert.False(t, m1.IsEnabled())
}
