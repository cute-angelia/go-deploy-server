package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/cute-angelia/go-utils/logger"
	"go-deploy/pkg/consts"
	"go-deploy/pkg/utils"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var listening *string
var debug *string
var usage = `Usage: /pathto/client -l :8081 -d false`

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
	}
	listening = flag.String("l", ":8081", usage)
	debug = flag.String("d", "true", usage)
	flag.Parse()
	if *listening == "" {
		flag.Usage()
		os.Exit(1)
	}

	// 日志
	logger.NewLogger("go-deploy-client", *debug == "false")

	//start tcp server
	log.Printf("Start tcp server listening %s", *listening)
	ln, err := net.Listen("tcp", *listening)
	if err != nil {
		log.Println("Error listening:", err)
		os.Exit(1)
	}
	defer ln.Close()

	// run loop forever (or until ctrl-c)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Error accepting: ", err)
			continue
		}
		log.Printf("Received new connection %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		// will listen for message to process ending in newline (\n)
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println("Server closed", err.Error())
			break
		}
		// output message received
		if *debug == "true" {
			log.Print(conn.RemoteAddr(), " -> Message Received:", message)
		}

		if strings.TrimSpace(message) == "PING" {
			message = "PONG"
			// send new string back to client
			conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
			_, err = conn.Write([]byte(message + "\n"))
			if err != nil {
				log.Println("Error writing to stream.", err)
				break
			}
		} else {
			ret, err := processTask(message)
			if err != nil {
				log.Println("Process error", err)
				// ret = []byte(err.Error())
			}

			//replace \n with special chars
			ret = bytes.Replace(ret, []byte{10}, []byte("{CRLF}"), -1)
			ret = append(ret, '\n')
			conn.SetWriteDeadline(time.Now().Add(3 * time.Second))
			_, err = conn.Write(ret)
			if err != nil {
				log.Println("Error writing to stream:", err)
				break
			}
		}
	}
}

// git reset --hard 4f32685 || svn up -r 999
// git pull ||svn up
func processTask(message string) ([]byte, error) {
	msg := &consts.CmdMessage{}
	err := json.Unmarshal([]byte(message), msg)
	if err != nil {
		log.Print("Json decode error: " + err.Error())
		return nil, err
	}

	var command string
	if msg.Action == "update" {
		if msg.Type == "git" {
			command = fmt.Sprintf("cd %s && git pull", msg.Path)
		} else {
			command = fmt.Sprintf("cd %s && svn up", msg.Path)
		}
	} else if msg.Action == "rollback" {
		if msg.Type == "git" {
			command = fmt.Sprintf("cd %s && git reset --hard %s", msg.Path, msg.Reversion)
		} else {
			command = fmt.Sprintf("cd %s && svn up -r %s", msg.Path, msg.Reversion)
		}
	} else if msg.Action == "showlog" {
		if msg.Type == "git" {
			command = fmt.Sprintf(`cd %s && git log -20 --pretty="%%h %s %%an %s %%at %s %%s"`, msg.Path, consts.Separator, consts.Separator, consts.Separator)
		} else {
			command = fmt.Sprintf("svn log --limit 20")
		}
	}

	if command != "" {
		bytes := make([]byte, 0)

		//exec pre script
		if strings.TrimSpace(msg.BeforDeploy) != "" {
			log.Println("exec pre command:", command)
			byt, err := utils.RunShellCmd(msg.BeforDeploy)
			bytes = append(bytes, byt...)
			if err != nil {
				return bytes, err
			}
		}

		//exec command
		byt, err := utils.RunShellCmd(command)
		bytes = append(bytes, byt...)
		if err != nil {
			return bytes, err
		}

		//exec post script
		if strings.TrimSpace(msg.AfterDeploy) != "" {
			log.Println("exec post command:", command)
			byt, err := utils.RunShellCmd(msg.AfterDeploy)
			bytes = append(bytes, byt...)
			if err != nil {
				return bytes, err
			}
		}
		return bytes, nil
	}
	return nil, errors.New("command invalid")
}
