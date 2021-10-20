package client

import (
	"github.com/khbminus/tscli/config"
	"github.com/khbminus/tscli/cookiejar"
	"net/http"
)

type Client struct {
	Jar      *cookiejar.Jar `json:"cookies"`
	User     string         `json:"login"`
	Password string         `json:"password"`
	client   *http.Client
	path     string
}
type Submission struct {
	ID       string
	Problem  *config.Problem
	Attempt  string
	Time     string
	Compiler *config.Compiler
	Result   string
}
type Test struct {
	ID         string
	Result     string
	TimeUsed   string
	MemoryUsed string
	Comment    string
}
type Contest struct {
	ContestId           string
	ContestName         string
	ContestStatus       string
	ContestStarted      string
	ContestStatementURL string
}
