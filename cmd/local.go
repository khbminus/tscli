package cmd

import (
	"fmt"
	"github.com/khbminus/tscli/client"
	"github.com/khbminus/tscli/config"
	"github.com/khbminus/tscli/util"
	"github.com/logrusorgru/aurora/v3"
)

const (
	ConfigPath = "./.tscli.local"
)

func ShowLocalConfig() error {
	cfg, err := config.NewConfig(ConfigPath)
	if err != nil {
		return err
	}
	var compilerName string = "Doesn't set"
	if cfg.Compilers != nil && cfg.DefaultLang != -1 {
		compilerName = cfg.Compilers[cfg.DefaultLang].CompilerName
	}
	fmt.Printf("Contest ID: %v\nDefault Compiler: %v\n", cfg.Contest, compilerName)
	fmt.Println("\nAvailable problems:")

	for _, problem := range cfg.Problems {
		fmt.Printf("%v: %v\n", problem.ProblemId, problem.ProblemName)
	}
	fmt.Println("\nAvaliable Compiler:")

	for i, compiler := range cfg.Compilers {
		fmt.Printf("%v, %v", compiler.CompilerLang, compiler.CompilerName)
		if i == cfg.DefaultLang {
			fmt.Println("*")
		} else {
			fmt.Println()
		}
	}

	return nil
}

func ParseConfig() error {
	fmt.Println(aurora.Yellow("Getting new config..."))
	cfg, err := client.Instance.GetConfig(ConfigPath)
	if err != nil {
		return err
	}
	fmt.Println(aurora.Green("Got! Saving..."))
	if err := cfg.Save(); err != nil {
		return err
	}
	fmt.Println(aurora.Green("Saved. Use ts-cli local show to see current local config"))
	return nil
}

func ChooseContest(contestId string) error {
	fmt.Println(aurora.Yellow("Getting contests..."))
	contests, err := client.Instance.GetAvailableContests()
	if err != nil {
		fmt.Println(aurora.Red("Error at getting contests!"))
		return err
	}
	index := -1
	if contestId == "" {
		for i, v := range contests {
			fmt.Printf("%v) %v\n", i, v)
		}
		index = util.ChooseIndex(len(contests))
	} else {
		for i, contest := range contests {
			if contest.ContestId == contestId {
				index = i
				break
			}
		}
	}
	if index == -1 {
		fmt.Println(aurora.Red("Can't find any suitable contest!"))
		return nil
	}
	err = client.Instance.ChangeContest(contests[index].ContestId)

	if err != nil {
		fmt.Println(aurora.Red("Error at changing contests!"))
		return err
	}
	err = ParseConfig()
	if err != nil {
		return err
	}
	return nil
}

func ChangeDefaultLang() error {
	cfg, err := config.NewConfig(ConfigPath)
	if err != nil {
		return err
	}
	for i, v := range cfg.Compilers {
		fmt.Printf("%v) %v: %v\n", i, v.CompilerLang, v.CompilerName)
	}
	cfg.DefaultLang = util.ChooseIndex(len(cfg.Compilers))
	return cfg.Save()
}
