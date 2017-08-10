package configuration

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

type Settings struct {
	ID   string
	Path string
}

type Configuration struct {
	Settings Settings
	root     string
}

func (c Configuration) Root() string {
	return c.root
}

func Exist(configFile string) bool {
	pwd, err := os.Getwd()
	if err != nil {
		return false
	}

	single := filepath.Join(pwd, configFile)
	_, err = os.Stat(single)
	return err == nil
}

func Load(configFile string) (*Configuration, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	var c Configuration
	c.root = pwd

	fi, err := os.Stat(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &c, nil
		}
		return nil, err
	}

	if fi.IsDir() {
		return nil, fmt.Errorf("%s cannot be a directory", configFile)
	}

	if _, err := toml.DecodeFile(configFile, &c); err != nil {
		return nil, err
	}

	logrus.Debugf("Parsed configuration: %+v", c)

	return &c, nil
}

func Save(configFile string, conf *Configuration) error {
	single := filepath.Join(conf.root, configFile)
	f, err := os.OpenFile(single, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	return toml.NewEncoder(f).Encode(conf)
}
