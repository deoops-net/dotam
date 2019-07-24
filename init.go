package main

import (
	"github.com/flosch/pongo2"
	"github.com/hashicorp/hcl"
	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
)

func InitCMDFactor() (cli.Command, error) {
	return InitCmd{}, nil
}

type InitCmd struct {
}

func (r InitCmd) Help() string {
	return "this is help message"
}

func (r InitCmd) Run(args []string) int {
	var dotamFile string
	var renderData pongo2.Context

	if len(args) == 0 {
		dotamFile = Abs("Dotamfile.hcl")
	} else {
		dotamFile = Abs(args[0])
	}

	data := ReadFile(dotamFile)
	config := DotamConf{}
	err := hcl.Decode(&config, string(data))
	if err != nil {
		log.Error(err)
	}

	if config.Var != nil {
		renderData = VarToTplContext(config.Var)
	}

	newDotamSrc, err := Render(string(data), renderData)
	if err != nil {
		log.Error(err)
		return -1
	}
	log.Debug(newDotamSrc)

	newConfig := DotamConf{}
	err = hcl.Decode(&newConfig, newDotamSrc)
	if err != nil {
		panic(err)
	}

	// log.Debug(newConfig.Temp)
	if err = RunTasks(newConfig); err != nil {
		log.Error(err)
		return -1
	}

	log.Info("Congratulations! All works done!")

	return 0
}

func (r InitCmd) Synopsis() string {
	return "initial a demo conf file in current dir"
}
