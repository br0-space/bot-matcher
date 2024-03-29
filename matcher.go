package matcher

import (
	"fmt"
	"regexp"
	"strings"

	logger "github.com/br0-space/bot-logger"
	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/spf13/viper"
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

func (m Matcher) WithConfig(cfg *Config) Matcher {
	m.cfg = cfg

	return m
}

func (m Matcher) IsEnabled() bool {
	if m.cfg == nil || m.cfg.enabled == nil {
		return true
	}

	return *m.cfg.enabled
}

func (m Matcher) Identifier() string {
	return m.identifier
}

func (m Matcher) Help() []HelpStruct {
	return m.help
}

func (m Matcher) DoesMatch(messageIn telegramclient.WebhookMessageStruct) bool {
	return m.regexp.MatchString(messageIn.TextOrCaption())
}

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

func (m Matcher) HandleError(_ telegramclient.WebhookMessageStruct, identifier string, err error) {
	m.log.Error(identifier, err.Error())
}

func LoadMatcherConfig(identifier string, cfg interface{}) {
	v := viper.New()

	v.SetConfigFile(fmt.Sprintf("config/%s.yaml", identifier))

	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}

	if err := v.Unmarshal(cfg); err != nil {
		panic(err)
	}
}
