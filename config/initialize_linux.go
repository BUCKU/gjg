//go:build linux

package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	cfgDir  = ".local/share/gjg"
	cfgFile = "gjg.conf"
)

func initializeConfig(cfg *Config) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to recognize homedir")
	}

	standardPaths := []string{homeDir, "/snap"}
	var golandPaths []string

	for _, v := range standardPaths {
		cmd := exec.Command("find", v, "-name", "goland.sh")
		b, err := cmd.Output()
		if len(b) == 0 {
			log.Warn().Err(errors.New("find didn't find anything")).Str("search", "goland.sh").Str("path", v).Send()
			continue
		}
		if err != nil && len(b) == 0 {
			log.Error().Err(err).Str("search", "goland.sh").Msg("errors finding goland.sh")
		}
		output := string(b)
		paths := strings.Split(output, "\n")
		golandPaths = append(golandPaths, paths...)
	}

	if len(golandPaths) == 0 {
		log.Fatal().Msg("didn't find any goland.sh on that computer")
	}

	var golandPath string
	if len(golandPaths) > 1 {
		fmt.Println("choose your goland")
		for i, p := range golandPaths {
			if p == "" {
				i--
				continue
			}
			fmt.Printf("[%d] %s\n", i, p)
		}
		choose := 0
		_, err = fmt.Scan(&choose)
		if err != nil {
			fmt.Println("dont try to trick me nigger")
		}
		if choose > len(golandPaths)-1 {
			fmt.Println("dont try to trick me nigger")
		}
		golandPath = strings.Trim(golandPaths[choose], "\n")
	}
	cfg.GolandPath = golandPath
}
