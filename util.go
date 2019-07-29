package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/flosch/pongo2"
	docker "github.com/fsouza/go-dockerclient"
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

func VarToTplContext(vars map[string]Var, args CmdArgs) pongo2.Context {
	c := pongo2.Context{}
	c["_args"] = args
	for k, v := range vars {
		c[k] = v
	}

	return c
}

//func AppendToTplContext(c *pongo2.Context, conf map[string]interface{}) {
//	c["$args"] = map[string]interface{}{}
//	for k, v := range conf {
//		c["$args"] = m
//	}
//}

func TempVarToTplContext(vars map[string]interface{}) pongo2.Context {
	c := pongo2.Context{}
	for k, v := range vars {
		c[k] = v
	}

	return c
}

func RunTasks(conf DotamConf) (err error) {
	if err = ProcessTemp(conf.Temp); err != nil {
		return
	}

	if err = ProcessDocker(conf.Docker); err != nil {
		return
	}

	return
}

func ProcessTemp(temps map[string]Temp) error {
	if temps == nil {
		return nil
	}

	for k, v := range temps {
		var destFile string
		// srcFile := Abs(v.Src)
		if v.Dest == "." || v.Dest == "./" {
			destFile = Abs(k)
		}
		tempVars := TempVarToTplContext(v.Var)
		log.WithFields(log.Fields{"PROCESS": "TEMP VARS"}).Debug(tempVars)
		tpl := ReadFile(Abs(v.Src))
		destData, err := Render(string(tpl), tempVars)
		if err != nil {
			return err
		}

		log.WithFields(log.Fields{"PROCESS": "TEMP DEST"}).Debug(destData)

		if err = WriteFile(destData, destFile); err != nil {
			return err
		}
	}

	return nil
}

func ProcessDocker(d Docker) (err error) {
	if d == (Docker{}) {
		return
	}

	c, err := docker.NewClientFromEnv()
	if err != nil {
		log.Error("current env doesn't support docker, pls check your docker installation")
		panic(err)
	}

	if err = BuildImage(d, c); err != nil {
		return
	}

	if d.Auth == (Auth{}) {
		return
	}

	if err = PushImage(d, c); err != nil {
		return
	}

	return
}

func Exist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func ParseBuildArgs(src []string) (dotamFile string, dest []string) {

	str := strings.Join(src, " ")
	for _, d := range DEFAULT_DOTAMFILES {
		if strings.Contains(str, d) {
			dotamFile = d
		}
		str = strings.Replace(str, d, "", -1)
		str = strings.TrimSpace(str)
	}
	dest = strings.Split(str, " ")
	return
}

func BuildImage(d Docker, c *docker.Client) (err error) {

	log.WithFields(log.Fields{"PROCESS": "DOCKER REPO"}).Debug(d.Repo)
	log.WithFields(log.Fields{"PROCESS": "DOCKER TAG"}).Debug(d.Tag)
	log.WithFields(log.Fields{"PROCESS": "DOCKER AUTH"}).Debug(d.Auth)

	imageName := d.Repo + ":" + d.Tag
	log.Debug(imageName)
	// TODO parse Dockerfile from conf
	if err = c.BuildImage(docker.BuildImageOptions{
		Name:                imageName,
		ContextDir:          ".",
		Dockerfile:          "Dockerfile",
		SuppressOutput:      false,
		OutputStream:        os.Stdout,
		RmTmpContainer:      true,
		ForceRmTmpContainer: true,
	}); err != nil {
		return
	}
	return
}

// PushImage ...
// TODO
func PushImage(d Docker, c *docker.Client) (err error) {

	if err = c.PushImage(docker.PushImageOptions{
		Name:         d.Repo,
		Tag:          d.Tag,
		OutputStream: os.Stdout,
	}, docker.AuthConfiguration{
		Username: d.Auth.Username,
		Password: d.Auth.Password,
	}); err != nil {
		return
	}

	return
}

func ArgsToMiddleTemp(conf *DotamConf, args []string) {
	for _, v := range args {
		p := strings.Split(v, "=")
		conf.CmdArgs[p[0]] = p[1]
	}
}
