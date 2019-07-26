package main

import (
	"os"

	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
)

var CWD string
var RELEASE_VERSION []byte

func init() {
	RELEASE_VERSION = ReadFile("RELEASE")
	CWD, _ = os.Getwd()
	initLogLevel()
}

func initLogLevel() {
	level := os.Getenv("LOG_LEVEL")

	if level != "debug" {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}

	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	initCli()
}

func initCli() {
	c := cli.NewCLI("dotam", string(RELEASE_VERSION))
	c.Args = os.Args[1:]

	c.Commands = map[string]cli.CommandFactory{
		"build": BuildCMDFactor,
		"init":  InitCMDFactor,
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Error(err)
	}

	os.Exit(exitStatus)
}
