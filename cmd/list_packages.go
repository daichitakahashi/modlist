package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func expand(ss []string, fn func(s string) ([]string, error)) ([]string, error) {
	var result []string
	uniq := map[string]bool{}
loop:
	for _, s := range ss {
		r, err := fn(s)
		if err != nil {
			return nil, err
		}
		for _, rr := range r {
			if uniq[rr] {
				continue loop
			}
			uniq[rr] = true
		}
		if len(r) > 0 {
			result = append(result, r...)
		}
	}
	return result, nil
}

func listPackageNames(modDir string) ([]string, error) {
	result, err := exec.Command("go", "list", fmt.Sprintf("%s%s...", modDir, string(os.PathSeparator))).CombinedOutput()
	if err != nil {
		log.Println(string(result))
		return nil, err
	}

	var packageNames []string
	for _, line := range bytes.Split(result, []byte{'\n'}) {
		line = bytes.TrimSpace(line)
		if len(line) > 0 {
			packageNames = append(packageNames, string(line))
		}
	}

	return packageNames, nil
}

func listPackagePaths(modDir string) ([]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	result, err := exec.Command("go", "list",
		"-f", "{{ .Dir }}",
		fmt.Sprintf("%s%s...", modDir, string(os.PathSeparator)),
	).CombinedOutput()
	if err != nil {
		log.Println(string(result))
		return nil, err
	}

	var packagePaths []string
	for _, line := range bytes.Split(result, []byte{'\n'}) {
		line = bytes.TrimSpace(line)
		if len(line) > 0 {
			path, err := filepath.Rel(wd, string(line))
			if err != nil {
				return nil, err
			}
			if !strings.HasPrefix(path, ".") {
				path = fmt.Sprintf(".%s%s", string(os.PathSeparator), path)
			}
			packagePaths = append(packagePaths, path)
		}
	}

	return packagePaths, nil
}
