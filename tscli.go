package main

import (
	"fmt"
	"github.com/devfacet/gocmd/v3"
	"github.com/khbminus/tscli/client"
	"github.com/khbminus/tscli/cmd"
	"github.com/logrusorgru/aurora/v3"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"path"
)

var (
	ClientPath, _ = homedir.Expand("~/.tscli.global")
)

func main() {

	flags := struct {
		Help      bool `short:"h" long:"help" description:"Display Usage" global:"true"`
		Version   bool `short:"v" long:"version" description:"Display version"`
		VersionEx bool `long:"vv" description:"Display version (extended)"`
		Local     struct {
			Show       struct{} `command:"show" description:"Show local config"`
			Parse      struct{} `command:"parse" description:"Parse config for current config"`
			SetContest struct {
				ContestId string `long:"id" description:"contest's id'"`
			} `command:"set-contest" description:"Optional, set contest for local config"`
			SetCompiler struct{} `command:"set-compiler" description:"Set compiler for local config"`
		} `command:"local" description:"Functions for local config" nonempty:"true"`
		Submit struct {
			Settings bool   `settings:"true" allow-unknown-arg:"true"`
			Problem  string `short:"p" long:"problem" description:"Optional, problem's id' to send (if filename is not the same as problem id)" required:"false"`
			//File string `name:"file" required:"false"`
		} `command:"submit" description:"Submit solution"`
		Login struct{} `command:"login" description:"Login to TSWeb and remember credentials"`
	}{}
	_, _ = gocmd.HandleFlag("Local.Show", func(gcmd *gocmd.Cmd, args []string) error {
		client.Init(ClientPath)
		return cmd.ShowLocalConfig()
	})
	_, _ = gocmd.HandleFlag("Local.SetContest", func(gcmd *gocmd.Cmd, args []string) error {
		client.Init(ClientPath)
		return cmd.ChooseContest(flags.Local.SetContest.ContestId)
	})

	_, _ = gocmd.HandleFlag("Local.Parse", func(gcmd *gocmd.Cmd, args []string) error {
		client.Init(ClientPath)
		return cmd.ParseConfig()
	})
	_, _ = gocmd.HandleFlag("Local.SetCompiler", func(cgmd *gocmd.Cmd, args []string) error {
		client.Init(ClientPath)
		return cmd.ChangeDefaultLang()
	})
	_, _ = gocmd.HandleFlag("Login", func(gocmd *gocmd.Cmd, args []string) error {
		client.Init(ClientPath)
		return client.Instance.EnterCredentials()
	})
	_, _ = gocmd.HandleFlag("Submit", func(gcmd *gocmd.Cmd, args []string) error {
		client.Init(ClientPath)
		filename := ""
		for _, v := range args[1:] {
			if len(v) > 0 && v[0] != '-' {
				if filename != "" {
					fmt.Println(aurora.Red("There are more than 1 file"))
					return nil
				}
				filename = v
			}
		}
		if filename == "" {
			fmt.Println(aurora.Red("Unknown filename"))
			return nil
		}
		cfg, err := cmd.GetConfig()
		if err != nil {
			fmt.Println(aurora.Red("Error at config load"))
			return err
		}
		problemId := path.Base(filename)
		problemId = problemId[:len(problemId)-len(path.Ext(problemId))]
		if flags.Submit.Problem != "" {
			problemId = flags.Submit.Problem
		}
		if cfg.Compilers == nil || cfg.DefaultLang == -1 {
			fmt.Println(aurora.Red("Please set default compiler and parse local config!"))
			return nil
		}
		body, err := ioutil.ReadFile(filename)
		if err != nil {
			fmt.Println(aurora.Red("Unable to read file!"))
			return err
		}
		for _, problem := range cfg.Problems {
			if problem.ProblemId == problemId {
				return cmd.SubmitAndWatch(problem, cfg.Compilers[cfg.DefaultLang], string(body))
			}
		}
		fmt.Println(aurora.Red("Unable to find problem with such id"))
		return nil
	})

	_, _ = gocmd.New(gocmd.Options{
		Name:        "tscli",
		Description: "TSWeb CLI client",
		Version:     fmt.Sprintf("%v (%v)", "v1.0.2", "Okay, let's go"),
		Flags:       &flags,
		ConfigType:  gocmd.ConfigTypeAuto,
	})
}
