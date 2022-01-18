package controller

import (
	"errors"
	"fmt"
	"github.com/cute-angelia/go-utils/syntax/istrings"
	"github.com/go-chi/chi/v5"
	"go-deploy/config"
	"go-deploy/helper"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ShowLog struct {
}

func (rs ShowLog) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", rs.index)
	return r
}

func (rs ShowLog) index(w http.ResponseWriter, r *http.Request) {
	list, err := showlog(r.PostFormValue("groupid"))
	if err != nil {
		fmt.Fprintf(w, "%s\n", helper.JsonResp(false, err.Error(), "", nil))
	} else {
		fmt.Fprintf(w, "%s\n", helper.JsonResp(true, "", "", list))
	}
}

func showlog(groupid string) (list LogList, err error) {
	groupid = strings.TrimSpace(groupid)
	if groupid != "" {
		for _, app := range config.C.Apps {
			if app.GroupId == groupid {
				if app.Type == "svn" {
					return showSvnLog(app)
				} else if app.Type == "git" {
					return showGitLog(app)
				}
			}
		}
	}
	return nil, errors.New("groupid invalid")
}

//svn log --limit 100 svn://x.x.x.x/path
func showSvnLog(app config.Apps) (list LogList, err error) {
	bytes, err := helper.RunShell(fmt.Sprintf("svn log --limit 20 %s", app.Url))
	if err != nil {
		return nil, err
	} else {
		var convertStr string
		var svnlogRegex = new(regexp.Regexp)
		if helper.IsWin() {
			convertStr = istrings.GbkToUtf8(string(bytes))
			svnlogRegex = regexp.MustCompile(`r(\d+) \| (\w+) \| (.*) \+0800(?:.*)\r\n\r\n(.*)\r\n--`)
		} else {
			convertStr = string(bytes)
			svnlogRegex = regexp.MustCompile(`r(\d+) \| (\w+) \| (.*) \+0800(?:.*)\n\n(.*)\n--`)
		}

		// log.Println(convertStr)
		match := svnlogRegex.FindAllStringSubmatch(convertStr, -1)

		logList := make(LogList, 0)
		for _, item := range match {
			svnlog := LogEntity{Reversion: string(item[1]), Author: string(item[2]), Time: string(item[3]), Content: string(item[4])}
			logList = append(logList, svnlog)
		}
		return logList, nil
	}
}

//cd /pathto/xx && git log -50 --pretty="%h {CRLF} %an {CRLF} %at {CRLF} %s"
func showGitLog(app config.Apps) (list LogList, err error) {
	cmd := fmt.Sprintf(`cd %s && git log -20 --pretty="%%h %s %%an %s %%at %s %%s"`, app.Fetchlogpath, separator, separator, separator)
	bytes, err := helper.RunShell(cmd)
	if err != nil {
		return nil, err
	} else {
		logs := strings.Split(string(bytes), "\n")
		logList := make(LogList, 0)
		for _, line := range logs {
			if strings.TrimSpace(line) != "" {
				fmt.Println(line)
				commitLog := strings.Split(line, separator)
				timeInt64, err := strconv.ParseInt(strings.TrimSpace(commitLog[2]), 10, 64)
				if err != nil {
					timeInt64 = time.Now().Unix()
				}
				logList = append(logList, LogEntity{Reversion: commitLog[0], Author: commitLog[1], Time: time.Unix(timeInt64, 0).Format(dataFormat), Content: commitLog[3]})
			}
		}
		return logList, nil
	}
}
