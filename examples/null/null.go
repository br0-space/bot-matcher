// Package null provides an example matcher that deliberately never matches any input.
// It is useful for demonstrating how a matcher is structured without affecting behavior.
package null

import (
	"errors"
	"regexp"

	matcher "github.com/br0-space/bot-matcher"
	telegramclient "github.com/br0-space/bot-telegramclient"
)

// identifier is the unique name of this matcher.
var identifier = "null"

// pattern is a regular expression that never matches. Using an impossible zero-width assertion
// combination (word boundary AND non-word boundary) ensures it can never match in Go's RE2 engine.
var pattern = regexp.MustCompile(`\b\B`)

// help describes how to use this matcher and can be rendered in help messages.
var help = []matcher.HelpStruct{{
	Description: `Example; never matches`,
}}

// template is the reply text that would be returned if Process were ever reached after a match.
const template = `You should never get this response.`

// Matcher is an example matcher that embeds the base matcher.
// It never matches any message and exists only for documentation/testing purposes.
type Matcher struct {
	matcher.Matcher
}

// MakeMatcher constructs a new null.Matcher with its identifier, pattern, and help configured.
func MakeMatcher() Matcher {
	return Matcher{
		Matcher: matcher.MakeMatcher(identifier, pattern, help),
	}
}

// Process implements the matcher processing pipeline.
// Since the null matcher never matches, this method always returns an error and no replies.
func (m Matcher) Process(messageIn telegramclient.WebhookMessageStruct) ([]telegramclient.MessageStruct, error) {
	if !m.DoesMatch(messageIn) {
		return nil, errors.New("message does not match")
	}

	return []telegramclient.MessageStruct{
		telegramclient.Reply(template, messageIn.ID),
	}, nil
}
