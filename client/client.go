package client

import (
	"encoding/json"
	"fmt"
	"github.com/khbminus/tscli/cookiejar"
	. "github.com/logrusorgru/aurora/v3"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

var Instance *Client

const (
	HOST = "https://acm.math.spbu.ru/tsweb"
)

func Init(path string) {
	jar, _ := cookiejar.New(nil)
	c := &Client{Jar: jar,
		User:     "",
		Password: "",
		client:   nil,
		path:     path,
	}
	if err := c.load(); err != nil {
		fmt.Println(Red(err.Error()))
		fmt.Println(Green(fmt.Sprintf("Creating a new client config at %v", path)))
		c.client = &http.Client{Jar: jar}
	} else {
		c.client = &http.Client{Jar: c.Jar}
	}

	if err := c.save(); err != nil {
		fmt.Println(Red("Error while saving!"))
		fmt.Println(Red(err.Error()))
	}
	Instance = c
	if res, err := c.FindLogin(); err != nil || res == "" {
		fmt.Println(Magenta("Not logged..."))
		if err := c.Login(); err != nil {
			return
		}
		fmt.Println("Logged")
	}
}

func (c *Client) load() (err error) {
	file, err := os.Open(c.path)
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, c)
}

func (c *Client) save() (err error) {
	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(c.path), os.ModePerm); err != nil {
		return err
	}
	err = os.WriteFile(c.path, data, 0644)
	return err
}
