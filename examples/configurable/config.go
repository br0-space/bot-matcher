package configurable

import (
	"fmt"
	"regexp"

	matcher "github.com/br0-space/bot-matcher"
)

type Config struct {
	matcher.Config

	CommandText string `mapstructure:"command"`
	ReplyText   string `mapstructure:"reply"`
	Description string `mapstructure:"description"`
}

// GetEmbeddedMatcherConfigPtr exposes a pointer to the embedded matcher.Config.
// This allows the generic matcher.MatcherWithCustomConfig to wire base-level config behaviors like IsEnabled.
func (c Config) GetEmbeddedMatcherConfigPtr() *matcher.Config {
	return &c.Config
}

// Command returns the configured command, or the default "configurable" if unset.
func (c Config) Command() string {
	if c.CommandText == "" {
		return "configurable"
	}

	return c.CommandText
}

// Reply returns the configured reply, or a default text if unset.
func (c Config) Reply() string {
	if c.ReplyText == "" {
		return "unconfigured reply"
	}

	return c.ReplyText
}

// Pattern returns the compiled regular expression used by MakeMatcher based on the
// configured command. If the command is empty, it falls back to "configurable".
func (c Config) Pattern() *regexp.Regexp {
	cmd := c.Command()

	return regexp.MustCompile(fmt.Sprintf(`(?i)^/(%s)(@\w+)?($| )`, regexp.QuoteMeta(cmd)))
}

// Help returns the help entry constructed from config values.
func (c Config) Help() []matcher.HelpStruct {
	cmd := c.Command()

	desc := c.Description
	if desc == "" {
		desc = "Responds with a configured reply"
	}

	return []matcher.HelpStruct{{
		Command:     cmd,
		Description: desc,
		Usage:       "/" + cmd,
		Example:     "/" + cmd,
	}}
}
