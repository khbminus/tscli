package cmd

import (
	"fmt"
	"github.com/khbminus/tscli/client"
	"github.com/khbminus/tscli/config"
	"github.com/logrusorgru/aurora/v3"
	"github.com/muesli/termenv"
	"github.com/olekukonko/tablewriter"
	"os"
	"time"
)

func SubmitAndWatch(problem config.Problem, compiler config.Compiler, solution string) error {
	err := client.Instance.Submit(compiler, problem, solution)
	if err != nil {
		fmt.Println(aurora.Red("Error while submitting!"))
		return err
	}

	cfg, err := GetConfig()
	submits, err := client.Instance.GetAllSubmits(*cfg)
	nowId := submits[0].ID
	if err != nil {
		fmt.Println(aurora.Red("Error while getting submits"))
		return err
	}
	termenv.Reset()
	termenv.ClearScreen()
	termenv.SaveCursorPosition()
	for {
		submits, err = client.Instance.GetAllSubmits(*cfg)
		if err != nil {
			fmt.Println(aurora.Red("Error while getting submits"))
			return err
		}
		index := -1
		for i, v := range submits {
			if v.ID == nowId {
				index = i
				break
			}
		}
		feedback, err := client.Instance.GetFeedback(submits[index])
		if err != nil {
			fmt.Println(aurora.Red("Error occurred while feedback was being watched..."))
			fmt.Println(aurora.Red(err.Error()))
			return err
		}
		termenv.RestoreCursorPosition()
		termenv.SaveCursorPosition()
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"#", "Time used", "Memory used", "Status", "Comment"})
		table.SetAutoWrapText(false)
		table.SetAlignment(tablewriter.ALIGN_CENTER)
		for _, test := range feedback {
			var status string
			switch test.Result {
			case "OK":
				status = aurora.Sprintf(aurora.Green(test.Result))
			case "ML":
			case "TL":
				status = aurora.Sprintf(aurora.Magenta(test.Result))
			case "WA":
			case "RT":
				status = aurora.Sprintf(aurora.Red(test.Result))
			default:
				status = test.Result
			}
			table.Append([]string{test.ID, test.TimeUsed, test.MemoryUsed, status, test.Comment})
		}
		table.Render()

		if submits[index].Result != "NO" {
			break
		}
		time.Sleep(time.Millisecond * 500)
	}
	return nil
}
