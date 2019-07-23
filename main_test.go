package main

import (
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
