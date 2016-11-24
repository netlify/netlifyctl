package configuration

import (
	"fmt"
	"os"
	"path/filepath"

	"strings"

	"github.com/BurntSushi/toml"
	"github.com/Sirupsen/logrus"
)

const globalConfigFileName = "netlify.toml"

type Settings struct {
	ID   string
	Path string
	root string
}

type Configuration struct {
	Settings Settings
}

func Load() (*Configuration, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var c Configuration
	c.Settings.Path = pwd
	c.Settings.root = pwd

	single := filepath.Join(pwd, globalConfigFileName)
	fi, err := os.Stat(single)
	if err != nil {
		if os.IsNotExist(err) {
			return &c, nil
		}
		return nil, err
	}

	if fi.IsDir() {
		return nil, fmt.Errorf("%s cannot be a directory", globalConfigFileName)
	}

	if _, err := toml.DecodeFile(single, &c); err != nil {
		return nil, err
	}

	logrus.Debugf("Parsed configuration: %+v", c)

	if c.Settings.Path != "" {
		if !strings.HasPrefix(c.Settings.Path, "/") {
			c.Settings.Path = filepath.Join(pwd, c.Settings.Path)
			logrus.Debugf("Relative path detected, going to deploy: '%s'", c.Settings.Path)
		}
	}

	return &c, nil
}

func Save(conf *Configuration) error {
	single := filepath.Join(conf.Settings.root, globalConfigFileName)
	f, err := os.OpenFile(single, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	return toml.NewEncoder(f).Encode(conf)
}
