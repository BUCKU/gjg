package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	cfgDir  = ".local/share/gjg"
	cfgFile = "gjg.conf"
)

type Config struct {
	GolandPath string `json:"goland_path"`
}

func ProcessConfig(path string, reinit bool) (cfg *Config, err error) {
	cfgContents, err := os.ReadFile(filepath.Join(path, cfgDir, cfgFile))
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		cfgContents = []byte("{}") // TODO: костыль, сделать нормально сразу чтоб было
	}

	cfg = &Config{}
	err = json.Unmarshal(cfgContents, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse cfg. %s", err)
	}

	if reinit {
		initializeConfig(cfg)
	}

	if cfg.GolandPath == "" {
		// must init
		initializeConfig(cfg)
	}

	_, err = os.Stat(cfg.GolandPath)
	if err != nil {
		fmt.Printf("no goland.sh detected. must re-init.")
		// must init
		initializeConfig(cfg)
	}

	if err = SaveConfig(path, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func SaveConfig(path string, c *Config) error {
	bytes, err := json.MarshalIndent(c, " ", "    ")
	if err != nil {
		return err
	}

	// create config directory if not exists
	cfgPath := filepath.Join(path, cfgDir)

	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		if err = os.Mkdir(cfgPath, os.ModePerm); err != nil {
			return err
		}
	}

	if err = os.WriteFile(filepath.Join(path, cfgDir, cfgFile), bytes, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func initializeConfig(cfg *Config) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to recognize homedir")
	}

	cmd := exec.Command("find", homeDir, "-name", "goland.sh")
	b, err := cmd.Output()
	if len(b) == 0 {
		log.Fatal().Err(errors.New("find didn't find anything")).Str("search", "goland.sh").Send()
	}
	if err != nil {
		log.Error().Err(err).Str("search", "goland.sh").Msg("errors finding goland.sh")
	}
	output := string(b)
	paths := strings.Split(output, "\n")

	var golandPath string
	if len(paths) > 1 {
		fmt.Println("choose your goland")
		for i, p := range paths {
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
		if choose > len(paths)-1 {
			fmt.Println("dont try to trick me nigger")
		}
		golandPath = strings.Trim(paths[choose], "\n")
	}
	cfg.GolandPath = golandPath
}
