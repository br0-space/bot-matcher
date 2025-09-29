package main

import (
	logger "github.com/br0-space/bot-logger"
	matcher "github.com/br0-space/bot-matcher"
	"github.com/br0-space/bot-matcher/examples/configurable"
	"github.com/br0-space/bot-matcher/examples/null"
	"github.com/br0-space/bot-matcher/examples/ping"
	telegramclient "github.com/br0-space/bot-telegramclient"
	"github.com/spf13/pflag"
)

// main initializes a matcher registry and registers each example matcher.
func main() {
	pflag.Bool("verbose", false, "enable verbose (debug) logging")
	pflag.Bool("quiet", false, "only log errors")
	pflag.Parse()

	log := logger.New()
	telegram := telegramclient.NewMockClient()

	log.Info("Starting matcher registry example...")
	log.Debug("Creating matcher registry...")

	r := matcher.NewRegistry(log, telegram)

	// Register example matchers.
	r.Register(configurable.MakeMatcher())
	r.Register(ping.MakeMatcher())
	r.Register(null.MakeMatcher())
}
