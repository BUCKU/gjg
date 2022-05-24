package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
