package client

import (
	"fmt"
	"github.com/khbminus/tscli/config"
	"github.com/khbminus/tscli/util"
	"github.com/logrusorgru/aurora/v3"
	"os"
	"path/filepath"
	"regexp"
)

var (
	findLoginRegex, _         = regexp.Compile("You are currently logged in as <B>([\\s\\S]+?)</B>\\.")
	findNameRegex, _          = regexp.Compile("You are (?P<name>[\\S\\s]+?)<BR />")
	findContestInfoRegex, _   = regexp.Compile("Assigned contest: ([\\d\\D]+?) \\(([\\d\\D]+?)\\)")
	findProblemsBlockRegex, _ = regexp.Compile("<SELECT NAME=prob>[\\d\\D]+?</SELECT>")
	findLangsBlockRegex, _    = regexp.Compile("<SELECT NAME=lang>[\\d\\D]+?</SELECT>")
	findOptionsRegex, _       = regexp.Compile("<OPTION VALUE=([a-zA-Z0-9]+)[\\d\\D]*?>([\\d\\D]+?)</OPTION>")
	findLangAndNameRegex, _   = regexp.Compile("([\\d\\D]+?): ([\\d\\D]+)")
	findTableRowRegex, _      = regexp.Compile("<tr>([\\d\\D]+?)</tr>")
	findTableColumnRegex, _   = regexp.Compile("<td>([\\d\\D]+?)</td>")
	removeAnchorRegex, _      = regexp.Compile("<a([\\d\\D]*?)>([\\d\\D]+?)<")
	getLinkRegex, _           = regexp.Compile("href=\"([\\d\\D]*?)\"")
)

func (c *Client) FindLogin() (string, error) {
	body, err := util.GetBody(c.client, HOST+"/index")
	res := findLoginRegex.FindSubmatch(body)
	if err != nil || res == nil {
		return "", err
	}
	return string(res[1]), nil
}
func (c *Client) FindName() (string, error) {
	body, err := util.GetBody(c.client, HOST+"/index")
	if err != nil {
		return "", err
	}
	return string(findNameRegex.FindAllSubmatch(body, 2)[1][1]), nil
}
func (c *Client) FindContestID() (string, error) {
	body, err := util.GetBody(c.client, HOST+"/index")
	if err != nil {
		return "", err
	}
	return string(findContestInfoRegex.FindSubmatch(body)[1]), nil
}
func (c *Client) FindContestName() (string, error) {
	body, err := util.GetBody(c.client, HOST+"/index")
	if err != nil {
		return "", err
	}
	return string(findContestInfoRegex.FindSubmatch(body)[2]), nil
}

func (c *Client) GetProblems() (res []config.Problem, err error) {
	body, err := util.GetBody(c.client, HOST+"/submit")
	if err != nil {
		return nil, err
	}
	langsBlock := findProblemsBlockRegex.Find(body)
	for _, v := range findOptionsRegex.FindAllSubmatch(langsBlock, -1) {
		res = append(res, config.Problem{
			ProblemId:   string(v[1]),
			ProblemName: string(v[2]),
		})
	}
	return
}
func (c *Client) GetCompilers() (res []config.Compiler, err error) {
	body, err := util.GetBody(c.client, HOST+"/submit")
	if err != nil {
		return nil, err
	}
	langsBlock := findLangsBlockRegex.Find(body)
	for _, v := range findOptionsRegex.FindAllSubmatch(langsBlock, -1) {
		vv := findLangAndNameRegex.FindSubmatch(v[2])
		res = append(res, config.Compiler{
			CompilerId:   string(v[1]),
			CompilerName: string(vv[2]),
			CompilerLang: string(vv[1]),
		})
	}
	return
}
func (c *Client) GetConfig(path string) (cfg *config.Config, err error) {
	cfg, err = config.NewConfig(path)
	if cfg.Contest != "" {
		err := c.ChangeContest(cfg.Contest)
		if err != nil {
			fmt.Println(aurora.Red("Can't change contest"))
			return nil, err
		}
	}
	if err != nil {
		return
	}
	contestId, _ := c.FindContestID()
	fmt.Println(aurora.Sprintf(aurora.Cyan("Fetching problems for contest %v"), contestId))
	problems, err := c.GetProblems()
	if err != nil {
		fmt.Println(aurora.Red("Error occurred while problems were being fetched"))
		fmt.Println(err.Error())
	}
	cfg.Problems = problems
	compilers, err := c.GetCompilers()
	if err != nil {
		fmt.Println(aurora.Red("Error occurred while compilers were being fetched"))
		fmt.Println(err.Error())
	}
	cfg.Compilers = compilers
	cfg.Contest = contestId
	err = cfg.Save()
	if err != nil {
		fmt.Println(aurora.Red("Error occurred while Config was being saved"))
		fmt.Println(err.Error())
	}
	return
}

func (c *Client) GetAvailableContests() (res []Contest, err error) {
	body, err := util.GetBody(c.client, HOST+"/contests?mask=1")
	if err != nil {
		return
	}
	rows := findTableRowRegex.FindAllSubmatch(body, -1)
	for _, row := range rows {
		columns := findTableColumnRegex.FindAllStringSubmatch(string(row[1]), -1)
		res = append(res, Contest{
			ContestId:           removeAnchorRegex.FindStringSubmatch(columns[0][0])[2],
			ContestStatementURL: getLinkRegex.FindStringSubmatch(removeAnchorRegex.FindStringSubmatch(columns[1][0])[1])[1],
			ContestName:         removeAnchorRegex.FindStringSubmatch(columns[1][0])[2],
			ContestStarted:      columns[3][1],
			ContestStatus:       columns[2][1],
		})
	}
	return
}

func (c *Client) DownloadStatements(contest Contest, cfg *config.Config) error {
	body, err := util.GetBinary(c.client, contest.ContestStatementURL)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(cfg.Path), os.ModePerm); err != nil {
		return err
	}

	return os.WriteFile(filepath.Dir(cfg.Path)+"/statements.pdf", body, 0644)

}
