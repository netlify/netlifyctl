package configuration

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/Sirupsen/logrus"
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

	var c Configuration
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

	return &c, nil
}

func Save(conf *Configuration) error {
	single := filepath.Join(conf.Settings.Root(), globalConfigFileName)
	f, err := os.OpenFile(single, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	return toml.NewEncoder(f).Encode(conf)
}
