package main

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
	Repo string `json:"repo" hcl:"repo" yaml:"repo"`
	Tag  string `json:"tag" hcl:"tag" yaml:"tag"`
	Auth Auth   `json:"auth" hcl:"auth" yaml:"auth"`
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
