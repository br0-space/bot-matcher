package matcher

import (
	"regexp"
	"strings"

	logger "github.com/br0-space/bot-logger"
	telegramclient "github.com/br0-space/bot-telegramclient"
)

type Matcher struct {
	log        logger.Interface
	identifier string
	regexp     *regexp.Regexp
	help       []HelpStruct
	cfg        *Config
}

type Config struct {
	enabled *bool
}

type HelpStruct struct {
	Command     string
	Description string
	Usage       string
	Example     string
}

// MakeMatcher constructs a new Matcher using the given identifier, regular expression
// pattern and help definitions. The returned matcher has a default logger and no config.
func MakeMatcher(
	identifier string,
	pattern *regexp.Regexp,
	help []HelpStruct,
) Matcher {
	return Matcher{
		log:        logger.New(),
		identifier: identifier,
		regexp:     pattern,
		help:       help,
		cfg:        nil,
	}
}

// WithConfig returns a copy of the Matcher with the provided configuration applied.
func (m Matcher) WithConfig(cfg *Config) Matcher {
	m.cfg = cfg

	return m
}

func (m Matcher) Config() Config {
	if m.cfg == nil {
		return Config{}
	}

	return *m.cfg
}

// IsEnabled reports whether the matcher is enabled.
// If no config is present or the enabled flag is not set, it defaults to true.
func (m Matcher) IsEnabled() bool {
	return m.Config().enabled != nil && *m.Config().enabled || m.Config().enabled == nil // default to true
}

// Identifier returns the unique identifier of the matcher.
func (m Matcher) Identifier() string {
	return m.identifier
}

// Help returns the list of help entries describing how to use this matcher.
func (m Matcher) Help() []HelpStruct {
	return m.help
}

// DoesMatch reports whether the message's text or caption matches the matcher's pattern.
func (m Matcher) DoesMatch(messageIn telegramclient.WebhookMessageStruct) bool {
	return m.regexp.MatchString(messageIn.TextOrCaption())
}

// CommandMatch returns the capturing groups of the first match against the message.
// If there is no match, it returns nil. If there are no capturing groups, it returns an empty slice.
func (m Matcher) CommandMatch(messageIn telegramclient.WebhookMessageStruct) []string {
	match := m.regexp.FindStringSubmatch(messageIn.TextOrCaption())
	if match == nil {
		return nil
	}

	if len(match) > 0 {
		return match[1:]
	}

	return match
}

// InlineMatches returns all matches of the pattern in the message's text or caption.
// The returned matches are trimmed and an empty slice is returned if there are none.
func (m Matcher) InlineMatches(messageIn telegramclient.WebhookMessageStruct) []string {
	matches := m.regexp.FindAllString(messageIn.TextOrCaption(), -1)
	if matches == nil {
		return []string{}
	}

	for i, match := range matches {
		matches[i] = strings.TrimSpace(match)
	}

	return matches
}

// HandleError logs an error that occurred while processing a message for this matcher.
func (m Matcher) HandleError(_ telegramclient.WebhookMessageStruct, identifier string, err error) {
	m.log.Error(identifier, err.Error())
}

// WithCustomConfigType is a generic matcher wrapper that carries a typed config T and
// exposes it via Config() without requiring each matcher to re-implement it.
//
// It embeds the base Matcher for all core behavior.
type WithCustomConfigType[T any] struct {
	Matcher

	cfg T
}

// MakeMatcherWithCustomConfigType constructs a new MatcherWithCustomConfig[T] matcher from the given base inputs
// and the typed config value.
func MakeMatcherWithCustomConfigType[T any](identifier string, pattern *regexp.Regexp, help []HelpStruct, cfg T) WithCustomConfigType[T] {
	m := WithCustomConfigType[T]{
		Matcher: MakeMatcher(identifier, pattern, help),
	}

	return m.WithTypedConfig(cfg)
}

// WithTypedConfig returns a copy with the typed config applied as well as wiring
// the base Matcher config if the typed config embeds matcher.Config.
func (m WithCustomConfigType[T]) WithTypedConfig(cfg T) WithCustomConfigType[T] {
	m.cfg = cfg

	// If T embeds matcher.Config, also pass its address to the base matcher so
	// IsEnabled and similar helpers work. We use a type assertion to an
	// interface that exposes the embedded Config by pointer.
	if c, ok := any(&cfg).(interface{ GetEmbeddedMatcherConfigPtr() *Config }); ok {
		m.Matcher = m.WithConfig(c.GetEmbeddedMatcherConfigPtr())
	}

	return m
}

// Config returns the typed configuration value.
func (m WithCustomConfigType[T]) Config() T {
	return m.cfg
}
