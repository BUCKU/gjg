package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gjg/config"
)

var (
	golandPath = "/home/brix/.local/share/JetBrains/Toolbox/apps/Goland/ch-0/212.5457.54/bin/goland.sh"
)

const (
	homePrefix = "HOME"
	gopathPrefix = "GOPATH"
)

func main() {
	reinit := flag.Bool("r", false, "reinit gjg")
	flag.Parse()

	args := os.Args
	if len(args) < 2 {
		fmt.Println("please provide project name")
	}

	envs := os.Environ()
	var home string
	for _, v := range envs {
		if strings.HasPrefix(v, homePrefix) {
			home = strings.TrimLeft(v, homePrefix+"=")
			break
		}
	}

	cfg, err := config.ProcessConfig(home, *reinit)
	if err != nil {
		err = fmt.Errorf("cant initialize gjg, %s", err)
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var gopath string
	for _, v := range envs {
		if strings.HasPrefix(v, gopathPrefix) {
			gopath = strings.TrimLeft(v, gopathPrefix+"=")
			break
		}
	}

	if len(gopath) <= 1 {
		fmt.Println("no gopath in envs")
		os.Exit(1)
	}

	sourcesPath := gopath + "/src"
	fmt.Printf("use path for scan '%s'\n", sourcesPath)

	paths, err := parseGoPath(sourcesPath)
	if err != nil {
		fmt.Printf("error at src dir crowling. %s", err.Error())
	}

	// TODO: verbose
	// for i, v := range paths {
	// 	fmt.Printf("%s: %v\n", i, v)
	// }
	// fmt.Printf("length %d\n", len(paths))

	repo := os.Args[len(os.Args)-1]

	if _, ok := paths[repo]; !ok {
		fmt.Printf("repo didnt exist or not being parsed: '%s'", repo)
		os.Exit(1)
	}
	fmt.Printf("starting: %s %s\n\n", cfg.GolandPath, paths[repo][0])
	cmdGoland := exec.Command(cfg.GolandPath, paths[repo][0])
	err = cmdGoland.Run()
	if err != nil {
		fmt.Println(err)
	}

}



func parseGoPath(path string) (map[string][]string, error) {
	pathsMap := make(map[string][]string)

	hosts, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, v := range hosts {
		err := detectGitRepos(path+"/"+v.Name(), pathsMap)
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
