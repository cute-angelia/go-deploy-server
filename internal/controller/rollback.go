package controller

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go-deploy/config"
	"go-deploy/helper"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Rollback struct {
}

func (rs Rollback) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", rs.index)
	return r
}

func (rs Rollback) index(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ret := rollback(r.PostFormValue("groupid"), r.PostFormValue("reversion"))
	fmt.Fprintf(w, "%s\n", helper.JsonResp(true, "", strconv.FormatFloat(time.Since(start).Seconds(), 'f', 3, 64), strings.Replace(ret, separator, "\n", -1)))
}

//send a rollback message to the group nodes
func rollback(groupid string, reversion string) string {
	for _, app := range config.C.Apps {
		if app.GroupId == groupid {
			jobExecChan := make(chan jobExecResult, len(app.Node))
			chanLen := 0

			for _, node := range app.Node {
				if node.Online {
					jobBody, _ := json.Marshal(struct {
						Type        string `json:"type"`
						Path        string `json:"path"`
						Action      string `json:"action"`
						Reversion   string `json:"reversion"`
						BeforDeploy string `json:"befor_deploy"`
						AfterDeploy string `json:"after_deploy"`
					}{Type: node.Type, Path: node.Path, BeforDeploy: node.BeforDeploy, AfterDeploy: node.AfterDeploy, Action: "rollback", Reversion: reversion})

					chanLen++
					go dispatchJob(jobBody, jobExecChan, node.Addr, app.Name, node.Alias)
				}
			}

			resp := ""
			for i := 0; i < chanLen; i++ {
				exeRet := <-jobExecChan
				if exeRet.Err != nil {
					resp += fmt.Sprintf("[%s:%s]\nERROR: %s", exeRet.AppName, exeRet.NodeName, exeRet.Err.Error())
				} else {
					resp += fmt.Sprintf("[%s:%s]\n%s", exeRet.AppName, exeRet.NodeName, exeRet.Message)
				}
			}
			return resp
		}
	}
	return ""
}
