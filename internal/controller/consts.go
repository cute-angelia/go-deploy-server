package controller

var separator = "{CRLF}"
var dataFormat = "2006-01-02 15:04:05"

type LogEntity struct {
	Reversion string
	Author    string
	Time      string
	Content   string
}

type LogList []LogEntity

type jobExecResult struct {
	Err      error
	Message  string
	AppName  string
	NodeName string
}
