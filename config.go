package main

type DotamConf struct {
	Plugin map[string]Plugin `json:"plugin" hcl:"plugin" yaml:"plugin"`
	Temp   map[string]Temp   `json:"temp" hcl:"temp" yaml:"temp"`
	Var    map[string]Var    `json:"var" hcl:"var" yaml:"var"`
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

type Var map[string]interface{}
