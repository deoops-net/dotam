package main

import (
	"reflect"
	"testing"

	"github.com/hashicorp/hcl"
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
	files := []string{"main.go", "Dotamfile.hcl", "Dotamfile.json"}
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
