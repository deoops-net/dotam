package cli

import (
	"strings"

	"github.com/techmesh/dotam/conf"
	"github.com/techmesh/dotam/util"

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

func (r InitCmd) Run(args []string) (exitCode int) {
	var destFile string
	var destData string
	var err error

	defer func() {
		if err != nil {
			log.Error(err)
			exitCode = -1
			return
		}
		log.Infof("Congratulations! %s generated.", destFile)
	}()
	log.WithFields(log.Fields{"CMD INIT": "RUN"}).Debug(args)

	// gen hcl as default
	if len(args) == 0 {
		destFile = conf.DEMO_HCL
		if err = genDemoFile(conf.DemoHcl, destFile); err != nil {
			return
		}
	} else {
		// TODO maybe we need a better parser
		switch strings.Join(args, "") {
		case "-tyaml", "-tyml":
			destFile = conf.DEMO_YAML
			destData = conf.DemoYaml
		case "-thcl":
			destFile = conf.DEMO_HCL
			destData = conf.DemoHcl
		case "-tjson":
			destFile = conf.DEMO_JSON
			destData = conf.DemoJson
		}
		log.Debug(destFile)
		genDemoFile(destData, destFile)
	}

	return
}

func (r InitCmd) Synopsis() string {
	return string(`initial a demo conf file in current dir, as default it will generate a .hcl file, for specific format use: -t [yaml|yml, json, hcl]	
	`)
}

func genDemoFile(data, file string) error {
	return util.WriteFile(data, file)
}
