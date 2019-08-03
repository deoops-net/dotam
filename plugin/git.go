package plugin

import (
	"fmt"
	"path/filepath"

	//"github.com/labstack/gommon/log"

	"gopkg.in/src-d/go-git.v4"
)

type Git struct {
}

func (g Git) AddTracked() {
	root, err := filepath.Abs("../")
	if err != nil {
		panic(err)
	}
	fmt.Println(root)

	repo, err := git.PlainOpen(root)
	if err != nil {
		panic(err)
	}

	fmt.Println(repo)
	//log.Debug(repo)
}

func (g Git) AddAll() {

}

func (g Git) Commit() {

}

func (g Git) Push() {

}

func (g Git) CallHooks() {

}
