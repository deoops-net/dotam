package cli

import (
	"encoding/json"
	"errors"
	"path/filepath"

	"github.com/techmesh/dotam/util"

	"github.com/techmesh/dotam/conf"

	"github.com/flosch/pongo2"
	"github.com/hashicorp/hcl"
	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
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
	var buildArgs []string
	config := conf.DotamConf{}
	log.WithFields(log.Fields{"BUILD": "ARGS"}).Debug(args)

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
		dotamFile = util.Abs(defaultConf)
	} else {
		config.CmdArgs = map[string]interface{}{}
		dotamFile, buildArgs = util.ParseBuildArgs(args)

		if dotamFile == "" {
			if !exist {
				err = errors.New("you need at least a conf file, pls see help doc")
				return
			} else {
				dotamFile = defaultConf
			}
		}
	}

	// read src dotamfile
	data := util.ReadFile(dotamFile)
	if err = parseConf(&config, string(data), dotamFile); err != nil {
		return
	}
	// convert args into middle variables append to middle conf
	// 1, convert cli args to map[string]interface
	// 2, append args to pongo config
	// 3, replace $variable to pongo mark {{}}
	// 4, render them to middle
	util.ArgsToMiddleTemp(&config, buildArgs)

	renderData = util.VarToTplContext(config.Var, config.CmdArgs)

	// after render middle remove this middle variables
	// render middle conf
	newDotamSrc, err := util.Render(string(data), renderData)
	if err != nil {
		return -1
	}
	log.WithFields(log.Fields{"BUILD": "RENDERED DOC"}).Debug(newDotamSrc)

	newConfig := conf.DotamConf{}
	if err = parseConf(&newConfig, newDotamSrc, dotamFile); err != nil {
		return
	}
	log.Debug(newConfig)

	if err = util.RunTasks(newConfig); err != nil {
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
	for _, v := range conf.DEFAULT_DOTAMFILES {
		if util.Exist(v) {
			f = v
			e = true
			return
		}
	}

	return
}

func parseConf(conf *conf.DotamConf, data, dotamFile string) error {
	// data := ReadFile(dotamFile)
	ext := filepath.Ext(dotamFile)
	log.Debug("parsing conf file ext is: ", ext)

	if ext == ".hcl" {
		return hcl.Decode(conf, string(data))
	}

	if ext == ".json" {
		return json.Unmarshal([]byte(data), conf)
	}

	if ext == ".yml" || ext == ".yaml" {
		return yaml.Unmarshal([]byte(data), conf)
	}

	return errors.New("not support file type")
}
