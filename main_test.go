package main

import (
	"fmt"
	"net/url"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/techmesh/dotam/conf"
	"github.com/techmesh/dotam/util"

	docker "github.com/fsouza/go-dockerclient"

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
		results = append(results, util.Exist(f))
	}
	if !reflect.DeepEqual(expects, results) {
		t.Fail()
	}
}

func TestParseYaml(t *testing.T) {
	// data := map[string]interface{}{}
	c := conf.DotamConf{}
	file := util.ReadFile("Dotamfile.yml")
	err := yaml.Unmarshal(file, &c)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}

}

func TestParseArgs(t *testing.T) {
	data := []string{"Dotamfile.hcl", "reg_user=tom", "reg_pass=foo"}
	data2 := []string{"reg_user=tom", "foo=bar"}

	f, a := util.ParseBuildArgs(data)
	if !reflect.DeepEqual([]string{"reg_user=tom", "reg_pass=foo"}, a) {
		t.Log(a)
		t.Fail()
	}

	if f != "Dotamfile.hcl" {
		t.Log(f)
		t.Fail()
	}

	f2, a2 := util.ParseBuildArgs(data2)
	if !reflect.DeepEqual([]string{"reg_user=tom", "foo=bar"}, a2) {
		t.Log(a2)
		t.Fail()
	}

	if f2 != "" {
		t.Log(f2)
		t.Fail()
	}
}

func TestBuildDockerImage(t *testing.T) {

	c, err := docker.NewClientFromEnv()
	if err != nil {
		// log.Error("current env doesn't support docker, pls check your docker installation")
		panic(err)
	}
	// var buf bytes.Buffer
	// dockerfile, err := filepath.Abs("Dockerfile")
	if err = c.BuildImage(docker.BuildImageOptions{
		Name:                "dotam",
		ContextDir:          ".",
		Dockerfile:          "Dockerfile",
		RmTmpContainer:      true,
		ForceRmTmpContainer: true,
		// Remote:       "https://reg.yunpro.cn",
		SuppressOutput: false,
		// InputStream:    &buf,
		OutputStream: os.Stdout,
	}); err != nil {
		panic(err)
	}
	// log.Debug(DockerClient)
}

func TestParseBuildArgs(t *testing.T) {
	args := []string{"reg=foo", "pass=bar"}
	qs := strings.Join(args, "&")
	data, err := url.ParseQuery(qs)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	fmt.Println(data)

}

func TestPushImage(t *testing.T) {
	c, err := docker.NewClientFromEnv()
	if err != nil {
		// log.Error("current env doesn't support docker, pls check your docker installation")
		panic(err)
	}

	d := conf.Docker{
		Repo: "dhub.yiliang.cn/adpro/g1-gdt-ios-stupidball",
		Tag:  "latest",
		//Dockerfile: "",
		//BuildArgs:  nil,
		Auth: conf.Auth{
			//Username: "abc",
			//Password: "bcd",

		},
		NotPrivate: true,
		//Caporal:    conf.Caporal{},
	}
	if err := util.PushImage(d, c); err != nil {
		t.Log(err)
		t.Fail()
	}
}

func TestReplaceDigest(t *testing.T) {
	a := ` support logDriver flag
    support git commit hook
    support git --add-hook cli
    
    signed by dotam:
    checksum: 3d8394dd4e1db5151011a707c884d6f9
    build status: success
	
	some other tools generated
`
	b := `something new`

	re := regexp.MustCompile(`(?m)signed by dotam:.*\n.*\n.*`)
	newStr := re.ReplaceAllLiteralString(a, b)
	t.Log(newStr)
}
