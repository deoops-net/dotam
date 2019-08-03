package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/techmesh/dotam/conf"
	"github.com/techmesh/dotam/util"

	"github.com/mitchellh/cli"
	log "github.com/sirupsen/logrus"
)

func GitCMDFactor() (cli.Command, error) {
	return GitCmd{}, nil
}

type GitCmd struct {
}

func (r GitCmd) Help() string {
	return "this is help message"
}

func (r GitCmd) Run(args []string) (exitCode int) {
	if len(args) == 0 {
		return
	}

	for _, arg := range args {
		if arg == "--add-hook" {
			hookPath, err := filepath.Abs(".git/hooks/commit-msg")
			if err != nil {
				exitCode = -1
				log.Error(err)
				return
			}
			util.WriteFile(conf.DemoHook, hookPath)
			log.Info("hook file added!")
			return
		}
	}

	lockFile, err := filepath.Abs("Dotamfile.lock")
	if err != nil {
		exitCode = -1
		log.Error(err)
		return
	}
	digest := util.ReadFile(lockFile)

	path, err := filepath.Abs(args[0])
	if err != nil {
		exitCode = -1
		log.Error(err)
		return
	}

	old := util.ReadFile(path)
	fmt.Println(string(old))
	news := ""
	if strings.Contains(string(old), "signed by dotam") {
		news = util.ReplaceDigest(string(old), string(digest))
	}

	//util.WriteFile(path, news)

	f, err := os.OpenFile(path, os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	info, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	if err = f.Truncate(info.Size()); err != nil {
		log.Fatal(err)
	}

	if _, err := f.Write([]byte(news)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	return
}

func (r GitCmd) Synopsis() string {
	return string(`this cmd is used by git's commit-msg hook`)
}
