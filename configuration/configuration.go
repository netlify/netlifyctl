package configuration

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

const DefaultConfigFileName = "netlify.toml"

type Settings struct {
	ID   string
	Path string `toml:"path,omitempty"`
}

type Context struct {
	Publish   string
	Functions string
}

type Redirect struct {
	// TODO: This implementation is incomplete, but it allow us to test a few things already.
	// This program doesn't really use the redirects for anything.
	From    string
	To      string
	Status  int
	Force   bool
	Headers map[string]string
}

type Configuration struct {
	Settings  Settings
	Build     Context
	Redirects []Redirect
	root      string
	filePath  string
}

func (c Configuration) Root() string {
	return c.root
}

func (c Configuration) ExistConfFile() bool {
	_, err := os.Stat(c.confFilePath())
	return err == nil
}

// CopyConfigFile will copy over the toml file if there isn't one already in the
// publish path. That means a user can create the file there and we won't override it
func (c Configuration) CopyConfigFile(pubPath string) (string, error) {
	dest := filepath.Join(pubPath, c.filePath)
	if _, err := os.Stat(dest); err == nil {
		return "", nil // file exists, don't overwrite
	}

	f, err := os.Open(c.confFilePath())
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	w, err := os.Create(dest)
	if err != nil {
		return "", err
	}
	defer func() { _ = w.Close() }()

	_, err = io.Copy(w, f)
	return dest, err
}

func (c Configuration) confFilePath() string {
	return filepath.Join(c.root, c.filePath)
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
	c.filePath = configFile

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
