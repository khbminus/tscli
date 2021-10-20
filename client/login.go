package client

import (
	"errors"
	"fmt"
	"github.com/khbminus/tscli/cookiejar"
	"github.com/khbminus/tscli/util"
	. "github.com/logrusorgru/aurora/v3"
	"golang.org/x/crypto/ssh/terminal"
	"net/url"
	"syscall"
)

func (c *Client) Login() (err error) {
	fmt.Println(Cyan(fmt.Sprintf("Login %v...\n", c.User)))

	// Clean up jar
	jar, _ := cookiejar.New(nil)
	c.client.Jar = jar
	body, err := util.PostForm(c.client, HOST+"/index", url.Values{
		"team":      {c.User},
		"password":  {c.Password},
		"op":        {"login"},
		"contestId": {""}, // TODO: local config
	})
	if tsError := findErrorRegex.FindSubmatch(body); tsError != nil {
		fmt.Println(Sprintf(Red("Submit error by TSweb: %v\n"), string(tsError[1])))
		return errors.New("TSWeb internal error")
	}
	if err != nil {
		return
	}
	fmt.Println(Green("Successfully logged in"))
	c.Jar = jar
	return c.save()
}

func (c *Client) EnterCredentials() (err error) {
	if c.User != "" {
		fmt.Println(Sprintf(BrightYellow("Current user: %v\n"), c.User))
	}
	fmt.Print("Login: ")
	c.User = util.ScanlineTrim()
	fmt.Print("Enter password: ")
	bytePassword, err := terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		return err
	}
	c.Password = string(bytePassword)
	return c.Login()
}
