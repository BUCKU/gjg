//go:build darwin

package config

import (
	"os"

	"github.com/rs/zerolog/log"
)

const (
	cfgDir  = ".config/gjg"
	cfgFile = "gjg.conf"
)

func initializeConfig(cfg *Config) {
	golandPath := "/Applications/GoLand.app/Contents/MacOS/goland"
	if _, err := os.Stat(golandPath); os.IsNotExist(err) {
		log.Fatal().Msgf("didn't find goland in %v", golandPath)
	}

	cfg.GolandPath = golandPath
}
