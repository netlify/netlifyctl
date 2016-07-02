package configuration

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const globalConfigFileName = "netlify.toml"

type Settings struct {
	ID   string
	root string
}

func (s Settings) Root() string {
	return s.root
}

type Configuration struct {
	Settings Settings
}

func Load() (*Configuration, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	single := filepath.Join(pwd, globalConfigFileName)
	fi, err := os.Stat(single)
	if err != nil {
		return nil, err
	}

	if fi.IsDir() {
		return nil, fmt.Errorf("%s cannot be a directory", globalConfigFileName)
	}

	var c Configuration
	if _, err := toml.DecodeFile(single, &c); err != nil {
		return nil, err
	}
	c.Settings.root = pwd

	return &c, nil
}
