package main

import (
	"fmt"
	"io/ioutil"

	"github.com/flosch/pongo2"

	// "log"
	"os"

	"github.com/hashicorp/hcl"
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
		"build": cmdFactor,
		"init":  cmdFactor,
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}

func cmdFactor() (cli.Command, error) {
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
	return "a test cmd"
}

func ReadFile(file string) []byte {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return data
}

func WriteFile(data, file string) error {
	err := ioutil.WriteFile(file, []byte(data), 0644)
	if err != nil {
		return err
	}
	return nil
}

func Abs(path string) string {
	return fmt.Sprintf("%s/%s", CWD, path)
}

func Render(src string, data pongo2.Context) (out string, err error) {

	tpl, err := pongo2.FromString(src)
	if err != nil {
		return
	}

	if out, err = tpl.Execute(data); err != nil {
		return
	}

	return
}

func VarToTplContext(vars map[string]Var) pongo2.Context {
	c := pongo2.Context{}
	for k, v := range vars {
		c[k] = v
	}

	return c
}

func TempVarToTplContext(vars map[string]interface{}) pongo2.Context {
	c := pongo2.Context{}
	for k, v := range vars {
		c[k] = v
	}

	return c
}

func RunTasks(conf DotamConf) error {
	return ProcessTemp(conf.Temp)
}

func ProcessTemp(temps map[string]Temp) error {
	if temps == nil {
		return nil
	}

	for k, v := range temps {
		var destFile string
		srcFile := Abs(v.Src)
		if v.Dest == "." || v.Dest == "./" {
			destFile = Abs(k)
		}
		log.Debug("source file is :", srcFile)
		log.Debug("destFile file is :", destFile)
		// log.Debug(k, v.Src)
		log.Debugf("%s var is: %v", k, v.Var)
		tempVars := TempVarToTplContext(v.Var)
		log.Debug(tempVars)
		tpl := ReadFile(Abs(v.Src))
		destData, err := Render(string(tpl), tempVars)
		if err != nil {
			return err
		}

		log.Debug("new rendered data:")
		log.Debug(destData)

		if err = WriteFile(destData, destFile); err != nil {
			return err
		}

	}

	return nil
}
