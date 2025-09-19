# bot-matcher

Small Go library for building Telegram bots by composing independent “matchers.”
A matcher declares what incoming messages it responds to (via regex) and returns zero or more outgoing messages. A central Registry fans incoming messages out to all registered matchers and sends replies via a Telegram client.

- Lightweight building blocks for commands and inline patterns
- Concurrent processing of multiple matchers per message
- Pluggable logger and Telegram client interfaces
- Optional per-matcher configuration loader with Viper

## Status

[![Build](https://github.com/br0-space/bot-matcher/actions/workflows/build.yml/badge.svg?branch=main)](https://github.com/br0-space/bot-matcher/actions/workflows/build.yml)
[![Test](https://github.com/br0-space/bot-matcher/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/br0-space/bot-matcher/actions/workflows/test.yml)
[![Lint](https://github.com/br0-space/bot-matcher/actions/workflows/lint.yml/badge.svg?branch=main)](https://github.com/br0-space/bot-matcher/actions/workflows/lint.yml)
[![Staticcheck](https://github.com/br0-space/bot-matcher/actions/workflows/staticcheck.yml/badge.svg?branch=main)](https://github.com/br0-space/bot-matcher/actions/workflows/staticcheck.yml)
[![Vet](https://github.com/br0-space/bot-matcher/actions/workflows/vet.yml/badge.svg?branch=main)](https://github.com/br0-space/bot-matcher/actions/workflows/vet.yml)
[![CodeQL](https://github.com/br0-space/bot-matcher/actions/workflows/codeql-analysis.yml/badge.svg?branch=main)](https://github.com/br0-space/bot-matcher/actions/workflows/codeql-analysis.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/br0-space/bot-matcher.svg)](https://pkg.go.dev/github.com/br0-space/bot-matcher)
[![Go Report Card](https://goreportcard.com/badge/github.com/br0-space/bot-matcher)](https://goreportcard.com/report/github.com/br0-space/bot-matcher)

## Installation

Requires Go 1.21+ (tested in CI on latest Go versions).

```
go get github.com/br0-space/bot-matcher
```

## Getting started

bot-matcher revolves around two concepts:

- Matcher: implements the matcher.Interface and contains logic to decide if a message is relevant and how to respond.
- Registry: holds a set of matchers and processes incoming Telegram webhook messages concurrently, sending responses with a Telegram client.

### Define a matcher

You can build matchers by composing the provided Matcher type with your own Process implementation.

```go
package hello

import (
    "regexp"

    matcher "github.com/br0-space/bot-matcher"
    telegramclient "github.com/br0-space/bot-telegramclient"
)

type HelloMatcher struct {
    matcher.Matcher
}

func New() HelloMatcher {
    return matcher.MakeMatcher(
        "hello",
        regexp.MustCompile(`(?i)^/hello\b`),
        []matcher.HelpStruct{{
            Command:     "/hello",
            Description: "Say hello",
            Usage:       "/hello",
            Example:     "/hello",
        }},
    )
}

// Process is called when DoesMatch returns true.
func (h HelloMatcher) Process(messageIn telegramclient.WebhookMessageStruct) ([]telegramclient.MessageStruct, error) {
    reply := telegramclient.PlainReply("Hello!", messageIn.ID)
    return []telegramclient.MessageStruct{reply}, nil
}
```

Notes:
- matcher.Matcher already provides helpful defaults for Identifier, DoesMatch, Help, InlineMatches, CommandMatch, IsEnabled and HandleError.
- You only need to implement Process when composing as shown above.

### Register and process messages

```go
package main

import (
    logger "github.com/br0-space/bot-logger"
    matcher "github.com/br0-space/bot-matcher"
    telegramclient "github.com/br0-space/bot-telegramclient"
)

func main() {
    log := logger.New()

    // Create your Telegram client (implementation depends on bot-telegramclient).
    var tg telegramclient.ClientInterface = telegramclient.New(/* ... */)

    reg := matcher.NewRegistry(log, tg)

    // Register matchers
    reg.Register(hello.New())

    // In your webhook handler:
    var incoming telegramclient.WebhookMessageStruct // fill from Telegram update
    reg.Process(incoming)
}
```

### Optional configuration per matcher

matcher.LoadMatcherConfig lets you unmarshal YAML into your own config struct. A small helper Config type is provided with an optional Enabled flag.

```go
// Inside your matcher package

type MyConfig struct {
    Enabled *bool `yaml:"enabled"`
    // your fields...
}

func (h HelloMatcher) WithConfig(cfg *matcher.Config) HelloMatcher {
    return h.Matcher.WithConfig(cfg).(HelloMatcher)
}

// At startup
var cfg MyConfig
matcher.LoadMatcherConfig("hello", &cfg)
```

This reads config/hello.yaml using Viper and unmarshals it into your struct.

## Concepts and API

- Interface: the contract for matchers (Identifier, Help, DoesMatch, Process, etc.). See type.go.
- Registry: coordinates concurrent execution of all registered matchers per message and sends outputs using the injected Telegram client.
- Error handling: if Process returns an error, Registry logs it and replies in chat with a markdown-formatted error message referencing the matcher.

## Development

- Run linters and tests via the provided GitHub Actions workflows, or locally using your preferred tooling.
- The library depends on:
  - github.com/br0-space/bot-logger
  - github.com/br0-space/bot-telegramclient
  - github.com/spf13/viper (for optional config loader)

## License

MIT - See [LICENSE](LICENSE).
