package conf

import (
	"reflect"

	docker "github.com/fsouza/go-dockerclient"
)

type DotamConf struct {
	Plugin  map[string]Plugin `json:"plugin" hcl:"plugin" yaml:"plugin"`
	Temp    map[string]Temp   `json:"temp" hcl:"temp" yaml:"temp"`
	Var     map[string]Var    `json:"var" hcl:"var" yaml:"var"`
	Docker  Docker            `json:"docker" hcl:"docker" yaml:"docker"`
	CmdArgs CmdArgs           `json:"cmdArgs" hcl:"cmdArgs" yaml: "cmdArgs"`
}

type Plugin struct {
	Command  string                 `json:"command" hcl:"command" yaml:"command"`
	Args     []string               `json:"args" hcl:"args" yaml:"args"`
	Settings map[string]interface{} `json:"settings" hcl:"settings" yaml:"settings"`
}

type Temp struct {
	Dest string                 `json:"dest" hcl:"dest" yaml:"dest"`
	Src  string                 `json:"src" hcl:"src" yaml:"src"`
	Var  map[string]interface{} `json:"var" hcl:"var" yaml:"var"`
}

type Docker struct {
	Repo       string     `json:"repo" hcl:"repo" yaml:"repo"`
	Tag        string     `json:"tag" hcl:"tag" yaml:"tag"`
	Dockerfile string     `json:"dockerfile,omitempty" hcl:"dockerfile,omitempty" yaml:"dockerfile,omitempty"`
	BuildArgs  []BuildArg `json:"buildArgs,omitempty" hcl:"buildArgs,omitempty" yaml:"buildArgs,omitempty"`
	Auth       Auth       `json:"auth" hcl:"auth" yaml:"auth"`
	NotPrivate bool       `json:"notPrivate,omitempty" hcl:"notPrivate,omitempty" yaml:"notPrivate,omitempty"`
	Caporal    Caporal    `json:"caporal" hcl:"caporal" yaml:"caporal"`
}

// BuildArg the docker --build-args
type BuildArg struct {
	Key   string `json:"key,omitempty" hcl:"key,omitempty" yaml:"key,omitempty"`
	Value string `json:"value,omitempty" hcl:"value,omitempty" yaml:"value,omitempty"`
}

// Caporal is a built-in plugin for scheduling containers remotely
type Caporal struct {
	Host    string         `json:"host" hcl:"host" yaml:"host"`
	Name    string         `json:"name" hcl:"name" yaml:"name"`
	Options CaporalOptions `json:"opts" hcl:"opts" yaml: "opts"`
}

// CaporalOptions are flags for docker run
type CaporalOptions struct {
	Publish []string    `json:"publish" hcl:"publish" yaml:"publish"`
	Network string      `json:"network" hcl:"network" yaml:"network"`
	Mount   []HostMount `json:"mount" hcl:"mount" yaml:"mount"`
}
type HostMount struct {
	Bind string `json:"bind" hcl:"bind" yaml:"bind"`
	Type string `json:"type" hcl:"type" yaml:"type"`
}

type Auth struct {
	Username string `json:"username" hcl:"username" yaml:"username"`
	Password string `json:"password" hcl:"password" yaml:"password"`
}

// CmdArgs ...
// the params the after `build` command
// e.g.
// $dotam build reg_user=tom reg_pass=pass [Dotamfile.hcl]
type CmdArgs map[string]interface{}

type Var map[string]interface{}

func (d Docker) CreateBuildArgs() []docker.BuildArg {
	ba := []docker.BuildArg{}
	if d.BuildArgs == nil {
		return nil
	}

	for _, v := range d.BuildArgs {
		b := docker.BuildArg{
			Name:  v.Key,
			Value: v.Value,
		}
		ba = append(ba, b)
	}

	return ba
}

func (d Docker) CreateAuthConfig() docker.AuthConfiguration {
	auth := docker.AuthConfiguration{
		//Username: "",
		//Password: "",
	}
	if d.NotPrivate == true || reflect.DeepEqual(d.Auth, (Docker{}).Auth) {
		auth.RegistryToken = " "
		return auth
	}

	return auth
}
