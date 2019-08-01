package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/flosch/pongo2"
	docker "github.com/fsouza/go-dockerclient"
	log "github.com/sirupsen/logrus"
)

func ReadFile(file string) []byte {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	return data
}

func WriteFile(data, file string) error {
	err := ioutil.WriteFile(file, []byte(data), 0644)
	if err != nil {
		return err
	}
	return nil
}

func Abs(path string) string {
	return fmt.Sprintf("%s/%s", CWD, path)
}

func Render(src string, data pongo2.Context) (out string, err error) {

	tpl, err := pongo2.FromString(src)
	if err != nil {
		return
	}

	if out, err = tpl.Execute(data); err != nil {
		return
	}

	return
}

func VarToTplContext(vars map[string]Var, args CmdArgs) pongo2.Context {
	c := pongo2.Context{}
	c["_args"] = args
	for k, v := range vars {
		c[k] = v
	}

	return c
}

//func AppendToTplContext(c *pongo2.Context, conf map[string]interface{}) {
//	c["$args"] = map[string]interface{}{}
//	for k, v := range conf {
//		c["$args"] = m
//	}
//}

func TempVarToTplContext(vars map[string]interface{}) pongo2.Context {
	c := pongo2.Context{}
	for k, v := range vars {
		c[k] = v
	}

	return c
}

func RunTasks(conf DotamConf) (err error) {
	if err = ProcessTemp(conf.Temp); err != nil {
		return
	}

	if err = ProcessDocker(conf.Docker); err != nil {
		log.Error(err)
		return
	}

	return
}

func ProcessTemp(temps map[string]Temp) error {
	if temps == nil {
		return nil
	}

	for k, v := range temps {
		var destFile string
		// srcFile := Abs(v.Src)
		if v.Dest == "." || v.Dest == "./" {
			destFile = Abs(k)
		}
		tempVars := TempVarToTplContext(v.Var)
		log.WithFields(log.Fields{"PROCESS": "TEMP VARS"}).Debug(tempVars)
		tpl := ReadFile(Abs(v.Src))
		destData, err := Render(string(tpl), tempVars)
		if err != nil {
			return err
		}

		log.WithFields(log.Fields{"PROCESS": "TEMP DEST"}).Debug(destData)

		if err = WriteFile(destData, destFile); err != nil {
			return err
		}
	}

	return nil
}

func ProcessDocker(d Docker) (err error) {
	if reflect.DeepEqual(d, Docker{}) {
		return
	}

	c, err := docker.NewClientFromEnv()
	if err != nil {
		log.Error("current env doesn't support docker, pls check your docker installation")
		panic(err)
	}

	if err = BuildImage(d, c); err != nil {
		return
	}

	if d.Auth == (Auth{}) {
		return
	}

	if err = PushImage(d, c); err != nil {
		return
	}

	if err = ScheduleContainer(d, c); err != nil {
		return
	}

	return
}

func Exist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func ParseBuildArgs(src []string) (dotamFile string, dest []string) {

	str := strings.Join(src, " ")
	for _, d := range DEFAULT_DOTAMFILES {
		if strings.Contains(str, d) {
			dotamFile = d
		}
		str = strings.Replace(str, d, "", -1)
		str = strings.TrimSpace(str)
	}
	dest = strings.Split(str, " ")
	return
}

func BuildImage(d Docker, c *docker.Client) (err error) {

	dockerfile := ""
	if len(d.Dockerfile) != 0 {
		dockerfile = d.Dockerfile
	} else {
		dockerfile = "Dockerfile"
	}
	log.WithFields(log.Fields{"PROCESS": "DOCKER REPO"}).Debug(d.Repo)
	log.WithFields(log.Fields{"PROCESS": "DOCKER TAG"}).Debug(d.Tag)
	log.WithFields(log.Fields{"PROCESS": "DOCKER AUTH"}).Debug(d.Auth)

	imageName := d.Repo + ":" + d.Tag
	log.Debug(imageName)
	log.WithFields(log.Fields{"DOCKER": "BUILD ARGS FLAG"}).Debug(d.CreateBuildArgs())
	// TODO parse Dockerfile from conf
	if err = c.BuildImage(docker.BuildImageOptions{
		Name:                imageName,
		ContextDir:          ".",
		Dockerfile:          dockerfile,
		SuppressOutput:      false,
		OutputStream:        os.Stdout,
		RmTmpContainer:      true,
		ForceRmTmpContainer: true,
		BuildArgs:           d.CreateBuildArgs(),
	}); err != nil {
		return
	}
	return
}

// PushImage ...
// TODO
func PushImage(d Docker, c *docker.Client) (err error) {

	if err = c.PushImage(docker.PushImageOptions{
		Name:         d.Repo,
		Tag:          d.Tag,
		OutputStream: os.Stdout,
	}, docker.AuthConfiguration{
		Username: d.Auth.Username,
		Password: d.Auth.Password,
	}); err != nil {
		return
	}

	return
}

// Schedule it with caporal
func ScheduleContainer(d Docker, c *docker.Client) (err error) {
	// just skip
	if reflect.DeepEqual(d.Caporal, Caporal{}) {
		return nil
	}

	j, err := json.Marshal(d)
	log.WithFields(log.Fields{"DOCKER": "SCHEDULE"}).Debug(string(j))

	// TODO move to a pre defined struct
	//{"repo": "nginx", "tag": "latest", "name": "mynginx-2", "opts": {"publish": ["10009:80"]}}
	postBody := struct {
		Repo string         `json:"repo"`
		Tag  string         `json:"tag"`
		Name string         `json:"name"`
		Opts CaporalOptions `json:"opts"`
	}{
		Repo: d.Repo,
		Tag:  d.Tag,
		Name: d.Caporal.Name,
		Opts: d.Caporal.Options,
	}
	payload, err := json.Marshal(postBody)
	if err != nil {
		return
	}

	// do request for caporal
	log.Debug(string(payload))
	req, err := http.NewRequest("PUT", d.Caporal.Host+"/container", bytes.NewBuffer(payload))
	authCode := base64.StdEncoding.EncodeToString(Encrypt([]byte(d.Auth.Username+":"+d.Auth.Password), "caoral-salt"))
	log.WithFields(log.Fields{"CAPORAL": "AUTH CODE"}).Debug(authCode)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-AUTH", authCode)
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	defer res.Body.Close()
	log.Debug(string(body))

	log.Debug(res.StatusCode)

	return
}

func ArgsToMiddleTemp(conf *DotamConf, args []string) {
	for _, v := range args {
		p := strings.Split(v, "=")
		conf.CmdArgs[p[0]] = p[1]
	}
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func Encrypt(data []byte, passphrase string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext
}

func Decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func encryptFile(filename string, data []byte, passphrase string) {
	f, _ := os.Create(filename)
	defer f.Close()
	f.Write(Encrypt(data, passphrase))
}

func decryptFile(filename string, passphrase string) []byte {
	data, _ := ioutil.ReadFile(filename)
	return Decrypt(data, passphrase)
}
