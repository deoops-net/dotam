package main

import (
	"os"

	cliTool "github.com/mitchellh/cli"
	"github.com/techmesh/dotam/cli"

	log "github.com/sirupsen/logrus"
)

func init() {
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
	c := cliTool.NewCLI("dotam", string(RELEASE_VERSION))
	c.Args = os.Args[1:]

	c.Commands = map[string]cliTool.CommandFactory{
		"build": cli.BuildCMDFactor,
		"init":  cli.InitCMDFactor,
		"git":   cli.GitCMDFactor,
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Error(err)
	}

	os.Exit(exitStatus)
}
