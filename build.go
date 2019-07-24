package main

import (
	"encoding/json"
	"errors"
	"path/filepath"

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

// load file
// pre-render
// run tasks
func (r RunCmd) Run(args []string) (extCode int) {
	var err error
	var dotamFile string
	var renderData pongo2.Context

	defer func() {
		if err != nil {
			log.Error(err)
			extCode = -1
			return
		}
	}()

	// read config
	defaultConf, exist := hasDefaultConf()
	if len(args) == 0 {
		if !exist {
			err = errors.New("you need at least a conf file, pls see help doc")
			return
		}
		dotamFile = Abs(defaultConf)
	} else {
		dotamFile = Abs(args[0])
	}

	data := ReadFile(dotamFile)
	config := DotamConf{}
	if err = parseConf(&config, string(data), dotamFile); err != nil {
		return
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
	if err = hcl.Decode(&newConfig, newDotamSrc); err != nil {
		return
	}

	if err = RunTasks(newConfig); err != nil {
		log.Error(err)
		return
	}

	log.Info("Congratulations! All works done!")
	return
}

func (r RunCmd) Synopsis() string {
	return "run pipline tasks, default config file is Dotamfile.{json,yml,hcl}"
}

func hasDefaultConf() (f string, e bool) {
	for _, v := range DEFAULT_DOTAMFILES {
		if Exist(v) {
			f = v
			e = true
			return
		}
	}

	return
}

func parseConf(conf *DotamConf, data, dotamFile string) error {
	// data := ReadFile(dotamFile)
	ext := filepath.Ext(dotamFile)
	log.Debug("parsing conf file ext is: ", ext)

	if ext == ".hcl" {
		return hcl.Decode(conf, string(data))
	}

	if ext == ".json" {
		return json.Unmarshal([]byte(data), conf)
	}

	return nil
}
