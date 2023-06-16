package cmd

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
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

func listPackages(modDir string) ([]string, error) {
	pkg := !strings.HasPrefix(modDir, ".")

	result, err := exec.Command("go", "list", fmt.Sprintf("%s%s...", modDir, string(os.PathSeparator))).CombinedOutput()
	if err != nil {
		log.Println(string(result))
		return nil, err
	}

	var packagePaths []string
	for _, line := range bytes.Split(result, []byte{'\n'}) {
		line = bytes.TrimSpace(line)
		if len(line) > 0 {
			packagePaths = append(packagePaths, string(line))
		}
	}
	if pkg {
		return packagePaths, nil
	}

	log.Println(modDir, packagePaths)

	return nil, nil
}
