package controller

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go-deploy/config"
	"go-deploy/helper"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Deply struct {
}

func (rs Deply) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", rs.index)
	return r
}

func (rs Deply) index(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	ret := deply(r.PostFormValue("groupid"))
	fmt.Fprintf(w, "%s\n", helper.JsonResp(true, "", strconv.FormatFloat(time.Since(start).Seconds(), 'f', 3, 64), strings.Replace(ret, separator, "\n", -1)))
}

//send a update message to the group nodes
func deply(groupid string) string {
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
						BeforDeploy string `json:"befor_deploy"`
						AfterDeploy string `json:"after_deploy"`
					}{Type: node.Type, Path: node.Path, BeforDeploy: node.BeforDeploy, AfterDeploy: node.AfterDeploy, Action: "update"})

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

//send job to client and get execute result
func dispatchJob(jobBody []byte, jobExecChan chan jobExecResult, addr string, appName string, nodeName string) {
	execResult := jobExecResult{AppName: appName, NodeName: nodeName}
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println("Error connect to client:", err)
		execResult.Err = err
		jobExecChan <- execResult
		return
	}
	defer conn.Close()

	jobBody = append(jobBody, '\n')
	conn.SetWriteDeadline(time.Now().Add(30 * time.Second))
	_, err = conn.Write(jobBody)
	if err != nil {
		log.Println("Error writing to stream:", err)
		execResult.Err = err
		jobExecChan <- execResult
		return
	}

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println("Client closed", err.Error())
		execResult.Err = err
		jobExecChan <- execResult
		return
	}
	execResult.Message = message
	jobExecChan <- execResult
	return
}
