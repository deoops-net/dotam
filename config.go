package main

type DotamConf struct {
	Plugin map[string]Plugin `hcl:plugin`
	Temp   map[string]Temp   `hcl:temp`
	Var    map[string]Var    `hcl:var`
}

type Plugin struct {
	Command  string                 `hcl:command`
	Args     []string               `hcl:args`
	Settings map[string]interface{} `hcl:settings`
}

type Temp struct {
	Dest string                 `hcl:dest`
	Src  string                 `hcl:src`
	Var  map[string]interface{} `hcl:var`
}

type Var map[string]interface{}
