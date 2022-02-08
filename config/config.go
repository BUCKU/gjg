package config

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

var (
	homePrefix = "HOME"
	cfgPath    = ".local/share/gjg/gjg.conf"
)

type Config struct {
	GolandPath string `json:"goland_path"`
}

func ProcessConfig(path string, reinit bool) (cfg *Config, err error) {
	cfgContents, err := os.ReadFile(path + "/" + cfgPath)
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

	SaveConfig(cfg)
	return cfg, nil
}

func SaveConfig(c *Config) {
	bytes, err := json.MarshalIndent(c, " ", "    ")
	if err != nil {
		fmt.Errorf("failed to marshal config. %s", err)
	}

	err = os.WriteFile(cfgPath, bytes, 0664)
	if err != nil {
		fmt.Errorf("failed to marshal config. %s", err)
	}
}

func initializeConfig(cfg *Config) {
	envs := os.Environ()
	var home string
	for _, v := range envs {
		if strings.HasPrefix(v, homePrefix) {
			home = strings.TrimLeft(v, homePrefix+"=")
			break
		}
	}
	cmd := exec.Command("find", home, "-name", "goland.sh")
	b, err := cmd.Output()
	if err != nil {
		fmt.Printf("failed to find goland. %s\n", err.Error())
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
