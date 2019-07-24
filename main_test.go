package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/hcl"
	"gopkg.in/yaml.v2"
)

type Account struct {
	Name string
}

func TestParseHcl(t *testing.T) {
	// bytesTest := []byte(`name = "foo"`)
	a := Account{}
	o := Account{Name: "tom"}
	// b := Account{}
	if hcl.Decode(&a, `name = "tom"`); a != o {
		t.Fail()
	}
}

func TestExist(t *testing.T) {
	files := []string{"main.go", "Dotamfile.hcl", "Dotamfile.jon"}
	expects := []bool{true, true, false}
	results := []bool{}

	for _, f := range files {
		t.Log(Exist(f))
		results = append(results, Exist(f))
	}
	if !reflect.DeepEqual(expects, results) {
		t.Fail()
	}
}

func TestParseYaml(t *testing.T) {
	// data := map[string]interface{}{}
	c := DotamConf{}
	file := ReadFile("Dotamfile.yml")
	err := yaml.Unmarshal(file, &c)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
	fmt.Println(c)

}
