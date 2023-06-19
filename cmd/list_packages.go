package cmd

import (
	"bytes"
	"fmt"
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

func goList(args ...string) ([]byte, error) {
	cmd := exec.Command("go", append([]string{"list"}, args...)...)
	cmd.Stderr = os.Stderr
	return cmd.Output()
}

func listPackageNames(modDir string) ([]string, error) {
	data, err := goList(fmt.Sprintf("%s%s...", modDir, string(os.PathSeparator)))
	if err != nil {
		return nil, err
	}

	var packageNames []string
	for _, line := range bytes.Split(data, []byte{'\n'}) {
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

	data, err := goList(
		"-f", "{{ .Dir }}",
		fmt.Sprintf("%s%s...", modDir, string(os.PathSeparator)),
	)
	if err != nil {
		return nil, err
	}

	var packagePaths []string
	for _, line := range bytes.Split(data, []byte{'\n'}) {
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
