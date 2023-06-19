package golangci

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"

	"github.com/BurntSushi/toml"
	"github.com/goccy/go-yaml"
	"github.com/golangci/golangci-lint/pkg/fsutils"
	"github.com/golangci/golangci-lint/pkg/packages"
)

var possibleConfigFiles = []string{
	".golangci.yml",
	".golangci.yaml",
	".golangci.toml",
	".golangci.json",
}

type Config struct {
	Run struct {
		SkipDirs           []string `json:"skip-dirs" yaml:"skip-dirs" toml:"skip-dirs"`
		SkipDirsUseDefault *bool    `json:"skip-dirs-use-default" yaml:"skip-dirs-use-default" toml:"skip-dirs-use-default"`
	} `json:"run" yaml:"run" toml:"run"`
}

func ReadConfig() (*Config, error) {
	filename, data, err := readConfigFile()
	if err != nil {
		return nil, err
	}

	var c Config
	switch filepath.Ext(filename) {
	case ".yml", ".yaml":
		err = yaml.Unmarshal(data, &c)
	case ".json":
		err = json.Unmarshal(data, &c)
	case ".toml":
		err = toml.Unmarshal(data, &c)
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func readConfigFile() (string, []byte, error) {
	var (
		data []byte
		err  error
	)
	for _, f := range possibleConfigFiles {
		data, err = os.ReadFile(f)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return "", nil, err
		}
		return f, data, nil
	}
	return "", nil, err
}

func (c *Config) SkipDirectories() ([]*regexp.Regexp, error) {
	var skipDefault = true
	if c.Run.SkipDirsUseDefault != nil {
		skipDefault = *c.Run.SkipDirsUseDefault
	}

	patterns := make([]string, 0, len(c.Run.SkipDirs))
	for _, p := range c.Run.SkipDirs {
		if p == "" {
			continue
		}
		patterns = append(patterns, fsutils.NormalizePathInRegex(p))
	}
	if skipDefault {
		patterns = append(patterns, packages.StdExcludeDirRegexps...)
	}

	rxs := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		rx, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		rxs = append(rxs, rx)
	}

	return rxs, nil
}
