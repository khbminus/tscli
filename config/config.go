package config

import (
	"encoding/json"
	"fmt"
	. "github.com/logrusorgru/aurora/v3"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	DefaultLang int        `json:"default_lang"`
	Contest     string     `json:"contest"`
	Problems    []Problem  `json:"problems"`
	Compilers   []Compiler `json:"compilers"`
	Path        string
}

func NewConfig(path string) (c *Config, err error) {
	c = &Config{
		DefaultLang: -1,
		Contest:     "",
		Path:        path,
	}
	if err := c.load(); err != nil {
		fmt.Println(Red(err.Error()))
		fmt.Println(Green(fmt.Sprintf("Creating a new local config at %v", path)))
	}

	if err := c.Save(); err != nil {
		fmt.Println(Red("Error while saving!"))
		fmt.Println(Red(err.Error()))
		return nil, err
	}
	return
}

func (c *Config) load() (err error) {
	file, err := os.Open(c.Path)
	if err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, c)
}

func (c *Config) Save() (err error) {
	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(c.Path), os.ModePerm); err != nil {
		return err
	}
	err = os.WriteFile(c.Path, data, 0644)
	return err
}
