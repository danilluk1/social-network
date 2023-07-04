package main

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	var baseDir string
	flag.StringVar(&baseDir, "basedir", "../../apps", "Path to app folder")
	flag.Parse()

	microservicesFolders, err := getFolders(baseDir)
	if err != nil {
		panic(err)
	}
	for _, mf := range microservicesFolders {
		runCommand(mf)
	}
}

func getFolders(dirPath string) ([]string, error) {
	var folders []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			folders = append(folders, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return folders, nil
}

func runCommand(folder string) error {
	cmd := exec.Command("task", "migrationup")
	cmd.Dir = folder

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
