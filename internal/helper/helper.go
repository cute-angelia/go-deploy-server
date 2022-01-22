package helper

import (
	"bytes"
	"errors"
	"github.com/go-cmd/cmd"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

func IsWin() bool {
	sysType := runtime.GOOS
	return sysType == "windows"
}

// commond dir
func FlapCmd(command string) (string, string) {
	cmdstr := command
	cdPath := ""
	if strings.Contains(command, "cd ") && strings.Contains(command, "&& ") {
		temps := strings.Split(command, "&&")
		cdPath = strings.TrimSpace(strings.Replace(temps[0], "cd ", "", -1))
		cmdstr = strings.TrimSpace(temps[1])
	}
	return cmdstr, cdPath
}

func RunShellCmd(command string) ([]byte, error) {
	// 优先处理 cd 的命令
	cmdstr, dir := FlapCmd(command)

	log.Println("cmd:", cmdstr)

	strs := strings.Split(cmdstr, " ")
	parms := []string{}
	if len(strs) > 1 {
		parms = strs[1:]
	}
	// Start a long-running process, capture stdout and stderr
	findCmd := cmd.NewCmdOptions(cmd.Options{Buffered: true}, strs[0], parms...)
	if len(dir) > 0 {
		findCmd.Dir = dir
	}
	statusChan := findCmd.Start() // non-blocking

	// Stop command after 1 hour
	go func() {
		<-time.After(3 * time.Minute)
		findCmd.Stop()
	}()

	// Block waiting for command to exit, be stopped, or be killed
	finalStatus := <-statusChan

	// log.Println(ijson.Pretty(finalStatus))

	// 处理结果
	rest := []string{}
	if len(finalStatus.Stderr) > 0 {
		rest = finalStatus.Stderr
		finalStatus.Error = errors.New("出错了！")
	} else {
		rest = finalStatus.Stdout
	}
	log.Println("rest:", rest)

	// 兼容并返回
	bf := bytes.Buffer{}
	for _, i2 := range rest {
		bf.WriteString(i2 + "\n")
	}
	return bf.Bytes(), finalStatus.Error
}

// Deprecated: 请使用 RunShellCmd
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
