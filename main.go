package main

import (
	"flag"
	"fmt"
	"go/build"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/BUCKU/gjg/config"
)

const (
	srcDir = "/src"
)

func main() {
	reinit := flag.Bool("r", false, "reinit gjg")
	flag.Parse()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	gopath := build.Default.GOPATH
	if len(gopath) <= 1 {
		log.Fatal("$GOPATH is not set")
	}

	cfg, err := config.ProcessConfig(homeDir, *reinit)
	if err != nil {
		log.Fatalf("cannot process config: %s", err)
	}

	cmdGoland := exec.Command(cfg.GolandPath)

	if len(os.Args) > 1 {
		paths, err := crawlSrcPath(filepath.Join(gopath, srcDir))
		if err != nil {
			log.Fatal(err)
		}

		repo := os.Args[len(os.Args)-1]

		if _, ok := paths[repo]; !ok {
			log.Fatalf("cannot find repo: '%s'", repo)
		}

		cmdGoland = exec.Command(cfg.GolandPath, paths[repo][0])
	}

	fmt.Printf("opening: %s\n", cmdGoland)
	err = cmdGoland.Start()
	if err != nil {
		fmt.Println(err)
	}
}

func crawlSrcPath(srcPath string) (map[string][]string, error) {
	pathsMap := make(map[string][]string)

	hosts, err := os.ReadDir(srcPath)
	if err != nil {
		return nil, err
	}

	for _, v := range hosts {
		err := detectGitRepos(filepath.Join(srcPath, v.Name()), pathsMap)
		if err != nil {
			return nil, err
		}
	}

	return pathsMap, nil
}

func detectGitRepos(path string, pathsMap map[string][]string) error {
	return filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		var dif []os.DirEntry

		stat, err := os.Stat(path)
		if err != nil {
			fmt.Printf("failed to get fsstat of path: %s, %s", path, err.Error())
			return nil
		}

		if stat.IsDir() {
			if dif, err = os.ReadDir(path); err != nil {
				fmt.Printf("failed to read dir %s\n", path)
			}
			for _, v := range dif {
				if v.Name() == ".git" {
					fp := filepath.Base(path)

					if _, ok := pathsMap[fp]; !ok {
						pathsMap[fp] = []string{path}
					} else {
						pathsMap[fp] = append(pathsMap[v.Name()], path)
					}
					break
				}
			}
		}
		return nil
	})
}
