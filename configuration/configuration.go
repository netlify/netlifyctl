package configuration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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

// SetupViper will setup the env prefixes and bind all the flags together
func SetupViper(persistentFlags, localFlags *pflag.FlagSet) error {
	viper.SetEnvPrefix("NETLIFY")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv()

	if persistentFlags != nil {
		err := viper.BindPFlags(persistentFlags)
		if err != nil {
			return err
		}
	}
	if localFlags != nil {
		err := viper.BindPFlags(localFlags)
		if err != nil {
			return err
		}
	}

	home := os.Getenv("HOME")
	legacyConfigPath := filepath.Join(home, ".netlify", "config")
	if stat, err := os.Stat(legacyConfigPath); err == nil {
		if !stat.IsDir() {
			viper.SetConfigFile(legacyConfigPath)
			viper.SetConfigType("json")
		}
	} else {
		viper.SetConfigName("netlify")
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.AddConfigPath(filepath.Join(home, ".netlify"))
	}

	return viper.ReadInConfig()
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
