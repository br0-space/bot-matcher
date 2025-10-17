package matcher

import (
	"fmt"
	"path/filepath"
	"strconv"

	logger "github.com/br0-space/bot-logger"
	"github.com/spf13/viper"
)

// LoadMatcherConfig loads configurations for a matcher per chat.
// It returns a map keyed by chatID (int64) to a value of type T.
// The fallback config is read from config/{identifier}.yml and stored under key 0.
// Additionally, all files matching config/{chatID}/{identifier}.yml are read and stored under their chatID key.
// Panics if reading or unmarshalling any relevant file fails.
func LoadMatcherConfig[T any](identifier string) map[int64]T {
	log := logger.New()
	out := make(map[int64]T)

	log.Debugf("%s: requested to load matcher config", identifier)

	// Per chat configs in config/{chatID}/{identifier}.yml
	pattern := fmt.Sprintf("config/*/%s.yml", identifier)

	matches, _ := filepath.Glob(pattern)
	for _, p := range matches {
		// extract chatID from the parent directory name
		dir := filepath.Base(filepath.Dir(p))

		chatID, err := strconv.ParseInt(dir, 10, 64)
		if err != nil {
			continue // skip non-numeric directories
		}

		log.Debugf("%s: reading per-chat config for chatID=%d from %s", identifier, chatID, p)

		v2 := viper.New()
		v2.SetConfigFile(p)

		if err := v2.ReadInConfig(); err != nil {
			log.Debugf("%s: failed to read per-chat config %s: %v", identifier, p, err)
			panic(err)
		}

		var cfg T
		if err := v2.Unmarshal(&cfg); err != nil {
			log.Debugf("%s: failed to unmarshal per-chat config %s: %v", identifier, p, err)
			panic(err)
		}

		out[chatID] = cfg
	}

	// Fallback config at key 0
	fallbackPath := fmt.Sprintf("config/%s.yml", identifier)
	log.Debugf("%s: reading fallback config: %s", identifier, fallbackPath)

	v := viper.New()
	v.SetConfigFile(fallbackPath)

	if err := v.ReadInConfig(); err != nil {
		log.Debugf("%s: failed to read fallback config %s: %v", identifier, fallbackPath, err)
		panic(err)
	}

	var base T
	if err := v.Unmarshal(&base); err != nil {
		log.Debugf("%s: failed to unmarshal fallback config %s: %v", identifier, fallbackPath, err)
		panic(err)
	}

	out[0] = base

	return out
}
