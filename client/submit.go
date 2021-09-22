package client

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/khbminus/tscli/config"
	"github.com/khbminus/tscli/util"
	. "github.com/logrusorgru/aurora/v3"
	"golang.org/x/text/encoding/charmap"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
)

var (
	findSubmissionsRegex, _ = regexp.Compile("<(TR|tr) (CLASS|class)[\\d\\D]*?>([\\d\\D]+?)</(TR|tr)>")
	findColumnsRegex, _     = regexp.Compile("<(TD|td)>([\\d\\D]+?)</(TD|td)>")
)

var (
	findErrorRegex, _ = regexp.Compile("CLASS=errm>([\\d\\D]+?)</")
	decoder           = charmap.KOI8R.NewDecoder()
)

func (c *Client) Submit(compiler config.Compiler, problem config.Problem, solution string) error {
	contestName, err := c.FindContestName()
	if err != nil {
		return err
	}

	fmt.Println(Sprintf(Cyan("Submitting task %v at contest \"%v\""), problem.ProblemName, contestName))
	fmt.Println(Sprintf(Yellow("Compiler: %v"), Magenta(compiler.CompilerName)))

	bodyRequest := &bytes.Buffer{}
	writer := multipart.NewWriter(bodyRequest)
	_ = writer.WriteField("prob", problem.ProblemId)
	_ = writer.WriteField("lang", compiler.CompilerId)
	_ = writer.WriteField("solution", solution)
	_ = writer.Close()
	r, _ := http.NewRequest("POST", HOST+"/submit", bodyRequest)
	r.Header.Add("Content-Type", writer.FormDataContentType())
	resp, _ := c.client.Do(r)

	defer resp.Body.Close()
	text, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	body, err := decoder.Bytes(text)
	if err != nil {
		return err
	}

	if tsError := findErrorRegex.FindSubmatch(body); tsError != nil {
		fmt.Println(Sprintf(Red("Submit error by TSweb: %v\n"), string(tsError[1])))
		return errors.New("TSWeb internal error")
	}
	fmt.Println(Green("Sent solution successfully!"))
	return err
}

func (c *Client) GetAllSubmits(cfg config.Config) (res []Submission, err error) {
	body, err := util.GetBody(c.client, HOST+"/allsubmits")
	if err != nil {
		fmt.Println(Red("Error occurred at getting submits"))
		fmt.Println(err.Error())
		return nil, err
	}
	for i, submit := range findSubmissionsRegex.FindAllSubmatch(body, -1) {
		fields := findColumnsRegex.FindAllStringSubmatch(string(submit[3]), -1)
		res = append(res, Submission{
			ID:      fields[0][2],
			Attempt: fields[2][2],
			Time:    fields[3][2],
			Result:  fields[5][2],
		})
		for _, compiler := range cfg.Compilers {
			if compiler.CompilerName == fields[4][2] {
				res[i].Compiler = &compiler
				break
			}
		}

		for _, problem := range cfg.Problems {
			if problem.ProblemId == fields[1][2] {
				res[i].Problem = &problem
				break
			}
		}
	}
	return
}

func (c *Client) GetFeedback(submit Submission) (res []Test, err error) {
	body, err := util.GetBody(c.client, fmt.Sprintf("%v?id=%v", HOST+"/feedback", submit.ID))
	if err != nil {
		fmt.Println(Red("Error occurred at getting feedback"))
		fmt.Println(err.Error())
		return nil, err
	}
	for _, v := range findSubmissionsRegex.FindAllSubmatch(body, -1) {
		test := findColumnsRegex.FindAllStringSubmatch(strings.ReplaceAll(string(v[3]), "&nbsp;", " "), -1)
		res = append(res, Test{
			ID:         test[0][2],
			Result:     test[1][2],
			TimeUsed:   test[2][2],
			MemoryUsed: test[3][2],
			Comment:    test[4][2],
		})
	}
	return
}
func (c *Client) ChangeContest(newContest string) (err error) {
	_, err = util.GetBody(c.client, fmt.Sprintf(HOST+"/index?op=changecontest&newcontestid=%v", newContest))
	if err != nil {
		return err
	}
	// TODO: parse error
	return err
}
