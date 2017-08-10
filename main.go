package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/netlify/netlifyctl/commands"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:    true,
		DisableTimestamp: false,
		TimestampFormat:  time.RFC822Z,
	})
	commands.Execute()
}
