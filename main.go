package main

import (
	"os"

	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
)

var CWD string

func init() {
	CWD, _ = os.Getwd()
	initLogLevel()
}

func initLogLevel() {
	level := os.Getenv("LOG_LEVEL")

	if level != "debug" {
		// fmt.Println("log at error level")
		log.SetLevel(log.InfoLevel)
	} else {
		// fmt.Println("log at debug level")
		log.SetLevel(log.DebugLevel)
	}

	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	c := cli.NewCLI("dopam", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"build": BuildCMDFactor,
		"init":  InitCMDFactor,
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
