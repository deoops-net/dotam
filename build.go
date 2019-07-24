package main

import (
	"github.com/flosch/pongo2"
	"github.com/hashicorp/hcl"
	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
)

func BuildCMDFactor() (cli.Command, error) {
	return RunCmd{}, nil
}

type RunCmd struct {
}

func (r RunCmd) Help() string {
	return "this is help message"
}

func (r RunCmd) Run(args []string) int {
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

func (r RunCmd) Synopsis() string {
	return "run pipline tasks, default config file is Dotamfile.{json,yml,hcl}"
}
