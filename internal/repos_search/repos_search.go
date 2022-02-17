package repos_search

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func CrawlPath(srcPath string) (map[string][]string, error) {
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
				if v.Name() == ".git" || v.Name() == "go.mod" {
					fp := filepath.Base(path)

					if _, ok := pathsMap[fp]; !ok {
						pathsMap[fp] = []string{path}
					} else {
						pathsMap[fp] = append(pathsMap[fp], path)
					}
					break
				}
			}
		}
		return nil
	})
}
