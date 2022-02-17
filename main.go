package main

import (
	"flag"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/BUCKU/gjg/config"
	"github.com/BUCKU/gjg/internal/consts"
	"github.com/BUCKU/gjg/internal/repos_search"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	verbose := flag.Bool("v", false, "verbose")

	reinit := flag.Bool("r", false, "reinit gjg")
	flag.Parse()

	if *verbose {
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Err(err).Msg("get user home directory")
	}

	gopath := build.Default.GOPATH
	if len(gopath) <= 1 {
		log.Fatal().Msg("$GOPATH is not set")
	}

	cfg, err := config.ProcessConfig(homeDir, *reinit)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot process config")
	}

	cmdGoland := exec.Command(cfg.GolandPath)

	if len(os.Args) > 1 {
		goSrcPath := filepath.Join(gopath, consts.GoSrcDir)
		reposPaths, err := repos_search.CrawlPath(goSrcPath)
		if err != nil {
			log.Fatal().Err(err).Str("path", goSrcPath).Msg("crawling go sources failed.")
		}

		repoToFind := os.Args[len(os.Args)-1]

		var reposWithName []string
		var ok bool
		if reposWithName, ok = reposPaths[repoToFind]; !ok {
			log.Fatal().Str("name", repoToFind).Msg("search repo failed. repo dir must contain either 'go.mod' or '.git' entry.")
		}

		repoOfChoice := reposPaths[repoToFind][0]
		if len(reposWithName) > 1 {
			fmt.Printf("There is different repos with the same name, choose one to open:\n")
			for i, v := range reposWithName {
				trimmedPath := v[strings.Index(v, "/src/")+len("/src/"):]
				fmt.Printf("[%d] %s\n", i, trimmedPath)
			}
			choose := 0
			_, err := fmt.Scanln(&choose)
			if err != nil {
				log.Fatal().Err(err).Msg("scan user input")
			}
			if choose > len(reposWithName) {
				log.Warn().Int("user input", choose).Msg("no repo with such index, using first repo in list")
				choose = 0
			}
			repoOfChoice = reposWithName[choose]
		}

		cmdGoland = exec.Command(cfg.GolandPath, repoOfChoice)
		log.Info().Str("command", cmdGoland.String()).Str("args", repoOfChoice).Msg("execute")
		err = cmdGoland.Start()
		if err != nil {
			log.Error().Err(err).Str("command", cmdGoland.String()).Str("args", repoOfChoice).Msg("execute")
		}
	}
}
