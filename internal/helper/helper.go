package helper

import (
	"encoding/json"
	"errors"
	"log"
	"os/exec"
	"runtime"
)

func IsWin() bool {
	sysType := runtime.GOOS
	return sysType == "windows"
}

func RunShell(command string) ([]byte, error) {
	//cmd := exec.Command("/bin/bash", "-c", `ps -eaf|grep "nginx: master"|grep -v "grep"|awk '{print $2}'`)
	cmd := new(exec.Cmd)
	if IsWin() {
		log.Println("exec command:", "cmd /C", command)
		cmd = exec.Command("cmd", "/C", command)
		//cmd = exec.Command("powershell", "-c", command)
	} else {
		log.Println("exec command:", "/bin/bash -c", command)
		cmd = exec.Command("/bin/bash", "-c", command)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		return nil, errors.New(string(out))
	}
	// log.Println(string(out))
	return out, nil
}

func JsonResp(status bool, msg string, elapsed string, data interface{}) []byte {
	bytes, _ := json.Marshal(struct {
		Status  bool
		Msg     string
		Elapsed string
		Data    interface{}
	}{Status: status, Msg: msg, Elapsed: elapsed, Data: data})
	return bytes
}
