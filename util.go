package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/flosch/pongo2"
	log "github.com/sirupsen/logrus"
)

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

func Exist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}

}
