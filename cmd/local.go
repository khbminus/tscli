package cmd

import (
	"errors"
	"fmt"
	"github.com/khbminus/tscli/client"
	"github.com/khbminus/tscli/config"
	"github.com/khbminus/tscli/util"
	"github.com/logrusorgru/aurora/v3"
	"github.com/olekukonko/tablewriter"
	"os"
	"path"
	"strconv"
)

const (
	ConfigName = ".tscli.local"
)

func GetConfig() (*config.Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	for {
		if cwd == "/" || cwd == "" {
			fmt.Println(aurora.Red("Can't find a local config. Please run tscli local parse"))
			return nil, errors.New("no config found")
		}
		if _, err := os.Stat(cwd + "/" + ConfigName); err == nil {
			return config.NewConfig(cwd + "/" + ConfigName)
		} else if os.IsNotExist(err) {
			cwd = path.Dir(cwd)
		} else {
			panic(err)
		}
	}
}

func ShowLocalConfig() error {
	cfg, err := GetConfig()
	if err != nil {
		return err
	}
	compilerName := "Doesn't set"
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
	cfg, err := client.Instance.GetConfig("./" + ConfigName)
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
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"index", "Id", "Name", "Started", "Status"})
		table.SetAutoWrapText(false)
		table.SetAlignment(tablewriter.ALIGN_CENTER)
		for i, v := range contests {
			table.Append([]string{strconv.Itoa(i), v.ContestId, v.ContestName, v.ContestStarted, v.ContestStatus})
		}
		table.Render()
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
	cfg, _ := GetConfig()
	if cfg == nil {
		cfg, err = config.NewConfig("./" + ConfigName)
		if err != nil {
			return err
		}
	}
	cfg.Contest = contests[index].ContestId
	err = client.Instance.DownloadStatements(contests[index], cfg)
	if err != nil {
		return err
	}
	_ = cfg.Save()

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
	cfg, err := GetConfig()
	if err != nil {
		return err
	}
	for i, v := range cfg.Compilers {
		fmt.Printf("%v) %v: %v\n", i, v.CompilerLang, v.CompilerName)
	}
	cfg.DefaultLang = util.ChooseIndex(len(cfg.Compilers))
	return cfg.Save()
}
