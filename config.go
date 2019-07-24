package main

type DotamConf struct {
	Plugin map[string]Plugin `json:"plugin" hcl:"plugin"`
	Temp   map[string]Temp   `json:"temp" hcl:"temp"`
	Var    map[string]Var    `json:"var" hcl:"var"`
}

type Plugin struct {
	Command  string                 `json:"command" hcl:"command"`
	Args     []string               `json:"args" hcl:"args"`
	Settings map[string]interface{} `json:"settings" hcl:"settings"`
}

type Temp struct {
	Dest string                 `json:"dest" hcl:"dest"`
	Src  string                 `json:"src" hcl:"src"`
	Var  map[string]interface{} `json:"var" hcl:"var"`
}

type Var map[string]interface{}
