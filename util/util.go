package util

import (
	"bufio"
	"fmt"
	"github.com/logrusorgru/aurora/v3"
	"golang.org/x/text/encoding/charmap"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var (
	decoder = charmap.KOI8R.NewDecoder()
)

func Scanline() string {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text()
	}
	fmt.Print("Interrupted")
	os.Exit(1)
	return ""
}

func ScanlineTrim() string {
	return strings.TrimSpace(Scanline())
}

func PostForm(c *http.Client, URL string, data url.Values) ([]byte, error) {
	resp, err := c.PostForm(URL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	text, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return decoder.Bytes(text)
}

func GetBody(c *http.Client, URL string) ([]byte, error) {
	resp, err := c.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	text, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return decoder.Bytes(text)
}

func GetBinary(c *http.Client, URL string) ([]byte, error) {
	resp, err := c.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	text, err := ioutil.ReadAll(resp.Body)
	return text, err
}

func ChooseIndex(maxN int) (res int) {
	fmt.Print(aurora.Cyan("Please choose one index: "))
	for {
		index := ScanlineTrim()
		i, err := strconv.Atoi(index)
		if err == nil && i >= 0 && i < maxN {
			return i
		}
		fmt.Print(aurora.Red("Invalid index! Try again: "))
	}
}
