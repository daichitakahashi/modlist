package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

func listModules(dir bool) ([]string, error) {
	result, err := exec.Command("go", "env", "GOWORK").CombinedOutput()
	if err != nil {
		log.Println(string(result))
		return nil, err
	}

	var modulePaths []string
	workPath := string(bytes.TrimSpace(result))
	if workPath == "" {
		modulePaths = append(modulePaths, ".")
	} else {
		modulePaths, err = loadWork(workPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load module paths from %s: %w", workPath, err)
		}
	}
	if dir {
		return modulePaths, nil
	}

	modules := make([]string, 0, len(modulePaths))
	for _, modDir := range modulePaths {
		name, found, err := loadModuleName(modDir)
		if err != nil {
			return nil, fmt.Errorf("failed to load module name: path=%q: %w", modDir, err)
		}
		if found {
			modules = append(modules, name)
		}
	}

	return modules, nil
}

func loadWork(workFile string) ([]string, error) {
	data, err := os.ReadFile(workFile)
	if err != nil {
		return nil, err
	}
	work, err := modfile.ParseWork(workFile, data, nil)
	if err != nil {
		return nil, err
	}
	modules := make([]string, 0, len(work.Use))
	for _, u := range work.Use {
		modules = append(modules, u.Path)
	}
	return modules, nil
}

func loadModuleName(modDir string) (string, bool, error) {
	modFile := filepath.Join(modDir, "go.mod")
	data, err := os.ReadFile(modFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", false, nil
		}
		return "", false, err
	}
	mod, err := modfile.Parse(modFile, data, nil)
	if err != nil {
		return "", false, err
	}
	return mod.Module.Mod.String(), true, nil
}
