package cli

import (
	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
)

func GitCMDFactor() (cli.Command, error) {
	return InitCmd{}, nil
}

type GitCmd struct {
}

func (r GitCmd) Help() string {
	return "this is help message"
}

func (r GitCmd) Run(args []string) (exitCode int) {
	log.Debug(args)

	return
}

func (r GitCmd) Synopsis() string {
	return string(`this cmd is used by git's commit-msg hook`)
}
